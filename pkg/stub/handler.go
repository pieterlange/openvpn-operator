package stub

import (
	"context"

	"github.com/pieterlange/openvpn-operator/pkg/apis/ptlc/v1alpha1"
	"github.com/pieterlange/openvpn-operator/pkg/openvpn"

	"github.com/operator-framework/operator-sdk/pkg/sdk"
)

func NewHandler() sdk.Handler {
	return &Handler{}
}

type Handler struct {
	// Fill me
}

func (h *Handler) Handle(ctx context.Context, event sdk.Event) error {
	switch o := event.Object.(type) {
	case *v1alpha1.OpenVPN:
		return openvpn.Reconcile(o)
	}
	return nil
}
