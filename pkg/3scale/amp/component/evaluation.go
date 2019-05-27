package component

import (
	appsv1 "github.com/openshift/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type Evaluation struct {
}

func NewEvaluation() *Evaluation {
	return &Evaluation{}
}

func (evaluation *Evaluation) RemoveContainersResourceRequestsAndLimits(objects []runtime.RawExtension) {
	for _, rawExtension := range objects {
		obj := rawExtension.Object
		dc, ok := obj.(*appsv1.DeploymentConfig)
		if ok {
			for containerIdx := range dc.Spec.Template.Spec.Containers {
				container := &dc.Spec.Template.Spec.Containers[containerIdx]
				container.Resources = v1.ResourceRequirements{}
			}
		}
	}
}
