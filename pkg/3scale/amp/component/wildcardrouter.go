package component

import (
	appsv1 "github.com/openshift/api/apps/v1"
	routev1 "github.com/openshift/api/route/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type WildcardRouter struct {
	Options *WildcardRouterOptions
}

func NewWildcardRouter(options *WildcardRouterOptions) *WildcardRouter {
	return &WildcardRouter{Options: options}
}

func (wr *WildcardRouter) Objects() []runtime.RawExtension {
	wildcardRouterDeploymentConfig := wr.buildWildcardRouterDeploymentConfig()
	wildcardRouterService := wr.buildWildcardRouterService()
	wildcardRouterRoute := wr.buildWildcardRouterRoute()

	objects := []runtime.RawExtension{
		runtime.RawExtension{Object: wildcardRouterDeploymentConfig},
		runtime.RawExtension{Object: wildcardRouterService},
		runtime.RawExtension{Object: wildcardRouterRoute},
	}

	return objects
}

func (wr *WildcardRouter) buildWildcardRouterRoute() *routev1.Route {
	return &routev1.Route{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Route",
			APIVersion: "route.openshift.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   "apicast-wildcard-router",
			Labels: map[string]string{"app": wr.Options.appLabel, "threescale_component": "apicast", "threescale_component_element": "wildcard-router"},
		},
		Spec: routev1.RouteSpec{
			Host: "apicast-wildcard." + wr.Options.wildcardDomain,
			To: routev1.RouteTargetReference{
				Kind: "Service",
				Name: "apicast-wildcard-router",
			},
			Port: &routev1.RoutePort{
				TargetPort: intstr.FromString("http"),
			},
			WildcardPolicy: routev1.WildcardPolicyType(wr.Options.wildcardPolicy),
			TLS: &routev1.TLSConfig{
				Termination:                   routev1.TLSTerminationEdge,
				InsecureEdgeTerminationPolicy: routev1.InsecureEdgeTerminationPolicyAllow},
		},
	}
}

func (wr *WildcardRouter) buildWildcardRouterService() *v1.Service {
	return &v1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "apicast-wildcard-router",
			Labels: map[string]string{
				"app":                          wr.Options.appLabel,
				"threescale_component":         "apicast",
				"threescale_component_element": "wildcard-router",
			},
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				v1.ServicePort{
					Name:       "http",
					Protocol:   v1.ProtocolTCP,
					Port:       8080,
					TargetPort: intstr.FromString("http"),
				},
			},
			Selector: map[string]string{"deploymentConfig": "apicast-wildcard-router"},
		},
	}
}

func (wr *WildcardRouter) buildWildcardRouterDeploymentConfig() *appsv1.DeploymentConfig {
	return &appsv1.DeploymentConfig{
		TypeMeta: metav1.TypeMeta{APIVersion: "apps.openshift.io/v1", Kind: "DeploymentConfig"},
		ObjectMeta: metav1.ObjectMeta{
			Name: "apicast-wildcard-router",
			Labels: map[string]string{
				"app":                          wr.Options.appLabel,
				"threescale_component":         "apicast",
				"threescale_component_element": "wildcard-router",
			},
		},
		Spec: appsv1.DeploymentConfigSpec{
			Replicas: 1,
			Selector: map[string]string{
				"deploymentConfig": "apicast-wildcard-router",
			},
			Strategy: appsv1.DeploymentStrategy{
				RollingParams: &appsv1.RollingDeploymentStrategyParams{
					IntervalSeconds: &[]int64{1}[0],
					MaxSurge: &intstr.IntOrString{
						Type:   intstr.Type(intstr.String),
						StrVal: "25%",
					},
					MaxUnavailable: &intstr.IntOrString{
						Type:   intstr.Type(intstr.String),
						StrVal: "25%",
					},
					TimeoutSeconds:      &[]int64{1800}[0],
					UpdatePeriodSeconds: &[]int64{1}[0],
				},
				Type: appsv1.DeploymentStrategyTypeRolling,
			},
			Triggers: appsv1.DeploymentTriggerPolicies{
				appsv1.DeploymentTriggerPolicy{
					Type: appsv1.DeploymentTriggerOnConfigChange,
				},
				appsv1.DeploymentTriggerPolicy{
					Type: appsv1.DeploymentTriggerOnImageChange,
					ImageChangeParams: &appsv1.DeploymentTriggerImageChangeParams{
						Automatic: true,
						ContainerNames: []string{
							"apicast-wildcard-router",
						},
						From: v1.ObjectReference{
							Kind: "ImageStreamTag",
							Name: "amp-wildcard-router:latest",
						},
					},
				},
			},
			Template: &v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"deploymentConfig":             "apicast-wildcard-router",
						"app":                          wr.Options.appLabel,
						"threescale_component":         "apicast",
						"threescale_component_element": "wildcard-router",
					},
				},
				Spec: v1.PodSpec{
					ServiceAccountName: "amp",
					Containers: []v1.Container{
						v1.Container{
							Ports: []v1.ContainerPort{
								v1.ContainerPort{
									ContainerPort: 8080,
									Protocol:      v1.ProtocolTCP,
									Name:          "http",
								},
							},
							Env:             wr.buildWildcardRouterEnv(),
							Image:           "amp-wildcard-router:latest",
							ImagePullPolicy: v1.PullIfNotPresent,
							Name:            "apicast-wildcard-router",
							Resources: v1.ResourceRequirements{
								Limits: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("500m"),
									v1.ResourceMemory: resource.MustParse("64Mi"),
								},
								Requests: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("120m"),
									v1.ResourceMemory: resource.MustParse("32Mi"),
								},
							},
							LivenessProbe: &v1.Probe{
								Handler: v1.Handler{TCPSocket: &v1.TCPSocketAction{
									Port: intstr.FromString("http"),
								}},
								InitialDelaySeconds: 30,
								PeriodSeconds:       10,
							},
						},
					},
				},
			},
		},
	}
}

func (wr *WildcardRouter) buildWildcardRouterEnv() []v1.EnvVar {
	return []v1.EnvVar{
		envVarFromSecret("API_HOST", "system-master-apicast", "BASE_URL"),
	}
}
