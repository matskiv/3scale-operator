package component

import (
	imagev1 "github.com/openshift/api/image/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	insecureImportPolicy = false
)

type AmpImages struct {
	Options *AmpImagesOptions
}

func NewAmpImages(options *AmpImagesOptions) *AmpImages {
	return &AmpImages{Options: options}
}

func (ampImages *AmpImages) Objects() []runtime.RawExtension {
	backendImageStream := ampImages.buildAmpBackendImageStream()
	zyncImageStream := ampImages.buildAmpZyncImageStream()
	apicastImageStream := ampImages.buildApicastImageStream()
	wildcardRouterImageStream := ampImages.buildWildcardRouterImageStream()
	systemImageStream := ampImages.buildAmpSystemImageStream()
	postgreSQLImageStream := ampImages.buildPostgreSQLImageStream()
	backendRedisImageStream := ampImages.buildBackendRedisImageStream()
	systemRedisImageStream := ampImages.buildSystemRedisImageStream()
	systemMemcachedImageStream := ampImages.buildSystemMemcachedImageStream()
	systemMySQLImageStream := ampImages.buildSystemMySQLImageStream()

	deploymentsServiceAccount := ampImages.buildDeploymentsServiceAccount()

	objects := []runtime.RawExtension{
		runtime.RawExtension{Object: backendImageStream},
		runtime.RawExtension{Object: zyncImageStream},
		runtime.RawExtension{Object: apicastImageStream},
		runtime.RawExtension{Object: wildcardRouterImageStream},
		runtime.RawExtension{Object: systemImageStream},
		runtime.RawExtension{Object: postgreSQLImageStream},
		runtime.RawExtension{Object: backendRedisImageStream},
		runtime.RawExtension{Object: systemRedisImageStream},
		runtime.RawExtension{Object: systemMemcachedImageStream},
		runtime.RawExtension{Object: systemMySQLImageStream},
		runtime.RawExtension{Object: deploymentsServiceAccount},
	}
	return objects
}

func (ampImages *AmpImages) buildAmpBackendImageStream() *imagev1.ImageStream {
	return &imagev1.ImageStream{
		ObjectMeta: metav1.ObjectMeta{
			Name: "amp-backend",
			Labels: map[string]string{
				"app":                  ampImages.Options.appLabel,
				"threescale_component": "backend",
			},
			Annotations: map[string]string{
				"openshift.io/display-name": "AMP backend",
			},
		},
		TypeMeta: metav1.TypeMeta{APIVersion: "image.openshift.io/v1", Kind: "ImageStream"},
		Spec: imagev1.ImageStreamSpec{
			Tags: []imagev1.TagReference{
				imagev1.TagReference{
					Name: "latest",
					Annotations: map[string]string{
						"openshift.io/display-name": "amp-backend (latest)",
					},
					From: &v1.ObjectReference{
						Kind: "ImageStreamTag",
						Name: ampImages.Options.ampRelease,
					},
				},
				imagev1.TagReference{
					Name: ampImages.Options.ampRelease,
					Annotations: map[string]string{
						"openshift.io/display-name": "amp-backend " + ampImages.Options.ampRelease,
					},
					From: &v1.ObjectReference{
						Kind: "DockerImage",
						Name: ampImages.Options.backendImage,
					},
					ImportPolicy: imagev1.TagImportPolicy{
						// TODO this was originally a double brace expansion from a variable, that is not possible
						// natively with kubernetes so we replaced it with a const
						Insecure: insecureImportPolicy,
					},
				},
			},
		},
	}
}

func (ampImages *AmpImages) buildAmpZyncImageStream() *imagev1.ImageStream {
	return &imagev1.ImageStream{
		ObjectMeta: metav1.ObjectMeta{
			Name: "amp-zync",
			Labels: map[string]string{
				"app":                  ampImages.Options.appLabel,
				"threescale_component": "zync",
			},
			Annotations: map[string]string{
				"openshift.io/display-name": "AMP Zync",
			},
		},
		TypeMeta: metav1.TypeMeta{APIVersion: "image.openshift.io/v1", Kind: "ImageStream"},
		Spec: imagev1.ImageStreamSpec{
			Tags: []imagev1.TagReference{
				imagev1.TagReference{
					Name: "latest",
					Annotations: map[string]string{
						"openshift.io/display-name": "AMP Zync (latest)",
					},
					From: &v1.ObjectReference{
						Kind: "ImageStreamTag",
						Name: ampImages.Options.ampRelease,
					},
				},
				imagev1.TagReference{
					Name: ampImages.Options.ampRelease,
					Annotations: map[string]string{
						"openshift.io/display-name": "AMP Zync " + ampImages.Options.ampRelease,
					},
					From: &v1.ObjectReference{
						Kind: "DockerImage",
						Name: ampImages.Options.zyncImage,
					},
					ImportPolicy: imagev1.TagImportPolicy{
						// TODO this was originally a double brace expansion from a variable, that is not possible
						// natively with kubernetes so we replaced it with a const
						Insecure: insecureImportPolicy,
					},
				},
			},
		},
	}
}

func (ampImages *AmpImages) buildApicastImageStream() *imagev1.ImageStream {
	return &imagev1.ImageStream{
		ObjectMeta: metav1.ObjectMeta{
			Name: "amp-apicast",
			Labels: map[string]string{
				"app":                  ampImages.Options.appLabel,
				"threescale_component": "apicast",
			},
			Annotations: map[string]string{
				"openshift.io/display-name": "AMP APIcast",
			},
		},
		TypeMeta: metav1.TypeMeta{APIVersion: "image.openshift.io/v1", Kind: "ImageStream"},
		Spec: imagev1.ImageStreamSpec{
			Tags: []imagev1.TagReference{
				imagev1.TagReference{
					Name: "latest",
					Annotations: map[string]string{
						"openshift.io/display-name": "AMP APIcast (latest)",
					},
					From: &v1.ObjectReference{
						Kind: "ImageStreamTag",
						Name: ampImages.Options.ampRelease,
					},
				},
				imagev1.TagReference{
					Name: ampImages.Options.ampRelease,
					Annotations: map[string]string{
						"openshift.io/display-name": "AMP APIcast " + ampImages.Options.ampRelease,
					},
					From: &v1.ObjectReference{
						Kind: "DockerImage",
						Name: ampImages.Options.apicastImage,
					},
					ImportPolicy: imagev1.TagImportPolicy{
						// TODO this was originally a double brace expansion from a variable, that is not possible
						// natively with kubernetes so we replaced it with a const
						Insecure: insecureImportPolicy,
					},
				},
			},
		},
	}
}

func (ampImages *AmpImages) buildWildcardRouterImageStream() *imagev1.ImageStream {
	return &imagev1.ImageStream{
		ObjectMeta: metav1.ObjectMeta{
			Name: "amp-wildcard-router",
			Labels: map[string]string{
				"app":                  ampImages.Options.appLabel,
				"threescale_component": "wildcard-router",
			},
			Annotations: map[string]string{
				"openshift.io/display-name": "AMP APIcast Wildcard Router",
			},
		},
		TypeMeta: metav1.TypeMeta{APIVersion: "image.openshift.io/v1", Kind: "ImageStream"},
		Spec: imagev1.ImageStreamSpec{
			Tags: []imagev1.TagReference{
				imagev1.TagReference{
					Name: "latest",
					Annotations: map[string]string{
						"openshift.io/display-name": "AMP APIcast Wildcard Router (latest)",
					},
					From: &v1.ObjectReference{
						Kind: "ImageStreamTag",
						Name: ampImages.Options.ampRelease,
					},
				},
				imagev1.TagReference{
					Name: ampImages.Options.ampRelease,
					Annotations: map[string]string{
						"openshift.io/display-name": "AMP APIcast Wildcard Router " + ampImages.Options.ampRelease,
					},
					From: &v1.ObjectReference{
						Kind: "DockerImage",
						Name: ampImages.Options.routerImage,
					},
					ImportPolicy: imagev1.TagImportPolicy{
						// TODO this was originally a double brace expansion from a variable, that is not possible
						// natively with kubernetes so we replaced it with a const
						Insecure: insecureImportPolicy,
					},
				},
			},
		},
	}
}

func (ampImages *AmpImages) buildAmpSystemImageStream() *imagev1.ImageStream {
	return &imagev1.ImageStream{
		ObjectMeta: metav1.ObjectMeta{
			Name: "amp-system",
			Labels: map[string]string{
				"app":                  ampImages.Options.appLabel,
				"threescale_component": "system",
			},
			Annotations: map[string]string{
				"openshift.io/display-name": "AMP System",
			},
		},
		TypeMeta: metav1.TypeMeta{APIVersion: "image.openshift.io/v1", Kind: "ImageStream"},
		Spec: imagev1.ImageStreamSpec{
			Tags: []imagev1.TagReference{
				imagev1.TagReference{
					Name: "latest",
					Annotations: map[string]string{
						"openshift.io/display-name": "AMP System (latest)",
					},
					From: &v1.ObjectReference{
						Kind: "ImageStreamTag",
						Name: ampImages.Options.ampRelease,
					},
				},
				imagev1.TagReference{
					Name: ampImages.Options.ampRelease,
					Annotations: map[string]string{
						"openshift.io/display-name": "AMP system " + ampImages.Options.ampRelease,
					},
					From: &v1.ObjectReference{
						Kind: "DockerImage",
						Name: ampImages.Options.systemImage,
					},
					ImportPolicy: imagev1.TagImportPolicy{
						Insecure: insecureImportPolicy,
					},
				},
			},
		},
	}
}

func (ampImages *AmpImages) buildPostgreSQLImageStream() *imagev1.ImageStream {
	return &imagev1.ImageStream{
		ObjectMeta: metav1.ObjectMeta{
			Name: "postgresql",
			Labels: map[string]string{
				"app":                  ampImages.Options.appLabel,
				"threescale_component": "system",
			},
			Annotations: map[string]string{
				"openshift.io/display-name": "Zync database",
			},
		},
		TypeMeta: metav1.TypeMeta{APIVersion: "image.openshift.io/v1", Kind: "ImageStream"},
		Spec: imagev1.ImageStreamSpec{
			Tags: []imagev1.TagReference{
				imagev1.TagReference{
					Name: "latest",
					Annotations: map[string]string{
						"openshift.io/display-name": "Zync PostgreSQL (latest)",
					},
					From: &v1.ObjectReference{
						Kind: "ImageStreamTag",
						Name: ampImages.Options.ampRelease,
					},
				},
				imagev1.TagReference{
					Name: ampImages.Options.ampRelease,
					Annotations: map[string]string{
						"openshift.io/display-name": "Zync " + ampImages.Options.ampRelease + " PostgreSQL",
					},
					From: &v1.ObjectReference{
						Kind: "DockerImage",
						Name: ampImages.Options.postgreSQLImage,
					},
					ImportPolicy: imagev1.TagImportPolicy{
						Insecure: insecureImportPolicy,
					},
				},
			},
		},
	}
}

func (ampImages *AmpImages) buildBackendRedisImageStream() *imagev1.ImageStream {
	return &imagev1.ImageStream{
		ObjectMeta: metav1.ObjectMeta{
			Name: "backend-redis",
			Labels: map[string]string{
				"app":                  ampImages.Options.appLabel,
				"threescale_component": "backend",
			},
			Annotations: map[string]string{
				"openshift.io/display-name": "Backend Redis",
			},
		},
		TypeMeta: metav1.TypeMeta{APIVersion: "image.openshift.io/v1", Kind: "ImageStream"},
		Spec: imagev1.ImageStreamSpec{
			Tags: []imagev1.TagReference{
				imagev1.TagReference{
					Name: "latest",
					Annotations: map[string]string{
						"openshift.io/display-name": "Backend Redis (latest)",
					},
					From: &v1.ObjectReference{
						Kind: "ImageStreamTag",
						Name: ampImages.Options.ampRelease,
					},
				},
				imagev1.TagReference{
					Name: ampImages.Options.ampRelease,
					Annotations: map[string]string{
						"openshift.io/display-name": "Backend " + ampImages.Options.ampRelease + " Redis",
					},
					From: &v1.ObjectReference{
						Kind: "DockerImage",
						Name: ampImages.Options.backendRedisImage,
					},
					ImportPolicy: imagev1.TagImportPolicy{
						Insecure: insecureImportPolicy,
					},
				},
			},
		},
	}
}

func (ampImages *AmpImages) buildSystemRedisImageStream() *imagev1.ImageStream {
	return &imagev1.ImageStream{
		ObjectMeta: metav1.ObjectMeta{
			Name: "system-redis",
			Labels: map[string]string{
				"app":                  ampImages.Options.appLabel,
				"threescale_component": "system",
			},
			Annotations: map[string]string{
				"openshift.io/display-name": "System Redis",
			},
		},
		TypeMeta: metav1.TypeMeta{APIVersion: "image.openshift.io/v1", Kind: "ImageStream"},
		Spec: imagev1.ImageStreamSpec{
			Tags: []imagev1.TagReference{
				imagev1.TagReference{
					Name: "latest",
					Annotations: map[string]string{
						"openshift.io/display-name": "System Redis (latest)",
					},
					From: &v1.ObjectReference{
						Kind: "ImageStreamTag",
						Name: ampImages.Options.ampRelease,
					},
				},
				imagev1.TagReference{
					Name: ampImages.Options.ampRelease,
					Annotations: map[string]string{
						"openshift.io/display-name": "System " + ampImages.Options.ampRelease + " Redis",
					},
					From: &v1.ObjectReference{
						Kind: "DockerImage",
						Name: ampImages.Options.backendRedisImage,
					},
					ImportPolicy: imagev1.TagImportPolicy{
						Insecure: insecureImportPolicy,
					},
				},
			},
		},
	}
}

func (ampImages *AmpImages) buildSystemMemcachedImageStream() *imagev1.ImageStream {
	return &imagev1.ImageStream{
		ObjectMeta: metav1.ObjectMeta{
			Name: "system-memcached",
			Labels: map[string]string{
				"app":                  ampImages.Options.appLabel,
				"threescale_component": "system",
			},
			Annotations: map[string]string{
				"openshift.io/display-name": "System Memcached",
			},
		},
		TypeMeta: metav1.TypeMeta{APIVersion: "image.openshift.io/v1", Kind: "ImageStream"},
		Spec: imagev1.ImageStreamSpec{
			Tags: []imagev1.TagReference{
				imagev1.TagReference{
					Name: "latest",
					Annotations: map[string]string{
						"openshift.io/display-name": "System Memcached (latest)",
					},
					From: &v1.ObjectReference{
						Kind: "ImageStreamTag",
						Name: ampImages.Options.ampRelease,
					},
				},
				imagev1.TagReference{
					Name: ampImages.Options.ampRelease,
					Annotations: map[string]string{
						"openshift.io/display-name": "System " + ampImages.Options.ampRelease + " Memcached",
					},
					From: &v1.ObjectReference{
						Kind: "DockerImage",
						Name: ampImages.Options.systemMemcachedImage,
					},
					ImportPolicy: imagev1.TagImportPolicy{
						Insecure: insecureImportPolicy,
					},
				},
			},
		},
	}
}

func (ampImages *AmpImages) buildSystemMySQLImageStream() *imagev1.ImageStream {
	return &imagev1.ImageStream{
		ObjectMeta: metav1.ObjectMeta{
			Name: "system-mysql",
			Labels: map[string]string{
				"app":                  ampImages.Options.appLabel,
				"threescale_component": "system",
			},
			Annotations: map[string]string{
				"openshift.io/display-name": "System MySQL",
			},
		},
		TypeMeta: metav1.TypeMeta{APIVersion: "image.openshift.io/v1", Kind: "ImageStream"},
		Spec: imagev1.ImageStreamSpec{
			Tags: []imagev1.TagReference{
				imagev1.TagReference{
					Name: "latest",
					Annotations: map[string]string{
						"openshift.io/display-name": "System MySQL (latest)",
					},
					From: &v1.ObjectReference{
						Kind: "ImageStreamTag",
						Name: ampImages.Options.ampRelease,
					},
				},
				imagev1.TagReference{
					Name: ampImages.Options.ampRelease,
					Annotations: map[string]string{
						"openshift.io/display-name": "System " + ampImages.Options.ampRelease + " MySQL",
					},
					From: &v1.ObjectReference{
						Kind: "DockerImage",
						Name: ampImages.Options.systemMySQLImage,
					},
					ImportPolicy: imagev1.TagImportPolicy{
						Insecure: insecureImportPolicy,
					},
				},
			},
		},
	}
}

func (ampImages *AmpImages) buildDeploymentsServiceAccount() *v1.ServiceAccount {
	return &v1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ServiceAccount",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "amp",
		},
		ImagePullSecrets: []v1.LocalObjectReference{
			v1.LocalObjectReference{
				Name: "threescale-registry-auth"}}}
}
