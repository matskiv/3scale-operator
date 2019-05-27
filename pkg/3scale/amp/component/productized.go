package component

import (
	imagev1 "github.com/openshift/api/image/v1"
	templatev1 "github.com/openshift/api/template/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type Productized struct {
	Options *ProductizedOptions
}

func NewProductized(options *ProductizedOptions) *Productized {
	return &Productized{Options: options}
}

func (productized *Productized) UpdateAmpImagesParameters(template *templatev1.Template) {
	for paramIdx := range template.Parameters {
		param := &template.Parameters[paramIdx]
		switch param.Name {
		case "AMP_SYSTEM_IMAGE":
			param.Value = "registry.access.redhat.com/3scale-amp25/system"
		case "AMP_BACKEND_IMAGE":
			param.Value = "registry.access.redhat.com/3scale-amp25/backend"
		case "AMP_APICAST_IMAGE":
			param.Value = "registry.access.redhat.com/3scale-amp25/apicast-gateway"
		case "AMP_ROUTER_IMAGE":
			param.Value = "registry.access.redhat.com/3scale-amp22/wildcard-router"
		case "AMP_ZYNC_IMAGE":
			param.Value = "registry.access.redhat.com/3scale-amp25/zync"
		}
	}
}

func (productized *Productized) UpdateAmpImagesURIs(objects []runtime.RawExtension) []runtime.RawExtension {
	res := objects

	for _, rawExtension := range res {
		obj := rawExtension.Object
		is, ok := obj.(*imagev1.ImageStream)
		if ok {
			for tagIdx := range is.Spec.Tags {
				// Only change the ImageStream tag name that has the ampRelease
				// value. We do not modify the latest tag
				if is.Spec.Tags[tagIdx].Name == productized.Options.ampRelease {
					switch is.Name {
					case "amp-apicast":
						is.Spec.Tags[tagIdx].From.Name = productized.Options.apicastImage
					case "amp-system":
						is.Spec.Tags[tagIdx].From.Name = productized.Options.systemImage
					case "amp-backend":
						is.Spec.Tags[tagIdx].From.Name = productized.Options.backendImage
					case "amp-wildcard-router":
						is.Spec.Tags[tagIdx].From.Name = productized.Options.routerImage
					case "amp-zync":
						is.Spec.Tags[tagIdx].From.Name = productized.Options.zyncImage
					}
				}
			}
		}
	}

	return res
}
