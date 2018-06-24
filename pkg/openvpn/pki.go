package openvpn

import (
	"crypto/rsa"
	"crypto/x509"
	"fmt"

	action "github.com/operator-framework/operator-sdk/pkg/sdk"

	api "github.com/pieterlange/openvpn-operator/pkg/apis/ptlc/v1alpha1"
	"github.com/pieterlange/openvpn-operator/pkg/pki"
	"k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	orgForTLSCert = []string{"openvpn.ptlc.nl"}
)

func prepareDefaultOpenVPNPKI(vr *api.OpenVPN) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("prepare default OpenVPN TLS secrets failed: %v", err)
		}
	}()

	caKey, caCrt, err := newCACert()
	if err != nil {
		return err
	}
	se, err := newOpenVPNServerPKISecret(vr, caKey, caCrt)
	if err != nil {
		return err
	}
	addOwnerRefToObject(se, asOwner(vr))
	err = action.Create(se)

	if err != nil && !apierrors.IsAlreadyExists(err) {
		return err
	}

	return nil
}

func newOpenVPNServerPKISecret(vr *api.OpenVPN, caKey *rsa.PrivateKey, caCrt *x509.Certificate) (*v1.Secret, error) {
	return newTLSSecret(vr, caKey, caCrt, vr.Spec.PublicEndpoint.Hostname, vr.Name+"-pki",
		map[string]string{
			"key":  "private.key",
			"cert": "certificate.crt",
			"ca":   "ca.crt",
			"ta":   "ta.key",
		})
}

func newCACert() (*rsa.PrivateKey, *x509.Certificate, error) {
	key, err := pki.NewPrivateKey()
	if err != nil {
		return nil, nil, err
	}

	config := pki.CertConfig{
		CommonName:   "OpenVPN operator CA",
		Organization: orgForTLSCert,
	}

	cert, err := pki.NewSelfSignedCACertificate(config, key)
	if err != nil {
		return nil, nil, err
	}

	return key, cert, err
}

// newTLSSecret is a common utility for creating a secret containing TLS assets.
func newTLSSecret(vr *api.OpenVPN, caKey *rsa.PrivateKey, caCrt *x509.Certificate, commonName, secretName string,
	fieldMap map[string]string) (*v1.Secret, error) {
	tc := pki.CertConfig{
		CommonName:   commonName,
		Organization: orgForTLSCert,
	}
	key, crt, err := newKeyAndCert(caCrt, caKey, tc)
	if err != nil {
		return nil, fmt.Errorf("creating new TLS secret failed: %v", err)
	}
	ta, err := pki.NewPrivateKey()
	if err != nil {
		return nil, fmt.Errorf("creating new private key for TLS-auth failed: %v", err)
	}

	secret := &v1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: vr.Namespace,
			Labels:    labelsForOpenVPN(vr.Spec.PublicEndpoint.Hostname),
		},
		Data: map[string][]byte{
			fieldMap["key"]:  pki.EncodePrivateKeyPEM(key),
			fieldMap["cert"]: pki.EncodeCertificatePEM(crt),
			fieldMap["ca"]:   pki.EncodeCertificatePEM(caCrt),
			fieldMap["ta"]:   pki.EncodePrivateKeyPEM(ta),
		},
	}
	return secret, nil
}

func newKeyAndCert(caCert *x509.Certificate, caPrivKey *rsa.PrivateKey, config pki.CertConfig) (*rsa.PrivateKey, *x509.Certificate, error) {
	key, err := pki.NewPrivateKey()
	if err != nil {
		return nil, nil, err
	}
	cert, err := pki.NewSignedCertificate(config, key, caCert, caPrivKey)
	if err != nil {
		return nil, nil, err
	}
	return key, cert, nil
}
