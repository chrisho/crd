package v1alpha1

import "k8s.io/apimachinery/pkg/runtime"

func (in *Redis) DeepCopyObject() runtime.Object {
	out := &Redis{}
	in.DeepCopyInto(out)

	return out
}

func (in *Redis) DeepCopyInto(out *Redis){
	out.TypeMeta = in.TypeMeta
	out.ObjectMeta = in.ObjectMeta
	out.Spec = RedisSpec{
		Replicas: in.Spec.Replicas,
	}
}

func (in *RedisList) DeepCopyObject() runtime.Object {
	out := &RedisList{}
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		out.Items = make([]Redis, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&out.Items[i])
		}
	}

	return out
}
