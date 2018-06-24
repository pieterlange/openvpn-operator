package openvpn

import (
	"fmt"
	"path/filepath"
	"strings"

	api "github.com/pieterlange/openvpn-operator/pkg/apis/ptlc/v1alpha1"

	"github.com/operator-framework/operator-sdk/pkg/sdk"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	openvpnPKIVolName     = "openvpn-pki"
	openvpnCRLVolName     = "openvpn-crl"
	openvpnCCDVolName     = "openvpn-ccd"
	openvpnPortMapVolName = "openvpn-portmap"
	openvpnStatusVolName  = "openvpn-status"
	openvpnPKIVolPath     = "/etc/openvpn/pki/"
	openvpnCRLVolPath     = "/etc/openvpn/crl/"
	openvpnCCDVolPath     = "/etc/openvpn/ccd/"
	openvpnPortMapVolPath = "/etc/openvpn/portmapping/"
	openvpnStatusVolPath  = "/etc/openvpn/status/"
	openvpnPort           = 1194
	exporterPort          = 9176
)

// deployOpenVPN deploys an OpenVPN Server.
// deployOpenVPN is a multi-step process. It creates the a deployment, the service
//
// deployOpenVPN is idempotent. If an object already exists, this function will ignore creating
// it and return no error. It is safe to retry on this function.
func deployOpenVPN(v *api.OpenVPN) error {
	selector := labelsForOpenVPN(v.GetName())

	replicas := int32(1)

	podTempl := v1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Name:      v.GetName(),
			Namespace: v.GetNamespace(),
			Labels:    selector,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{openvpnContainer(v), exporterContainer(v)},
			Volumes: []v1.Volume{
				{
					Name: openvpnPKIVolName,
					VolumeSource: v1.VolumeSource{
						Secret: &v1.SecretVolumeSource{
							SecretName: v.GetName() + "-pki",
						},
					},
				},
				{
					Name: openvpnStatusVolName,
					VolumeSource: v1.VolumeSource{
						EmptyDir: &v1.EmptyDirVolumeSource{},
					},
				},
			},
		},
	}

	d := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      v.GetName(),
			Namespace: v.GetNamespace(),
			Labels:    selector,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{MatchLabels: selector},
			Template: podTempl,
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RollingUpdateDeploymentStrategyType,
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxUnavailable: func(a intstr.IntOrString) *intstr.IntOrString { return &a }(intstr.FromInt(1)),
					MaxSurge:       func(a intstr.IntOrString) *intstr.IntOrString { return &a }(intstr.FromInt(1)),
				},
			},
		},
	}
	addOwnerRefToObject(d, asOwner(v))
	err := sdk.Create(d)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		return err
	}

	svcAnnotations := annotationsForLoadbalancer(v.Spec.PublicEndpoint.Hostname)
	svc := &v1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        v.GetName(),
			Namespace:   v.GetNamespace(),
			Annotations: svcAnnotations,
			Labels:      selector,
		},
		Spec: v1.ServiceSpec{
			Selector: selector,
			Type:     "LoadBalancer",
			Ports: []v1.ServicePort{
				{
					Name:       "openvpn",
					Protocol:   v1.Protocol(strings.ToUpper(v.Spec.PublicEndpoint.Protocol)),
					Port:       v.Spec.PublicEndpoint.Port,
					TargetPort: func(a intstr.IntOrString) intstr.IntOrString { return a }(intstr.FromInt(openvpnPort)),
				},
			},
		},
	}
	addOwnerRefToObject(svc, asOwner(v))
	err = sdk.Create(svc)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		return fmt.Errorf("failed to create openvpn service: %v", err)
	}
	return nil
}

func openvpnContainer(v *api.OpenVPN) v1.Container {
	return v1.Container{
		Name:  "openvpn",
		Image: v.Spec.Images.OpenVPNImage,
		Env: []v1.EnvVar{
			{
				Name: "PODIPADDR",
				ValueFrom: &v1.EnvVarSource{
					FieldRef: &v1.ObjectFieldSelector{
						FieldPath: "status.podIP",
					},
				},
			},
		},
		EnvFrom: []v1.EnvFromSource{{
			ConfigMapRef: &v1.ConfigMapEnvSource{
				LocalObjectReference: v1.LocalObjectReference{
					Name: configMapNameForOpenVPN(v)}}}},
		VolumeMounts: []v1.VolumeMount{
			{
				Name:      openvpnPKIVolName,
				MountPath: filepath.Dir(openvpnPKIVolPath),
			},
			{
				Name:      openvpnStatusVolName,
				MountPath: filepath.Dir(openvpnStatusVolPath),
			},
			//			{
			//				Name:      openvpnCRLVolName,
			//				MountPath: filepath.Dir(openvpnCRLVolPath),
			//		    },
			//			{
			//				Name:      openvpnCCDVolName,
			//				MountPath: filepath.Dir(openvpnCCDVolPath),
			//		    },
			//			{
			//				Name:      openvpnPortMapVolName,
			//				MountPath: filepath.Dir(openvpnPortMapVolPath),
			//		    },

		},
		SecurityContext: &v1.SecurityContext{
			Capabilities: &v1.Capabilities{
				// OpenVPN requires NET_ADMIN in order to create a TUN device
				Add: []v1.Capability{"NET_ADMIN"},
			},
		},
		Ports: []v1.ContainerPort{{
			Name:          "openvpn",
			ContainerPort: int32(openvpnPort),
		}},
	}
}

func exporterContainer(v *api.OpenVPN) v1.Container {
	return v1.Container{
		Name:  "metrics",
		Image: v.Spec.Images.OpenVPNExporterImage,
		VolumeMounts: []v1.VolumeMount{
			{
				Name:      openvpnStatusVolName,
				MountPath: filepath.Dir(openvpnStatusVolPath),
			},
		},
		Ports: []v1.ContainerPort{{
			Name:          "metrics",
			ContainerPort: int32(exporterPort),
		}},
	}
}
