package adapters

import (
	"github.com/3scale/3scale-operator/pkg/3scale/amp/component"
	templatev1 "github.com/openshift/api/template/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type Zync struct {
}

func NewZyncAdapter(options []string) Adapter {
	return NewAppenderAdapter(&Zync{})
}

func (z *Zync) Parameters() []templatev1.Parameter {
	return []templatev1.Parameter{
		templatev1.Parameter{
			Name:        "ZYNC_DATABASE_PASSWORD",
			DisplayName: "PostgreSQL Connection Password",
			Description: "Password for the PostgreSQL connection user.",
			Generate:    "expression",
			From:        "[a-zA-Z0-9]{16}",
			Required:    true,
		},
		templatev1.Parameter{
			Name:     "ZYNC_SECRET_KEY_BASE",
			Generate: "expression",
			From:     "[a-zA-Z0-9]{16}",
			Required: true,
		},
		templatev1.Parameter{
			Name:     "ZYNC_AUTHENTICATION_TOKEN",
			Generate: "expression",
			From:     "[a-zA-Z0-9]{16}",
			Required: true,
		},
	}
}

func (z *Zync) Objects() ([]runtime.RawExtension, error) {
	zyncOptions, err := z.options()
	if err != nil {
		return nil, err
	}
	zyncComponent := component.NewZync(zyncOptions)
	return zyncComponent.Objects(), nil
}

func (z *Zync) options() (*component.ZyncOptions, error) {
	zob := component.ZyncOptionsBuilder{}
	zob.AppLabel("${APP_LABEL}")
	zob.AuthenticationToken("${ZYNC_AUTHENTICATION_TOKEN}")
	zob.DatabasePassword("${ZYNC_DATABASE_PASSWORD}")
	zob.SecretKeyBase("${ZYNC_SECRET_KEY_BASE}")
	return zob.Build()
}
