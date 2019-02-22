package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const GroupName = "my.example.com"
const GroupVersion = "v1alpha1"

var (
	SchemeGroupVersion = schema.GroupVersion{Group:GroupName, Version:GroupVersion}
	SchemeBuilder = runtime.NewSchemeBuilder(addToScheme)
	AddToScheme = SchemeBuilder.AddToScheme
)

func addToScheme(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&Redis{},)

	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}