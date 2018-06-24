package openvpn

import (
	"fmt"
	"strings"

	"github.com/operator-framework/operator-sdk/pkg/sdk"
	api "github.com/pieterlange/openvpn-operator/pkg/apis/ptlc/v1alpha1"
	"k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// prepareOpenVPNConfig applies our settings to an openvpn configmap
// - If given user configmap, appends into user provided openvpn config
//   and creates another configmap "${configMapName}-copy" for it.
// - Otherwise, creates a new configmap "${openvpnName}-copy" with our section.
func prepareOpenVPNConfig(vr *api.OpenVPN) error {
	cm := &v1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: vr.Namespace,
		},
	}
	if len(vr.Spec.ConfigMapName) != 0 {
		cm.Name = vr.Spec.ConfigMapName
		err := sdk.Get(cm)
		if err != nil {
			return fmt.Errorf("prepare openvpn config error: get configmap (%s) failed: %v", vr.Spec.ConfigMapName, err)
		}
		// DO something with this at some point
	}
	cm.Name = configMapNameForOpenVPN(vr)
	cm.Labels = labelsForOpenVPN(vr.Name)
	cm.Data = map[string]string{
		"OVPN_SERVER_URL":          openVPNEndpointURL(vr.Spec.PublicEndpoint.Protocol, vr.Spec.PublicEndpoint.Hostname, vr.Spec.PublicEndpoint.Port),
		"OVPN_K8S_SERVICE_NETWORK": vr.Spec.ServiceCIDR,
		"OVPN_K8S_POD_NETWORK":     vr.Spec.PodCIDR,
		"OVPN_K8S_DOMAIN":          vr.Spec.ClusterDomain,
		"OVPN_STATUS":              vr.Spec.StatusFilePath,
	}

	addOwnerRefToObject(cm, asOwner(vr))
	err := sdk.Create(cm)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		return fmt.Errorf("prepare OpenVPN settings error: create new configmap (%s) failed: %v", cm.Name, err)
	}
	return nil
}

func openVPNEndpointURL(proto string, hostname string, port int32) string {
	protocol := strings.ToLower(proto)
	return fmt.Sprintf("%s://%s:%d", protocol, hostname, port)
}

func configMapNameForOpenVPN(v *api.OpenVPN) string {
	n := v.Spec.ConfigMapName
	if len(n) == 0 {
		n = v.Name
	}
	return n + "-copy"
}
