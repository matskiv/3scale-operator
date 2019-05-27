package adapters

import (
	"github.com/3scale/3scale-operator/pkg/3scale/amp/component"
	templatev1 "github.com/openshift/api/template/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type WildcardRouter struct {
}

func NewWildcardRouterAdapter(options []string) Adapter {
	return NewAppenderAdapter(&WildcardRouter{})
}

func (w *WildcardRouter) Parameters() []templatev1.Parameter {
	return []templatev1.Parameter{
		templatev1.Parameter{
			Name:        "WILDCARD_DOMAIN",
			Description: "Root domain for the wildcard routes. Eg. example.com will generate 3scale-admin.example.com.",
			Required:    true,
		},
		templatev1.Parameter{
			Name:        "WILDCARD_POLICY",
			Description: "Use \"Subdomain\" to create a wildcard route for apicast wildcard router",
			Value:       "None",
			Required:    true,
		},
	}
}

func (w *WildcardRouter) Objects() ([]runtime.RawExtension, error) {
	wildcardRouterOptions, err := w.options()
	if err != nil {
		return nil, err
	}
	wildcardrouterComponent := component.NewWildcardRouter(wildcardRouterOptions)
	return wildcardrouterComponent.Objects(), nil
}

func (w *WildcardRouter) options() (*component.WildcardRouterOptions, error) {
	wrob := component.WildcardRouterOptionsBuilder{}
	wrob.AppLabel("${APP_LABEL}")
	wrob.WildcardDomain("${WILDCARD_DOMAIN}")
	wrob.WildcardPolicy("${WILDCARD_POLICY}")
	return wrob.Build()
}
