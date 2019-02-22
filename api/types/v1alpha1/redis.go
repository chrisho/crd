package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type RedisSpec struct {
	DeploymentName string `json:"deploymentName"`
	Replicas *int32 `json:"replicas"`
}

type Redis struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec RedisSpec `json:"spec"`
}

type RedisList struct {
	metav1.TypeMeta		`json:",inline"`
	metav1.ListMeta		`json:"metadata,omitempty"`

	Items []Redis		`json:"items"`
}