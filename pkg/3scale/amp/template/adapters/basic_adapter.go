package adapters

import (
	templatev1 "github.com/openshift/api/template/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type AppenderElement interface {
	Parameters() []templatev1.Parameter
	Objects() ([]runtime.RawExtension, error)
}

type AppenderAdapter struct {
	AppenderElement AppenderElement
}

func NewAppenderAdapter(s AppenderElement) Adapter {
	return &AppenderAdapter{AppenderElement: s}
}

func (b *AppenderAdapter) Adapt(template *templatev1.Template) {
	parameters := b.AppenderElement.Parameters()
	template.Parameters = append(template.Parameters, parameters...)
	objects, err := b.AppenderElement.Objects()
	if err != nil {
		panic(err)
	}
	template.Objects = append(template.Objects, objects...)
}
