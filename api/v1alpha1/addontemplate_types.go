/*
Copyright 2021.

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

// AddonTemplateSpec defines the desired state of AddonTemplate
// View markers: $ controller-gen -www crd
// TODO add missing fields from schema, simplified type for POC
type AddonTemplateSpec struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^[A-Za-z0-9][A-Za-z0-9-]{0,30}[A-Za-z0-9]$`
	ID string `json:"id"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^[0-9A-Z][A-Za-z0-9-_ ()]+$`
	Name string `json:"name"`

	// +kubebuilder:validation:Required
	Description string `json:"description"`

	// +optional
	// +kubebuilder:validation:Pattern=`^http[s]?://(?:[a-zA-Z]|[0-9]|[$-_@.&+]|[!*\(\),]|(?:%[0-9a-fA-F][0-9a-fA-F]))+$`
	Link string `json:"link"`

	// +kubebuilder:validation:Required
	// TODO base64
	Icon string `json:"icon"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^api\.openshift\.com/addon-[0-9a-z][0-9a-z-]{0,30}[0-9a-z]$`
	Label string `json:"label"`

	// +kubebuilder:validation:Required
	Enabled bool `json:"enabled"`
}

// AddonTemplateStatus defines the observed state of AddonTemplate
type AddonTemplateStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// AddonTemplate is the Schema for the addontemplates API
type AddonTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AddonTemplateSpec   `json:"spec,omitempty"`
	Status AddonTemplateStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AddonTemplateList contains a list of AddonTemplate
type AddonTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AddonTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AddonTemplate{}, &AddonTemplateList{})
}
