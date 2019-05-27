package component

import (
	appsv1 "github.com/openshift/api/apps/v1"
	imagev1 "github.com/openshift/api/image/v1"
	templatev1 "github.com/openshift/api/template/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type HighAvailability struct {
	Options *HighAvailabilityOptions
}

const (
	HighlyAvailableReplicas = 2
)

var highlyAvailableExternalDatabases = map[string]bool{
	"backend-redis": true,
	"system-redis":  true,
	"system-mysql":  true,
}

func NewHighAvailability(options *HighAvailabilityOptions) *HighAvailability {
	return &HighAvailability{Options: options}
}

func (ha *HighAvailability) Objects() []runtime.RawExtension {
	systemDatabaseSecrets := ha.createSystemDatabaseSecret()

	objects := []runtime.RawExtension{
		runtime.RawExtension{Object: systemDatabaseSecrets},
	}
	return objects
}

func (ha *HighAvailability) createSystemDatabaseSecret() *v1.Secret {
	return &v1.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Secret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: SystemSecretSystemDatabaseSecretName,
			Labels: map[string]string{
				"app":                  ha.Options.appLabel,
				"threescale_component": "system",
			},
		},
		StringData: map[string]string{
			SystemSecretSystemDatabaseURLFieldName: ha.Options.systemDatabaseURL,
		},
		Type: v1.SecretTypeOpaque,
	}
}

func (ha *HighAvailability) IncreaseReplicasNumber(objects []runtime.RawExtension) {
	// We do not increase the number of replicas in database DeploymentConfigs
	excludedDeploymentConfigs := map[string]bool{
		"system-memcache": true,
		"system-sphinx":   true,
		"zync-database":   true,
	}

	for _, rawExtension := range objects {
		obj := rawExtension.Object
		dc, ok := obj.(*appsv1.DeploymentConfig)
		if ok {
			if _, isExcluded := excludedDeploymentConfigs[dc.Name]; !isExcluded {
				dc.Spec.Replicas = HighlyAvailableReplicas
			}
		}
	}
}

func (ha *HighAvailability) DeleteInternalDatabasesObjects(objects []runtime.RawExtension) []runtime.RawExtension {
	keepObjects := []runtime.RawExtension{}

	for rawExtIdx, rawExtension := range objects {
		switch obj := (rawExtension.Object).(type) {
		case *appsv1.DeploymentConfig:
			if _, ok := highlyAvailableExternalDatabases[obj.ObjectMeta.Name]; !ok {
				// We create a new array and add to it the elements that will
				//NOT have to be deleted
				keepObjects = append(keepObjects, objects[rawExtIdx])
			}
		case *v1.Service:
			if _, ok := highlyAvailableExternalDatabases[obj.ObjectMeta.Name]; !ok {
				keepObjects = append(keepObjects, objects[rawExtIdx])
			}
		case *v1.PersistentVolumeClaim:
			if obj.ObjectMeta.Name != "backend-redis-storage" && obj.ObjectMeta.Name != "system-redis-storage" &&
				obj.ObjectMeta.Name != "mysql-storage" {
				keepObjects = append(keepObjects, objects[rawExtIdx])
			}
		case *v1.ConfigMap:
			if obj.ObjectMeta.Name != "mysql-main-conf" && obj.ObjectMeta.Name != "mysql-extra-conf" &&
				obj.ObjectMeta.Name != "redis-config" {
				keepObjects = append(keepObjects, objects[rawExtIdx])
			}
		case *imagev1.ImageStream:
			if _, ok := highlyAvailableExternalDatabases[obj.ObjectMeta.Name]; !ok {
				keepObjects = append(keepObjects, objects[rawExtIdx])
			}
		default:
			keepObjects = append(keepObjects, objects[rawExtIdx])
		}
	}

	return keepObjects
}

func (ha *HighAvailability) DeleteDBRelatedParameters(template *templatev1.Template) {
	keepParams := []templatev1.Parameter{}
	dbParamsToDelete := map[string]bool{
		"MYSQL_IMAGE":         true,
		"REDIS_IMAGE":         true,
		"MYSQL_USER":          true,
		"MYSQL_PASSWORD":      true,
		"MYSQL_DATABASE":      true,
		"MYSQL_ROOT_PASSWORD": true,
	}

	for paramIdx := range template.Parameters {
		paramName := template.Parameters[paramIdx].Name
		if _, ok := dbParamsToDelete[paramName]; !ok {
			keepParams = append(keepParams, template.Parameters[paramIdx])
		}
	}

	template.Parameters = keepParams
}

func (ha *HighAvailability) UpdateDatabasesURLS(objects []runtime.RawExtension) {
	for rawExtIdx := range objects {
		obj := objects[rawExtIdx].Object
		secret, ok := obj.(*v1.Secret)
		if ok {
			switch secret.Name {
			case "system-redis":
				secret.StringData["URL"] = ha.Options.systemRedisURL
			case "backend-redis":
				secret.StringData["REDIS_STORAGE_URL"] = ha.Options.backendRedisStorageEndpoint
				secret.StringData["REDIS_QUEUES_URL"] = ha.Options.backendRedisQueuesEndpoint
			}
		}
	}
}
