package component

import (
	appsv1 "github.com/openshift/api/apps/v1"
	templatev1 "github.com/openshift/api/template/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	S3SecretAWSAccessKeyIdFieldName     = "AWS_ACCESS_KEY_ID"
	S3SecretAWSSecretAccessKeyFieldName = "AWS_SECRET_ACCESS_KEY"
)

type S3 struct {
	Options *S3Options
}

func NewS3(options *S3Options) *S3 {
	return &S3{Options: options}
}

func (s3 *S3) Objects() []runtime.RawExtension {
	s3AWSSecret := s3.buildS3AWSSecret()

	objects := []runtime.RawExtension{
		runtime.RawExtension{Object: s3AWSSecret},
	}
	return objects
}

func (s3 *S3) RemoveSystemStorageReferences(objects []runtime.RawExtension) {
	for _, rawExtension := range objects {
		obj := rawExtension.Object
		dc, ok := obj.(*appsv1.DeploymentConfig)
		if ok {
			if dc.ObjectMeta.Name == "system-app" || dc.ObjectMeta.Name == "system-sidekiq" {

				// Remove system-storage references in the VolumeMount fields of the containers
				for containerIdx := range dc.Spec.Template.Spec.Containers {
					container := &dc.Spec.Template.Spec.Containers[containerIdx]
					resIdx := -1
					for vmIdx, vm := range container.VolumeMounts {
						if vm.Name == "system-storage" {
							resIdx = vmIdx
							break
						}
					}
					if resIdx != -1 {
						container.VolumeMounts = append(container.VolumeMounts[:resIdx], container.VolumeMounts[resIdx+1:]...)
					}
				}

				// Remove system-storage references in the Volumes fields of the containers
				resIdx := -1
				for volIdx := range dc.Spec.Template.Spec.Volumes {
					vol := &dc.Spec.Template.Spec.Volumes[volIdx]
					if vol.Name == "system-storage" {
						resIdx = volIdx
						break
					}
				}
				if resIdx != -1 {
					dc.Spec.Template.Spec.Volumes = append(dc.Spec.Template.Spec.Volumes[:resIdx], dc.Spec.Template.Spec.Volumes[resIdx+1:]...)
				}

				// Remove system-storage references in the Volumes fields of the pre-hook in system-app
				if dc.ObjectMeta.Name == "system-app" {
					resIdx = -1
					for volIdx := range dc.Spec.Strategy.RollingParams.Pre.ExecNewPod.Volumes {
						vol := &dc.Spec.Strategy.RollingParams.Pre.ExecNewPod.Volumes[volIdx]
						if *vol == "system-storage" {
							resIdx = volIdx
							break
						}
					}
					if resIdx != -1 {
						dc.Spec.Strategy.RollingParams.Pre.ExecNewPod.Volumes = append(dc.Spec.Strategy.RollingParams.Pre.ExecNewPod.Volumes[:resIdx], dc.Spec.Strategy.RollingParams.Pre.ExecNewPod.Volumes[resIdx+1:]...)
					}
				}
			}
		}
	}
}

// Remove the RWX_STORAGE_CLASS parameter because it is used only for the system-storage PersistentVolumeClaim
func (s3 *S3) RemoveRWXStorageClassParameter(template *templatev1.Template) {
	for paramIdx, param := range template.Parameters {
		if param.Name == "RWX_STORAGE_CLASS" {
			template.Parameters = append(template.Parameters[:paramIdx], template.Parameters[paramIdx+1:]...)
			break
		}
	}
}

func (s3 *S3) getNewCfgMapElements() []v1.EnvVar {
	return []v1.EnvVar{
		envVarFromConfigMap("FILE_UPLOAD_STORAGE", "system-environment", "FILE_UPLOAD_STORAGE"),
		envVarFromSecret("AWS_ACCESS_KEY_ID", s3.Options.awsCredentialsSecret, S3SecretAWSAccessKeyIdFieldName),
		envVarFromSecret("AWS_SECRET_ACCESS_KEY", s3.Options.awsCredentialsSecret, S3SecretAWSSecretAccessKeyFieldName),
		envVarFromConfigMap("AWS_BUCKET", "system-environment", "AWS_BUCKET"),
		envVarFromConfigMap("AWS_REGION", "system-environment", "AWS_REGION"),
	}
}

func (s3 *S3) AddCfgMapElemsToSystemBaseEnv(objects []runtime.RawExtension) {
	newCfgMapElements := s3.getNewCfgMapElements()
	for _, rawExtension := range objects {
		obj := rawExtension.Object
		dc, ok := obj.(*appsv1.DeploymentConfig)
		if ok {
			if dc.ObjectMeta.Name == "system-app" || dc.ObjectMeta.Name == "system-sidekiq" {
				if dc.ObjectMeta.Name == "system-app" {
					dc.Spec.Strategy.RollingParams.Pre.ExecNewPod.Env = append(dc.Spec.Strategy.RollingParams.Pre.ExecNewPod.Env, newCfgMapElements...)
				}

				for containerIdx := range dc.Spec.Template.Spec.Containers {
					container := &dc.Spec.Template.Spec.Containers[containerIdx]
					container.Env = append(container.Env, newCfgMapElements...)
				}
			}
		}
	}
}

func (s3 *S3) AddS3PostprocessOptionsToSystemEnvironmentCfgMap(objects []runtime.RawExtension) {
	var systemEnvCfgMap *v1.ConfigMap

	for rawExtIdx := range objects {
		obj := objects[rawExtIdx].Object
		cfgmap, ok := obj.(*v1.ConfigMap)
		if ok {
			if cfgmap.Name == "system-environment" {
				systemEnvCfgMap = cfgmap
				break
			}
		}
	}

	systemEnvCfgMap.Data["FILE_UPLOAD_STORAGE"] = s3.Options.fileUploadStorage
	systemEnvCfgMap.Data["AWS_BUCKET"] = s3.Options.awsBucket
	systemEnvCfgMap.Data["AWS_REGION"] = s3.Options.awsRegion
}

func (s3 *S3) RemoveSystemStoragePVC(objects []runtime.RawExtension) []runtime.RawExtension {
	res := objects

	for idx, rawExtension := range res {
		obj := rawExtension.Object
		pvc, ok := obj.(*v1.PersistentVolumeClaim)
		if ok {
			if pvc.ObjectMeta.Name == "system-storage" {
				res = append(res[:idx], res[idx+1:]...) // This deletes the element in the array
				break
			}
		}
	}

	return res
}

func (s3 *S3) buildS3AWSSecret() *v1.Secret {
	return &v1.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Secret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: s3.Options.awsCredentialsSecret,
		},
		StringData: map[string]string{
			S3SecretAWSAccessKeyIdFieldName:     s3.Options.awsAccessKeyId,
			S3SecretAWSSecretAccessKeyFieldName: s3.Options.awsSecretAccessKey,
		},
		Type: v1.SecretTypeOpaque,
	}
}
