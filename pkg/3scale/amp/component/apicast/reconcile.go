package apicast

import (
	"context"
	"fmt"
	"reflect"

	appsv1operator "github.com/3scale/3scale-operator/pkg/apis/apps"
	appsv1alpha1 "github.com/3scale/3scale-operator/pkg/apis/apps/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	extensions "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/3scale/3scale-operator/pkg/k8sutils"
)

type Reconciler struct {
	Client client.Client
	Logger logr.Logger
}

func (r *Reconciler) Reconcile(request reconcile.Request) error {
	apicastCR, err := r.getAPIcast(request)
	if err != nil {
		return err
	}
	// if apicast is nil and we did not have an error it means it does
	// not exist and we don't want to return an error because we
	// don't want to requeue the request
	if apicastCR == nil {
		return nil
	}

	err = r.setMissingOptionalDefaultValues(apicastCR)
	if err != nil {
		return err
	}

	desiredAPIcast, err := r.internalAPIcast(apicastCR)
	if err != nil {
		return err
	}

	desiredAPIcastAdminPortalEndpointSecret, err := desiredAPIcast.AdminPortalEndpointSecret()
	if err != nil {
		return err
	}

	err = r.reconcileAdminPortalEndpointSecret(desiredAPIcastAdminPortalEndpointSecret)
	if err != nil {
		return err
	}

	err = r.reconcileDeployment(*desiredAPIcast.Deployment())
	if err != nil {
		return err
	}

	err = r.reconcileService(*desiredAPIcast.Service())
	if err != nil {
		return err
	}

	if apicastCR.Spec.ExposedHostname != nil {
		err = r.reconcileIngress(*desiredAPIcast.Ingress())
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Reconciler) adminPortalCredentials(existingAPIcast *appsv1alpha1.APIcast) (ApicastAdminPortalCredentials, error) {
	adminPortalURLSecretKeyRef := existingAPIcast.Spec.AdminPortal.URLSecretKeyRef
	if adminPortalURLSecretKeyRef.Name == "" {
		return ApicastAdminPortalCredentials{}, fmt.Errorf("Field 'name' not specified for URL Secret Key Reference")
	}
	adminPortalURLSecretNamespacedName := types.NamespacedName{
		Name:      adminPortalURLSecretKeyRef.Name,
		Namespace: existingAPIcast.Namespace,
	}

	adminPortalURLSecret := v1.Secret{}
	err := r.Client.Get(context.TODO(), adminPortalURLSecretNamespacedName, &adminPortalURLSecret)
	if err != nil {
		return ApicastAdminPortalCredentials{}, err
	}
	secretStringData := k8sutils.SecretStringDataFromData(adminPortalURLSecret)
	adminPortalURL, ok := secretStringData[adminPortalURLSecretKeyRef.Key]
	if !ok {
		return ApicastAdminPortalCredentials{}, fmt.Errorf("Key '%s' not found in secret '%s'", adminPortalURLSecretKeyRef.Key, adminPortalURLSecretKeyRef.Name)
	}

	adminPortalAccessTokenKeyRef := existingAPIcast.Spec.AdminPortal.AccessTokenSecretKeyRef
	if adminPortalAccessTokenKeyRef.Name == "" {
		return ApicastAdminPortalCredentials{}, fmt.Errorf("Field 'name' not specified for URL Access Token Key Reference")
	}
	adminPortalAccessTokenNamespacedName := types.NamespacedName{
		Name:      adminPortalAccessTokenKeyRef.Name,
		Namespace: existingAPIcast.Namespace,
	}
	adminPortalAccessTokenSecret := v1.Secret{}
	err = r.Client.Get(context.TODO(), adminPortalAccessTokenNamespacedName, &adminPortalAccessTokenSecret)
	if err != nil {
		return ApicastAdminPortalCredentials{}, err
	}

	secretStringData = k8sutils.SecretStringDataFromData(adminPortalAccessTokenSecret)
	adminPortalAccessToken, ok := secretStringData[adminPortalAccessTokenKeyRef.Key]
	if !ok {
		return ApicastAdminPortalCredentials{}, fmt.Errorf("Key '%s' not found in secret '%s'", adminPortalAccessTokenKeyRef.Key, adminPortalAccessTokenKeyRef.Name)
	}
	result := ApicastAdminPortalCredentials{
		URL:         adminPortalURL,
		AccessToken: adminPortalAccessToken,
	}
	return result, nil
}

func (r *Reconciler) checkAdditionalEnvironmentConfigurationSecretRef(existingAPIcast *appsv1alpha1.APIcast) error {
	environmentConfigurationSecretRef := existingAPIcast.Spec.EnvironmentConfigurationSecretRef
	if environmentConfigurationSecretRef == nil {
		return nil
	}

	if environmentConfigurationSecretRef.Name == "" {
		return fmt.Errorf("Field 'name' not specified for Additional Env Secret Reference")
	}

	environmentConfigurationSecretNamespacedName := types.NamespacedName{
		Name:      environmentConfigurationSecretRef.Name,
		Namespace: existingAPIcast.Namespace,
	}

	environmentConfigurationSecret := v1.Secret{}
	err := r.Client.Get(context.TODO(), environmentConfigurationSecretNamespacedName, &environmentConfigurationSecret)
	return err
}

func (r *Reconciler) internalAPIcast(existingAPIcast *appsv1alpha1.APIcast) (Apicast, error) {
	apicastFullName := "apicast-" + existingAPIcast.Name
	apicastExposedHostname := ""
	if existingAPIcast.Spec.ExposedHostname != nil {
		apicastExposedHostname = *existingAPIcast.Spec.ExposedHostname
	}
	apicastOwnerRef := asOwner(existingAPIcast)

	adminPortalCredentials, err := r.adminPortalCredentials(existingAPIcast)
	if err != nil {
		return Apicast{}, err
	}

	err = r.checkAdditionalEnvironmentConfigurationSecretRef(existingAPIcast)
	if err != nil {
		return Apicast{}, err
	}

	internalApicastResult := Apicast{
		deploymentName:         apicastFullName,
		serviceName:            apicastFullName,
		initialReplicas:        int32(*existingAPIcast.Spec.Replicas),
		appLabel:               "apicast",
		serviceAccountName:     *existingAPIcast.Spec.ServiceAccount,
		image:                  *existingAPIcast.Spec.Image,
		exposedHostname:        apicastExposedHostname,
		namespace:              existingAPIcast.Namespace,
		ownerReference:         &apicastOwnerRef,
		additionalEnvironment:  existingAPIcast.Spec.EnvironmentConfigurationSecretRef,
		adminPortalCredentials: adminPortalCredentials,
	}

	return internalApicastResult, err
}

func (r *Reconciler) namespacedName(object metav1.Object) types.NamespacedName {
	return types.NamespacedName{
		Name:      object.GetName(),
		Namespace: object.GetNamespace(),
	}
}

func (r *Reconciler) getAPIcast(request reconcile.Request) (*appsv1alpha1.APIcast, error) {
	instance := appsv1alpha1.APIcast{}
	err := r.Client.Get(context.TODO(), request.NamespacedName, &instance)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return &instance, nil
}

func (r *Reconciler) setMissingOptionalDefaultValues(apicastCR *appsv1alpha1.APIcast) error {
	var defaultAPIcastReplicas int64 = 1
	defaultServiceAccount := "default"
	defaultAPIcastImage := "registry.access.redhat.com/3scale-amp25/apicast-gateway"
	missingOptionalField := false

	if apicastCR.Spec.Replicas == nil {
		apicastCR.Spec.Replicas = &defaultAPIcastReplicas
		missingOptionalField = true
	}
	if apicastCR.Spec.ServiceAccount == nil {
		apicastCR.Spec.ServiceAccount = &defaultServiceAccount
		missingOptionalField = true
	}
	if apicastCR.Spec.Image == nil {
		apicastCR.Spec.Image = &defaultAPIcastImage
		missingOptionalField = true
	}

	if missingOptionalField {
		err := r.Client.Update(context.TODO(), apicastCR)
		if err != nil {
			return err
		}
		return fmt.Errorf("APIcast resource missed optional fields. Requeuing request after having set them")
	}

	return nil
}

// asOwner returns an owner reference set as the tenant CR
func asOwner(a *appsv1alpha1.APIcast) metav1.OwnerReference {
	trueVar := true
	return metav1.OwnerReference{
		APIVersion: appsv1alpha1.SchemeGroupVersion.String(),
		Kind:       appsv1operator.APICastKind,
		Name:       a.Name,
		UID:        a.UID,
		Controller: &trueVar,
	}
}

func (r *Reconciler) reconcileDeployment(desiredDeployment appsv1.Deployment) error {
	existingDeployment := appsv1.Deployment{}
	err := r.Client.Get(context.TODO(), r.namespacedName(&desiredDeployment), &existingDeployment)
	if err != nil {
		if errors.IsNotFound(err) {
			err = r.Client.Create(context.TODO(), &desiredDeployment)
			r.Logger.Info("Creating DeploymentConfig...")
		}
		return err
	} else {
		// Comparing Spec directly always detects changes because they are automatically
		// added by the controller
		// if !reflect.DeepEqual(existingDeployment.Spec, desiredDeployment.Spec) {
		// 	existingDeployment.Spec = desiredDeployment.Spec
		// 	err = r.Client.Update(context.TODO(), &existingDeployment)
		// 	r.Logger.Info("Updating Deployment...")
		// 	if err != nil {
		// 		return err
		// 	}
		// }
		changed := false
		// TODO reconcile the desired fields
		if !reflect.DeepEqual(existingDeployment.Spec.Template.Labels, desiredDeployment.Spec.Template.Labels) {
			existingDeployment.Spec.Template.Labels = desiredDeployment.Spec.Template.Labels
			changed = true
		}
		if !reflect.DeepEqual(existingDeployment.Spec.Template.Spec.Containers[0].Image, desiredDeployment.Spec.Template.Spec.Containers[0].Image) {
			existingDeployment.Spec.Template.Spec.Containers[0].Image = desiredDeployment.Spec.Template.Spec.Containers[0].Image
			changed = true
		}
		if !reflect.DeepEqual(existingDeployment.Spec.Template.Spec.ServiceAccountName, desiredDeployment.Spec.Template.Spec.ServiceAccountName) {
			changed = true
			existingDeployment.Spec.Template.Spec.ServiceAccountName = desiredDeployment.Spec.Template.Spec.ServiceAccountName
		}
		if !reflect.DeepEqual(existingDeployment.Spec.Template.Spec.Containers[0].EnvFrom, desiredDeployment.Spec.Template.Spec.Containers[0].EnvFrom) {
			changed = true
			existingDeployment.Spec.Template.Spec.Containers[0].EnvFrom = desiredDeployment.Spec.Template.Spec.Containers[0].EnvFrom
		}
		if changed {
			err = r.Client.Update(context.TODO(), &existingDeployment)
			r.Logger.Info("Updating Deployment...")
			return err
		}
		return nil
	}
}

func (r *Reconciler) reconcileAdminPortalEndpointSecret(desiredAdminPortalSecret v1.Secret) error {
	existingAdminPortalSecret := v1.Secret{}
	err := r.Client.Get(context.TODO(), r.namespacedName(&desiredAdminPortalSecret), &existingAdminPortalSecret)
	if err != nil {
		if errors.IsNotFound(err) {
			err = r.Client.Create(context.TODO(), &desiredAdminPortalSecret)
			r.Logger.Info("Creating Admin Portal Endpoint Secret...")
		}
		return err
	} else {
		existingAdminPortalSecretStringData := k8sutils.SecretStringDataFromData(existingAdminPortalSecret)
		if !reflect.DeepEqual(existingAdminPortalSecretStringData, desiredAdminPortalSecret.StringData) {
			existingAdminPortalSecret.StringData = desiredAdminPortalSecret.StringData
			err = r.Client.Update(context.TODO(), &existingAdminPortalSecret)
			r.Logger.Info("Updating Admin Portal Endpoint Secret...")
			return err
		}
	}
	return nil
}

func (r *Reconciler) reconcileService(desiredService v1.Service) error {
	existingService := v1.Service{}
	err := r.Client.Get(context.TODO(), r.namespacedName(&desiredService), &existingService)
	if err != nil {
		if errors.IsNotFound(err) {
			err = r.Client.Create(context.TODO(), &desiredService)
			r.Logger.Info("Creating Service...")
		}
		return err
	} else {
		if !reflect.DeepEqual(existingService.Spec.Ports, desiredService.Spec.Ports) {
			existingService.Spec.Ports = desiredService.Spec.Ports

			err = r.Client.Update(context.TODO(), &existingService)
			r.Logger.Info("Updating Service...")
			if err != nil {
				return err
			}
		}
		return err
	}

}

func (r *Reconciler) reconcileIngress(desiredIngress extensions.Ingress) error {
	existingIngress := extensions.Ingress{}
	err := r.Client.Get(context.TODO(), r.namespacedName(&desiredIngress), &existingIngress)
	if err != nil {
		if errors.IsNotFound(err) {
			err = r.Client.Create(context.TODO(), &desiredIngress)
			r.Logger.Info("Creating Ingress...")
		}
		return err
	} else {
		if !reflect.DeepEqual(existingIngress.Spec, desiredIngress.Spec) {
			existingIngress.Spec = desiredIngress.Spec
			err = r.Client.Update(context.TODO(), &existingIngress)
			r.Logger.Info("Updating Ingress...")
			if err != nil {
				return err
			}
		}
		return err
	}
}
