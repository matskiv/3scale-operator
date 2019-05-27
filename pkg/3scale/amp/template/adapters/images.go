package adapters

import (
	"github.com/3scale/3scale-operator/pkg/3scale/amp/component"
	templatev1 "github.com/openshift/api/template/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type ImagesAdapter struct {
}

func NewImagesAdapter(options []string) Adapter {
	return NewAppenderAdapter(&ImagesAdapter{})
}

func (i *ImagesAdapter) Parameters() []templatev1.Parameter {
	return []templatev1.Parameter{
		templatev1.Parameter{
			Name:     "AMP_BACKEND_IMAGE",
			Required: true,
			Value:    "quay.io/3scale/apisonator:nightly",
		},
		templatev1.Parameter{
			Name:     "AMP_ZYNC_IMAGE",
			Value:    "quay.io/3scale/zync:nightly",
			Required: true,
		},
		templatev1.Parameter{
			Name:     "AMP_APICAST_IMAGE",
			Value:    "quay.io/3scale/apicast:nightly",
			Required: true,
		},
		templatev1.Parameter{
			Name:     "AMP_ROUTER_IMAGE",
			Value:    "quay.io/3scale/wildcard-router:nightly",
			Required: true,
		},
		templatev1.Parameter{
			Name:     "AMP_SYSTEM_IMAGE",
			Value:    "quay.io/3scale/porta:nightly",
			Required: true,
		},
		templatev1.Parameter{
			Name:        "POSTGRESQL_IMAGE",
			Description: "Postgresql image to use",
			Value:       "registry.access.redhat.com/rhscl/postgresql-10-rhel7",
			Required:    true,
		},
		templatev1.Parameter{
			Name:        "MYSQL_IMAGE",
			Description: "Mysql image to use",
			Value:       "registry.access.redhat.com/rhscl/mysql-57-rhel7:5.7",
			Required:    true,
		},
		templatev1.Parameter{
			Name:        "MEMCACHED_IMAGE",
			Description: "Memcached image to use",
			Value:       "registry.access.redhat.com/3scale-amp20/memcached",
			Required:    true,
		},
		templatev1.Parameter{
			Name:        "IMAGESTREAM_TAG_IMPORT_INSECURE",
			Description: "Set to true if the server may bypass certificate verification or connect directly over HTTP during image import.",
			Value:       "false",
			Required:    true,
		},
	}
}

func (i *ImagesAdapter) Objects() ([]runtime.RawExtension, error) {
	imagesOptions, err := i.options()
	if err != nil {
		return nil, err
	}
	imagesComponent := component.NewAmpImages(imagesOptions)
	return imagesComponent.Objects(), nil
}

func (i *ImagesAdapter) options() (*component.AmpImagesOptions, error) {
	aob := component.AmpImagesOptionsBuilder{}
	aob.AppLabel("${APP_LABEL}")
	aob.AMPRelease("${AMP_RELEASE}")
	aob.ApicastImage("${AMP_APICAST_IMAGE}")
	aob.BackendImage("${AMP_BACKEND_IMAGE}")
	aob.RouterImage("${AMP_ROUTER_IMAGE}")
	aob.SystemImage("${AMP_SYSTEM_IMAGE}")
	aob.ZyncImage("${AMP_ZYNC_IMAGE}")
	aob.PostgreSQLImage("${POSTGRESQL_IMAGE}")
	aob.BackendRedisImage("${REDIS_IMAGE}")
	aob.SystemRedisImage("${REDIS_IMAGE}")
	aob.SystemMemcachedImage("${MEMCACHED_IMAGE}")
	aob.SystemMySQLImage("${MYSQL_IMAGE}")
	aob.InsecureImportPolicy(false)
	return aob.Build()
}
