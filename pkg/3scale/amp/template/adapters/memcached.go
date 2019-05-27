package adapters

import (
	"github.com/3scale/3scale-operator/pkg/3scale/amp/component"
	templatev1 "github.com/openshift/api/template/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type MemcachedAdapter struct {
}

func NewMemcachedAdapter(options []string) Adapter {
	return NewAppenderAdapter(&MemcachedAdapter{})
}

func (m *MemcachedAdapter) Parameters() []templatev1.Parameter {
	return []templatev1.Parameter{}
}

func (m *MemcachedAdapter) Objects() ([]runtime.RawExtension, error) {
	memcachedOptions, err := m.options()
	if err != nil {
		return nil, err
	}
	memcachedComponent := component.NewMemcached(memcachedOptions)
	return memcachedComponent.Objects(), nil
}

func (m *MemcachedAdapter) options() (*component.MemcachedOptions, error) {
	rob := component.MemcachedOptionsBuilder{}
	rob.AppLabel("${APP_LABEL}")
	return rob.Build()
}
