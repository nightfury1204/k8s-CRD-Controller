package util

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AddFinalizer(in metav1.ObjectMeta, finalizer string) metav1.ObjectMeta {
	if !HasFinalizer(in,finalizer) {
		in.Finalizers = append(in.Finalizers,finalizer)
	}
	return in
}

func RemoveFinalizer(in metav1.ObjectMeta, finalizer string) metav1.ObjectMeta {
	temp := in.Finalizers[:0]
	for _,name := range in.Finalizers {
		if name != finalizer {
			temp = append(temp,name)
		}
	}
	in.Finalizers = temp
	return in
}

func RemoveInitializer(in metav1.ObjectMeta) metav1.ObjectMeta {
	if in.GetInitializers()!=nil {
		if len(in.GetInitializers().Pending)==1{
			in.Initializers = nil
		} else {
			in.Initializers.Pending = append(in.Initializers.Pending[:0],in.Initializers.Pending[1:]...)
		}
	}
	return in
}

func HasFinalizer(in metav1.ObjectMeta, finalizer string) bool {
	for _,name := range in.Finalizers {
		if name == finalizer {
			return true
		}
	}
	return false
}