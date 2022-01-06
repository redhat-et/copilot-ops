/*
Copyright 2022.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// CompletionSpec defines the desired state of Completion
type CompletionSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Completion. Edit completion_types.go to remove/update
	UserPrompt string `json:"userPrompt"`
	//+kubebuilder:validation:Minimum=0
	//+kubebuilder:validation:Maximum=4096
	MaxTokens int `json:"maxTokens"`
	// Temperature float64 `json:"temperature"`
	// TopP        float64 `json:"topP"`
}

// CompletionStatus defines the observed state of Completion
type CompletionStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Completion         string `json:"completion"`
	ObservedGeneration int64  `json:"observedGeneration,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Completion is the Schema for the completions API
//+kubebuilder:subresource:status
type Completion struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CompletionSpec   `json:"spec,omitempty"`
	Status CompletionStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CompletionList contains a list of Completion
type CompletionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Completion `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Completion{}, &CompletionList{})
}

func (c *Completion) NeedsReconcile() bool {
	return c.Status.ObservedGeneration != c.GetObjectMeta().GetGeneration()
}
