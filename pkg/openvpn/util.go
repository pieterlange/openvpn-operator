package openvpn

import (
	api "github.com/pieterlange/openvpn-operator/pkg/apis/ptlc/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// addOwnerRefToObject appends the desired OwnerReference to the object
func addOwnerRefToObject(o metav1.Object, r metav1.OwnerReference) {
	o.SetOwnerReferences(append(o.GetOwnerReferences(), r))
}

// labelsForOpenVPN returns the labels for selecting the resources
// belonging to the given openvpn server
func labelsForOpenVPN(hostname string) map[string]string {
	return map[string]string{"app": "openvpn", "openvpn": hostname}
}

// annotationsForLoadbalancer returns the annotation for the loadbalancer service
func annotationsForLoadbalancer(hostname string) map[string]string {
	return map[string]string{"external-dns.alpha.kubernetes.io/hostname": hostname}
}

// asOwner returns an owner reference set as the vault cluster CR
func asOwner(v *api.OpenVPN) metav1.OwnerReference {
	trueVar := true
	return metav1.OwnerReference{
		APIVersion: api.SchemeGroupVersion.String(),
		Kind:       api.OpenVPNKind,
		Name:       v.Name,
		UID:        v.UID,
		Controller: &trueVar,
	}
}
