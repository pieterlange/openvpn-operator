package openvpn

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	action "github.com/operator-framework/operator-sdk/pkg/sdk"
	api "github.com/pieterlange/openvpn-operator/pkg/apis/ptlc/v1alpha1"
)

// Reconcile reconciles the openvpn cluster's state to the spec specified by vr
// by preparing the TLS secrets, deploying the etcd and openvpn cluster,
// and finally updating the openvpn deployment if needed.
func Reconcile(vr *api.OpenVPN) (err error) {
	vr = vr.DeepCopy()

	logrus.Infof("Got new object %+v\n", vr)
	changed := vr.SetDefaults()
	if changed {
		return action.Update(vr)
	}
	logrus.Infof("After setting defaults: %+v\n", vr)

	err = prepareOpenVPNConfig(vr)
	if err != nil {
		return fmt.Errorf("Something went wrong setting up the PKI: %v", err)
	}
	logrus.Infof("Created OpenVPN configuration for %v", vr.Name)

	logrus.Infof("Initializing PKI for %s", vr.Name)
	err = prepareDefaultOpenVPNPKI(vr)
	if err != nil {
		return fmt.Errorf("Something went wrong setting up the PKI: %v", err)
	}

	err = deployOpenVPN(vr)
	if err != nil {
		return err
	}
	//
	//	vcs, err := getOpenVPNStatus(vr)
	//	if err != nil {
	//		return err
	//	}
	return nil
	//return updateopenvpnStatus(vr, vcs)
}
