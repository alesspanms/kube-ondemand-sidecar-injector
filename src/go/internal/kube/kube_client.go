package kube

import (
	"context"
	"errors"
	"flag"
	"path/filepath"

	"go.uber.org/zap"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	"github.com/alesspanms/kube-ondemand-sidecar-injector/internal/logging"
	injectormodels "github.com/alesspanms/kube-ondemand-sidecar-injector/internal/models/injector"
)

type IKubeClient interface {
	GetDeployments(namespace string) ([]injectormodels.Deployment, error)
	GetSingleDeployment(namespace string, name string) (injectormodels.Deployment, error)
	SetSidecar(payload *injectormodels.SetSidecarPayload) (injectormodels.Deployment, error)
	ClearSidecar(payload *injectormodels.ClearSidecarPayload) (injectormodels.Deployment, error)
}

type KubeClient struct {
	logger            *logging.Logger
	sidecarNamePrefix string
	clientset         *kubernetes.Clientset
	config            *rest.Config
}

// NewKubeClient creates a new instance of the KubeClient
func New(logger *logging.Logger, sidecarNamePrefix string) IKubeClient {
	kubeClient := &KubeClient{
		logger:            logger,
		sidecarNamePrefix: sidecarNamePrefix + "-",
	}
	kubeClient.init()
	return kubeClient
}

func (kc *KubeClient) init() {
	var err error // Declare err variable

	// Create a Kubernetes client
	kc.config, err = rest.InClusterConfig()
	if err != nil {
		// panic(err.Error())
		// fallback to local kubeconfig, if any

		var kubeconfig *string
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		flag.Parse()

		// use the current context in kubeconfig
		kc.config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			panic(err.Error())
		}
	}

	clientset, err := kubernetes.NewForConfig(kc.config)
	if err != nil {
		panic(err.Error())
	}

	kc.clientset = clientset
}

func (kc *KubeClient) GetDeployments(namespace string) (result []injectormodels.Deployment, err error) {

	if namespace == "" {
		err = errors.New("namespace is required")
		return
	}

	kc.logger.Log().Info("Getting deployments", zap.String("namespace", namespace))

	deployments, err := kc.clientset.AppsV1().Deployments(namespace).List(context.Background(), v1.ListOptions{})

	if err != nil {
		//panic(err.Error())
		return
	}

	result = make([]injectormodels.Deployment, 0)
	for _, deployment := range deployments.Items {

		kc.logger.Log().Info("GetDeployments - Found Deployment", zap.String("name", deployment.Spec.Template.Name), zap.String("namespace", deployment.Namespace))

		internalDeployment := convertToInternalModel(deployment)

		result = append(result, internalDeployment)
	}

	return
}

func (kc *KubeClient) GetSingleDeployment(namespace string, name string) (result injectormodels.Deployment, err error) {

	if namespace == "" {
		err = errors.New("namespace is required")
		return
	}

	kc.logger.Log().Info("Getting Single deployment", zap.String("namespace", namespace), zap.String("name", name))

	deployment, err := kc.clientset.AppsV1().Deployments(namespace).Get(context.Background(), name, v1.GetOptions{})

	if err != nil {
		//panic(err.Error())
		return
	}

	kc.logger.Log().Info("GetSingleDeployment - Found Deployment", zap.String("name", deployment.Spec.Template.Name), zap.String("namespace", deployment.Namespace))

	result = convertToInternalModel(*deployment)

	return
}

func (kc *KubeClient) SetSidecar(payload *injectormodels.SetSidecarPayload) (result injectormodels.Deployment, err error) {

	kc.logger.Log().Info("SetSidecar ", zap.String("DeploymentName", payload.DeploymentName), zap.String("Namespace", payload.Namespace), zap.String("SidecarImage", payload.SidecarImage))

	if payload.Namespace == "" {
		err = errors.New("namespace is required")
		return
	}

	if payload.DeploymentName == "" {
		err = errors.New("DeploymentName is required")
		return
	}

	if payload.SidecarImage == "" {
		err = errors.New("SidecarImage is required")
		return
	}

	// will be used the container's image ENTRYPOINT if not provided
	// if payload.Command == nil {
	// 	payload.Command = []string{"/bin/sh", "-c", "while true; do sleep 10; done"}
	// }

	deployment, err := kc.clientset.AppsV1().Deployments(payload.Namespace).Get(context.Background(), payload.DeploymentName, v1.GetOptions{})

	if err != nil {
		//panic(err.Error())
		return
	}

	kc.logger.Log().Info("SetSidecar - Found Deployment", zap.String("name", deployment.Spec.Template.Name), zap.String("namespace", deployment.Namespace))

	v1Container := corev1.Container{
		Name:         kc.sidecarNamePrefix + payload.SidecarContainerName,
		Image:        payload.SidecarImage,
		VolumeMounts: make([]corev1.VolumeMount, len(payload.VolumeMounts)),
	}

	if payload.Command != nil && len(payload.Command) > 0 {
		v1Container.Command = payload.Command
	}

	for i, volumeMount := range payload.VolumeMounts {

		found := false
		for _, volume := range deployment.Spec.Template.Spec.Volumes {
			if volume.Name == volumeMount.Name {
				found = true
				break
			}
		}

		if !found {
			err = errors.New("volume '" + volumeMount.Name + "' not found. cannot continue.")
			return
		}

		v1Container.VolumeMounts[i] = corev1.VolumeMount{
			Name:      volumeMount.Name,
			MountPath: volumeMount.MountPath,
		}
	}

	deployment.Spec.Template.Spec.Containers = append(deployment.Spec.Template.Spec.Containers, v1Container)

	_, err = kc.clientset.AppsV1().Deployments(payload.Namespace).Update(context.Background(), deployment, v1.UpdateOptions{})
	if err != nil {
		//panic(err.Error())
		return
	}

	updatedDeployment, err := kc.clientset.AppsV1().Deployments(payload.Namespace).Get(context.Background(), payload.DeploymentName, v1.GetOptions{})
	if err != nil {
		//panic(err.Error())
		return
	}

	result = convertToInternalModel(*updatedDeployment)

	kc.logger.Log().Info("SetSidecar - Updated Deployment", zap.String("name", updatedDeployment.Spec.Template.Name), zap.String("namespace", updatedDeployment.Namespace))

	return
}

func (kc *KubeClient) ClearSidecar(payload *injectormodels.ClearSidecarPayload) (result injectormodels.Deployment, err error) {

	kc.logger.Log().Info("ClearSidecar ", zap.String("DeploymentName", payload.DeploymentName), zap.String("Namespace", payload.Namespace), zap.String("SidecarContainerName", payload.SidecarContainerName))

	if payload.Namespace == "" {
		err = errors.New("namespace is required")
		return
	}

	if payload.DeploymentName == "" {
		err = errors.New("DeploymentName is required")
		return
	}

	if payload.SidecarContainerName == "" {
		err = errors.New("SidecarContainerName is required")
		return
	}

	deployment, err := kc.clientset.AppsV1().Deployments(payload.Namespace).Get(context.Background(), payload.DeploymentName, v1.GetOptions{})

	if err != nil {
		//panic(err.Error())
		return
	}

	kc.logger.Log().Info("SetSidecar - Found Deployment", zap.String("name", deployment.Spec.Template.Name), zap.String("namespace", deployment.Namespace))

	//find index of named container
	index := -1
	for i, container := range deployment.Spec.Template.Spec.Containers {
		if container.Name == kc.sidecarNamePrefix+payload.SidecarContainerName {
			index = i
			break
		}
	}

	if index == -1 {
		err = errors.New("container ' " + payload.SidecarContainerName + "' not found ")
		return
	} else {
		deployment.Spec.Template.Spec.Containers = append(deployment.Spec.Template.Spec.Containers[:index], deployment.Spec.Template.Spec.Containers[index+1:]...)
	}

	_, err = kc.clientset.AppsV1().Deployments(payload.Namespace).Update(context.Background(), deployment, v1.UpdateOptions{})
	if err != nil {
		//panic(err.Error())
		return
	}

	updatedDeployment, err := kc.clientset.AppsV1().Deployments(payload.Namespace).Get(context.Background(), payload.DeploymentName, v1.GetOptions{})
	if err != nil {
		//panic(err.Error())
		return
	}

	result = convertToInternalModel(*updatedDeployment)

	kc.logger.Log().Info("SetSidecar - Updated Deployment", zap.String("name", updatedDeployment.Spec.Template.Name), zap.String("namespace", updatedDeployment.Namespace))

	return
}

// private functions and methods

func convertToInternalModel(deployment appsv1.Deployment) injectormodels.Deployment {

	volumeNames := make([]string, len(deployment.Spec.Template.Spec.Volumes))
	for i, volume := range deployment.Spec.Template.Spec.Volumes {
		volumeNames[i] = volume.Name
	}

	return injectormodels.Deployment{
		Name:        deployment.Name,
		Namespace:   deployment.Namespace,
		VolumeNames: volumeNames,
	}
}
