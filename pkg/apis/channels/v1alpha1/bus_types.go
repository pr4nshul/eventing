/*
 * Copyright 2018 The Knative Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package v1alpha1

import (
	"encoding/json"

	kapi "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:defaulter-gen=true

// Bus represents how channels and subscriptions should be managed and
// corresponds to the buses.channels.knative.dev CRD. Buses will frequently, but
// not always, be backed by an event broker.
type Bus struct {
	meta_v1.TypeMeta   `json:",inline"`
	meta_v1.ObjectMeta `json:"metadata"`
	Spec               BusSpec   `json:"spec"`
	Status             BusStatus `json:"status,omitempty"`
}

// BusSpec specifies the Bus' parameters for Channels and Subscriptions, how the
// provisioner and dispatcher for a bus should be run, and which volumes should
// be mounted into them.
type BusSpec struct {
	// Parameters defines the parameters that must be passed by this Bus'
	// Channels and their Subscriptions. Channels and Subscriptions fulfill
	// these parameters with Arguments.
	Parameters *BusParameters `json:"parameters,omitempty"`

	// Provisioner defines how the provisioner container for this bus should be
	// run. Provisioners are responsible for provisioning the underlying
	// infrastructure for Channels and Subscriptions. The exact work done by the
	// provisioner varies by Bus; one example of work done by a provisioner
	// could be creating a messaging topic that backs a channel.
	Provisioner *kapi.Container `json:"provisioner,omitempty"`

	// Dispatcher defines how the dispatcher container for this bus should be
	// run. Dispatchers are responsible for performing two types of event
	// dispatch: dispatching incoming events to the Bus' Channels and
	// dispatching events in the Channel to the Channel's Subscriptions.
	Dispatcher kapi.Container `json:"dispatcher"`

	// Volumes to be mounted inside the provisioner or dispatcher containers
	Volumes *[]kapi.Volume `json:"volumes,omitempty"`
}

// BusParameters represents the arguments that must be passed by Channels and
// Subscriptions.
type BusParameters struct {

	// Channel configuration params for channels on the bus
	Channel *[]Parameter `json:"channel,omitempty"`

	// Subscription configuration params for subscriptions on the bus
	Subscription *[]Parameter `json:"subscription,omitempty"`
}

// BusStatus (computed) for a bus
type BusStatus struct {
}

func (b *Bus) BacksChannel(channel *Channel) bool {
	return b.Namespace == channel.Namespace && b.Name == channel.Spec.Bus
}

func (b *Bus) GetSpec() *BusSpec {
	return &b.Spec
}

func (b *Bus) GetSpecJSON() ([]byte, error) {
	return json.Marshal(b.Spec)
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// BusList returned in list operations
type BusList struct {
	meta_v1.TypeMeta `json:",inline"`
	meta_v1.ListMeta `json:"metadata"`
	Items            []Bus `json:"items"`
}

// GenericBus may be backed by Bus or ClusterBus
type GenericBus interface {
	runtime.Object
	meta_v1.ObjectMetaAccessor
	BacksChannel(channel *Channel) bool
	GetSpec() *BusSpec
}