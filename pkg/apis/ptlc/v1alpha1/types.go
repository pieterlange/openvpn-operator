package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	defaultOpenVPNImage         = "quay.io/plange/openvpn:no-dh"
	defaultOpenVPNExporterImage = "quay.io/plange/openvpn_exporter:latest"
	defaultClusterDomain        = "svc.cluster.local"
	defaultStatusFilePath       = "/etc/openvpn/status/server.status"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type OpenVPNList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []OpenVPN `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type OpenVPN struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              OpenVPNSpec   `json:"spec"`
	Status            OpenVPNStatus `json:"status,omitempty"`
}

func (v *OpenVPN) SetDefaults() bool {
	changed := false
	ovpn := &v.Spec
	//	if len(ovpn.Issuer) == 0 {
	//		ovpn.Issuer = v.Name + "-pki"
	//		changed = true
	//	}
	if len(ovpn.Images.OpenVPNImage) == 0 {
		ovpn.Images.OpenVPNImage = defaultOpenVPNImage
		changed = true
	}
	if len(ovpn.Images.OpenVPNExporterImage) == 0 {
		ovpn.Images.OpenVPNExporterImage = defaultOpenVPNExporterImage
		changed = true
	}
	if len(ovpn.ClusterDomain) == 0 {
		ovpn.ClusterDomain = defaultClusterDomain
		changed = true
	}
	return changed
}

type OpenVPNSpec struct {
	//	Issuer         string         `json:"issuer"`
	PublicEndpoint PublicEndpoint `json:"publicEndpoint"`
	ServiceCIDR    string         `json:"serviceCIDR"`
	PodCIDR        string         `json:"podCIDR"`
	Routes         string         `json:"routes"`
	ClusterDomain  string         `json:"clusterDomain"`
	StatusFilePath string         `json:"statusFile"`
	Images         Images         `json:"images"`
	// Name of the ConfigMap for OpenVPN's configuration
	// If this is empty, operator will create a default config for OpenVPN.
	// If this is not empty, operator will create a new config overwriting
	// the CIDR, server url, domain and statusfile settings
	ConfigMapName string `json:"configMapName"`
}

type Images struct {
	OpenVPNImage         string `json:"openvpn"`
	OpenVPNExporterImage string `json:"openvpnExporter"`
}
type PublicEndpoint struct {
	Protocol string `json:"protocol"`
	Hostname string `json:"hostname"`
	Port     int32  `json:"port"`
}

type OpenVPNStatus struct {
	ActiveClients int32 `json:"activeClients"`
}
