package adapters

import (
	"github.com/3scale/3scale-operator/pkg/3scale/amp/component"
	templatev1 "github.com/openshift/api/template/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type RedisAdapter struct {
}

func NewRedisAdapter(options []string) Adapter {
	return NewAppenderAdapter(&RedisAdapter{})
}

func (a *RedisAdapter) Parameters() []templatev1.Parameter {
	return []templatev1.Parameter{
		{
			Name:        "REDIS_IMAGE",
			Description: "Redis image to use",
			Required:    true,
			Value:       "registry.access.redhat.com/rhscl/redis-32-rhel7:3.2",
		},
	}
}

func (r *RedisAdapter) Objects() ([]runtime.RawExtension, error) {
	redisOptions, err := r.options()
	if err != nil {
		return nil, err
	}
	redisComponent := component.NewRedis(redisOptions)
	return redisComponent.Objects(), nil
}

func (r *RedisAdapter) options() (*component.RedisOptions, error) {
	rob := component.RedisOptionsBuilder{}
	rob.AppLabel("${APP_LABEL}")

	return rob.Build()
}
