// +build !ignore_autogenerated

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Images) DeepCopyInto(out *Images) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Images.
func (in *Images) DeepCopy() *Images {
	if in == nil {
		return nil
	}
	out := new(Images)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OpenVPN) DeepCopyInto(out *OpenVPN) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OpenVPN.
func (in *OpenVPN) DeepCopy() *OpenVPN {
	if in == nil {
		return nil
	}
	out := new(OpenVPN)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *OpenVPN) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OpenVPNList) DeepCopyInto(out *OpenVPNList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]OpenVPN, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OpenVPNList.
func (in *OpenVPNList) DeepCopy() *OpenVPNList {
	if in == nil {
		return nil
	}
	out := new(OpenVPNList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *OpenVPNList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OpenVPNSpec) DeepCopyInto(out *OpenVPNSpec) {
	*out = *in
	out.PublicEndpoint = in.PublicEndpoint
	out.Images = in.Images
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OpenVPNSpec.
func (in *OpenVPNSpec) DeepCopy() *OpenVPNSpec {
	if in == nil {
		return nil
	}
	out := new(OpenVPNSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OpenVPNStatus) DeepCopyInto(out *OpenVPNStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OpenVPNStatus.
func (in *OpenVPNStatus) DeepCopy() *OpenVPNStatus {
	if in == nil {
		return nil
	}
	out := new(OpenVPNStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PublicEndpoint) DeepCopyInto(out *PublicEndpoint) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PublicEndpoint.
func (in *PublicEndpoint) DeepCopy() *PublicEndpoint {
	if in == nil {
		return nil
	}
	out := new(PublicEndpoint)
	in.DeepCopyInto(out)
	return out
}