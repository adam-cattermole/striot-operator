// +build !ignore_autogenerated

// Code generated by operator-sdk. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConnectType) DeepCopyInto(out *ConnectType) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConnectType.
func (in *ConnectType) DeepCopy() *ConnectType {
	if in == nil {
		return nil
	}
	out := new(ConnectType)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Partition) DeepCopyInto(out *Partition) {
	*out = *in
	out.ConnectType = in.ConnectType
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Partition.
func (in *Partition) DeepCopy() *Partition {
	if in == nil {
		return nil
	}
	out := new(Partition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Topology) DeepCopyInto(out *Topology) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Topology.
func (in *Topology) DeepCopy() *Topology {
	if in == nil {
		return nil
	}
	out := new(Topology)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Topology) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TopologyList) DeepCopyInto(out *TopologyList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Topology, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TopologyList.
func (in *TopologyList) DeepCopy() *TopologyList {
	if in == nil {
		return nil
	}
	out := new(TopologyList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *TopologyList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TopologySpec) DeepCopyInto(out *TopologySpec) {
	*out = *in
	if in.Partitions != nil {
		in, out := &in.Partitions, &out.Partitions
		*out = make([]Partition, len(*in))
		copy(*out, *in)
	}
	if in.Order != nil {
		in, out := &in.Order, &out.Order
		*out = make([]int32, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TopologySpec.
func (in *TopologySpec) DeepCopy() *TopologySpec {
	if in == nil {
		return nil
	}
	out := new(TopologySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TopologyStatus) DeepCopyInto(out *TopologyStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TopologyStatus.
func (in *TopologyStatus) DeepCopy() *TopologyStatus {
	if in == nil {
		return nil
	}
	out := new(TopologyStatus)
	in.DeepCopyInto(out)
	return out
}
