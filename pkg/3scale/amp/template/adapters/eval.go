package adapters

import (
	"github.com/3scale/3scale-operator/pkg/3scale/amp/component"
	templatev1 "github.com/openshift/api/template/v1"
)

type EvalAdapter struct {
}

func NewEvalAdapter(options []string) Adapter {
	return &EvalAdapter{}
}

func (e *EvalAdapter) Adapt(template *templatev1.Template) {
	// update metadata
	template.Name = "3scale-api-management-eval"
	template.ObjectMeta.Annotations["description"] = "3scale API Management main system (Evaluation)"

	evalComponent := component.NewEvaluation()
	evalComponent.RemoveContainersResourceRequestsAndLimits(template.Objects)
}
