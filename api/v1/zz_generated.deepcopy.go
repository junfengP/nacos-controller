//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by controller-gen. DO NOT EDIT.

package v1

import (
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AdditionalConfiguration) DeepCopyInto(out *AdditionalConfiguration) {
	*out = *in
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Properties != nil {
		in, out := &in.Properties, &out.Properties
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Tags != nil {
		in, out := &in.Tags, &out.Tags
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AdditionalConfiguration.
func (in *AdditionalConfiguration) DeepCopy() *AdditionalConfiguration {
	if in == nil {
		return nil
	}
	out := new(AdditionalConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DynamicConfiguration) DeepCopyInto(out *DynamicConfiguration) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DynamicConfiguration.
func (in *DynamicConfiguration) DeepCopy() *DynamicConfiguration {
	if in == nil {
		return nil
	}
	out := new(DynamicConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DynamicConfiguration) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DynamicConfigurationList) DeepCopyInto(out *DynamicConfigurationList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]DynamicConfiguration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DynamicConfigurationList.
func (in *DynamicConfigurationList) DeepCopy() *DynamicConfigurationList {
	if in == nil {
		return nil
	}
	out := new(DynamicConfigurationList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DynamicConfigurationList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DynamicConfigurationSpec) DeepCopyInto(out *DynamicConfigurationSpec) {
	*out = *in
	if in.DataIds != nil {
		in, out := &in.DataIds, &out.DataIds
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.AdditionalConf != nil {
		in, out := &in.AdditionalConf, &out.AdditionalConf
		*out = new(AdditionalConfiguration)
		(*in).DeepCopyInto(*out)
	}
	out.Strategy = in.Strategy
	in.NacosServer.DeepCopyInto(&out.NacosServer)
	out.ObjectRef = in.ObjectRef
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DynamicConfigurationSpec.
func (in *DynamicConfigurationSpec) DeepCopy() *DynamicConfigurationSpec {
	if in == nil {
		return nil
	}
	out := new(DynamicConfigurationSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DynamicConfigurationStatus) DeepCopyInto(out *DynamicConfigurationStatus) {
	*out = *in
	if in.SyncStatuses != nil {
		in, out := &in.SyncStatuses, &out.SyncStatuses
		*out = make([]SyncStatus, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	out.ObjectRef = in.ObjectRef
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DynamicConfigurationStatus.
func (in *DynamicConfigurationStatus) DeepCopy() *DynamicConfigurationStatus {
	if in == nil {
		return nil
	}
	out := new(DynamicConfigurationStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NacosServerConfiguration) DeepCopyInto(out *NacosServerConfiguration) {
	*out = *in
	if in.Endpoint != nil {
		in, out := &in.Endpoint, &out.Endpoint
		*out = new(string)
		**out = **in
	}
	if in.ServerAddr != nil {
		in, out := &in.ServerAddr, &out.ServerAddr
		*out = new(string)
		**out = **in
	}
	out.AuthRef = in.AuthRef
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NacosServerConfiguration.
func (in *NacosServerConfiguration) DeepCopy() *NacosServerConfiguration {
	if in == nil {
		return nil
	}
	out := new(NacosServerConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SyncStatus) DeepCopyInto(out *SyncStatus) {
	*out = *in
	in.LastSyncTime.DeepCopyInto(&out.LastSyncTime)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SyncStatus.
func (in *SyncStatus) DeepCopy() *SyncStatus {
	if in == nil {
		return nil
	}
	out := new(SyncStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SyncStrategy) DeepCopyInto(out *SyncStrategy) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SyncStrategy.
func (in *SyncStrategy) DeepCopy() *SyncStrategy {
	if in == nil {
		return nil
	}
	out := new(SyncStrategy)
	in.DeepCopyInto(out)
	return out
}
