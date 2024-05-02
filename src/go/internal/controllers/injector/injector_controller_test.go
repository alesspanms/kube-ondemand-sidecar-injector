package injector

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alesspanms/kube-ondemand-sidecar-injector/internal/kube"
	"github.com/alesspanms/kube-ondemand-sidecar-injector/internal/logging"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	injectormodels "github.com/alesspanms/kube-ondemand-sidecar-injector/internal/models/injector"
)

func TestNew(t *testing.T) {
	logger := logging.New()

	kubeClient := new(kube.KubeClientMock)

	controller := New(logger, kubeClient)

	assert.Equal(t, logger, controller.logger)
	assert.Equal(t, kubeClient, controller.kubeClient)
}

func TestGetDeployments(t *testing.T) {
	// You'll need to mock the kube.IKubeClient interface and its methods
	// to test this function. Here's a basic example:

	kubeClient := new(kube.KubeClientMock)
	expectedDeployments := []injectormodels.Deployment{{Name: "test-deployment"}}
	kubeClient.On("GetDeployments", "test-namespace").Return(expectedDeployments, nil)

	logger := logging.New()

	controller := New(logger, kubeClient)

	// Create a new gin context and request
	// gin.Default()
	w, context := createPostRequestFor("/api/injector/GetDeployments", strings.NewReader(`{"Namespace": "test-namespace"}`))

	controller.GetDeployments(context)

	// Check that the HTTP response status code is 200
	assert.Equal(t, http.StatusOK, w.Code)

	// Check that the HTTP response body contains the expected deployments
	var deployments []injectormodels.Deployment
	err := json.Unmarshal(w.Body.Bytes(), &deployments)
	assert.NoError(t, err)
	assert.Equal(t, expectedDeployments, deployments)
}

func TestGetDeploymentsErrorBinding(t *testing.T) {
	// You'll need to mock the kube.IKubeClient interface and its methods
	// to test this function. Here's a basic example:

	kubeClient := new(kube.KubeClientMock)
	kubeClient.On("GetDeployments", "test-namespace").Return(nil, assert.AnError)

	logger := logging.New()

	controller := New(logger, kubeClient)

	// Create a new gin context and request
	// gin.Default()
	w, context := createPostRequestFor("/api/injector/GetDeployments", strings.NewReader(`{"IncorrectNamespaceProperty": "test-namespace"}`))

	controller.GetDeployments(context)

	// Check that the HTTP response status code is 400
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Check that the HTTP response body contains the expected error message
	assert.JSONEq(t, `{"error": "Error binding JSON"}`, w.Body.String())
}

func TestGetDeploymentsErrorGettingDeployments(t *testing.T) {
	// You'll need to mock the kube.IKubeClient interface and its methods
	// to test this function. Here's a basic example:

	kubeClient := new(kube.KubeClientMock)
	kubeClient.On("GetDeployments", "test-namespace").Return(nil, assert.AnError)

	logger := logging.New()

	controller := New(logger, kubeClient)

	// Create a new gin context and request
	// gin.Default()
	w, context := createPostRequestFor("/api/injector/GetDeployments", strings.NewReader(`{"Namespace": "test-namespace"}`))

	controller.GetDeployments(context)

	// Check that the HTTP response status code is 400
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Check that the HTTP response body contains the expected error message
	assert.JSONEq(t, `{"error": "Error getting deployments: assert.AnError general error for testing"}`, w.Body.String())
}

func TestGetSingleDeployment(t *testing.T) {
	// You'll need to mock the kube.IKubeClient interface and its methods
	// to test this function. Here's a basic example:

	kubeClient := new(kube.KubeClientMock)
	expectedDeployment := injectormodels.Deployment{Name: "test-deployment"}
	kubeClient.On("GetSingleDeployment", "test-namespace", "test-deployment").Return(expectedDeployment, nil)

	logger := logging.New()

	controller := New(logger, kubeClient)

	// Create a new gin context and request
	// gin.Default()
	w, context := createPostRequestFor("/api/injector/GetSingleDeployment", strings.NewReader(`{"Namespace": "test-namespace", "DeploymentName": "test-deployment"}`))

	controller.GetSingleDeployment(context)

	// Check that the HTTP response status code is 200
	assert.Equal(t, http.StatusOK, w.Code)

	// Check that the HTTP response body contains the expected deployment
	var deployment injectormodels.Deployment
	err := json.Unmarshal(w.Body.Bytes(), &deployment)
	assert.NoError(t, err)
	assert.Equal(t, expectedDeployment, deployment)
}

func TestGetSingleDeploymentErrorBinding(t *testing.T) {
	// You'll need to mock the kube.IKubeClient interface and its methods
	// to test this function. Here's a basic example:

	kubeClient := new(kube.KubeClientMock)
	kubeClient.On("GetSingleDeployment", "test-namespace", "test-deployment").Return(injectormodels.Deployment{}, assert.AnError)

	logger := logging.New()

	controller := New(logger, kubeClient)

	// Create a new gin context and request
	// gin.Default()
	w, context := createPostRequestFor("/api/injector/GetSingleDeployment", strings.NewReader(`{"IncorrectNamespaceProperty": "test-namespace", "DeploymentName": "test-deployment"}`))

	controller.GetSingleDeployment(context)

	// Check that the HTTP response status code is 400
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Check that the HTTP response body contains the expected error message
	assert.JSONEq(t, `{"error": "Error binding JSON"}`, w.Body.String())
}

func TestGetSingleDeploymentErrorGettingDeployment(t *testing.T) {
	// You'll need to mock the kube.IKubeClient interface and its methods
	// to test this function. Here's a basic example:

	kubeClient := new(kube.KubeClientMock)
	kubeClient.On("GetSingleDeployment", "test-namespace", "test-deployment").Return(nil, assert.AnError)

	logger := logging.New()

	controller := New(logger, kubeClient)

	// Create a new gin context and request
	// gin.Default()
	w, context := createPostRequestFor("/api/injector/GetSingleDeployment", strings.NewReader(`{"Namespace": "test-namespace", "DeploymentName": "test-deployment"}`))

	controller.GetSingleDeployment(context)

	// Check that the HTTP response status code is 400
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Check that the HTTP response body contains the expected error message
	assert.JSONEq(t, `{"error": "Error getting single deployment: assert.AnError general error for testing"}`, w.Body.String())
}

func TestSetSidecar(t *testing.T) {
	// You'll need to mock the kube.IKubeClient interface and its methods
	// to test this function. Here's a basic example:

	kubeClient := new(kube.KubeClientMock)
	expectedDeployment := injectormodels.Deployment{Name: "test-deployment"}
	kubeClient.On("SetSidecar", mock.Anything).Return(expectedDeployment, nil)

	logger := logging.New()

	controller := New(logger, kubeClient)

	// Create a new gin context and request
	// gin.Default()
	w, context := createPostRequestFor("/api/injector/SetSidecar", strings.NewReader(`{"Namespace": "test-namespace", "DeploymentName": "test-deployment", "SidecarContainerName": "sidecar-container", "SidecarImage": "sidecar-image", "VolumeMounts": [], "Command": ["/bin/sh"]}`))

	controller.SetSidecar(context)

	// Check that the HTTP response status code is 200
	assert.Equal(t, http.StatusOK, w.Code)

	// Check that the HTTP response body contains the expected deployment
	var deployment injectormodels.Deployment
	err := json.Unmarshal(w.Body.Bytes(), &deployment)
	assert.NoError(t, err)
	assert.Equal(t, expectedDeployment, deployment)
}

func TestSetSidecarErrorBinding(t *testing.T) {
	// You'll need to mock the kube.IKubeClient interface and its methods
	// to test this function. Here's a basic example:

	kubeClient := new(kube.KubeClientMock)
	kubeClient.On("SetSidecar", mock.Anything).Return(injectormodels.Deployment{}, assert.AnError)

	logger := logging.New()

	controller := New(logger, kubeClient)

	// Create a new gin context and request
	// gin.Default()
	w, context := createPostRequestFor("/api/injector/SetSidecar", strings.NewReader(`{"IncorrectNamespaceProperty": "test-namespace", "DeploymentName": "test-deployment", "SidecarContainerName": "sidecar-container", "SidecarImage": "sidecar-image", "VolumeMounts": [], "Command": ["/bin/sh"]}`))

	controller.SetSidecar(context)

	// Check that the HTTP response status code is 400
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Check that the HTTP response body contains the expected error message
	assert.JSONEq(t, `{"error": "Error binding JSON"}`, w.Body.String())
}

func TestSetSidecarErrorSettingSidecar(t *testing.T) {
	// You'll need to mock the kube.IKubeClient interface and its methods
	// to test this function. Here's a basic example:

	kubeClient := new(kube.KubeClientMock)
	kubeClient.On("SetSidecar", mock.Anything).Return(nil, assert.AnError)

	logger := logging.New()

	controller := New(logger, kubeClient)

	// Create a new gin context and request
	// gin.Default()
	w, context := createPostRequestFor("/api/injector/SetSidecar", strings.NewReader(`{"Namespace": "test-namespace", "DeploymentName": "test-deployment", "SidecarContainerName": "sidecar-container", "SidecarImage": "sidecar-image", "VolumeMounts": [], "Command": ["/bin/sh"]}`))

	controller.SetSidecar(context)

	// Check that the HTTP response status code is 400
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Check that the HTTP response body contains the expected error message
	assert.JSONEq(t, `{"error": "Error on setting sidecar: assert.AnError general error for testing"}`, w.Body.String())
}

func TestClearSidecar(t *testing.T) {
	// You'll need to mock the kube.IKubeClient interface and its methods
	// to test this function. Here's a basic example:

	kubeClient := new(kube.KubeClientMock)
	expectedDeployment := injectormodels.Deployment{Name: "test-deployment"}
	kubeClient.On("ClearSidecar", mock.Anything).Return(expectedDeployment, nil)

	logger := logging.New()

	controller := New(logger, kubeClient)

	// Create a new gin context and request
	// gin.Default()
	w, context := createPostRequestFor("/api/injector/ClearSidecar", strings.NewReader(`{"Namespace": "test-namespace", "DeploymentName": "test-deployment", "SidecarContainerName": "sidecar-container"}`))

	controller.ClearSidecar(context)

	// Check that the HTTP response status code is 200
	assert.Equal(t, http.StatusOK, w.Code)

	// Check that the HTTP response body contains the expected deployment
	var deployment injectormodels.Deployment
	err := json.Unmarshal(w.Body.Bytes(), &deployment)
	assert.NoError(t, err)
	assert.Equal(t, expectedDeployment, deployment)
}

func TestClearSidecarErrorBinding(t *testing.T) {
	// You'll need to mock the kube.IKubeClient interface and its methods
	// to test this function. Here's a basic example:

	kubeClient := new(kube.KubeClientMock)
	kubeClient.On("ClearSidecar", mock.Anything).Return(injectormodels.Deployment{}, assert.AnError)

	logger := logging.New()

	controller := New(logger, kubeClient)

	// Create a new gin context and request
	// gin.Default()
	w, context := createPostRequestFor("/api/injector/ClearSidecar", strings.NewReader(`{"IncorrectNamespaceProperty": "test-namespace", "DeploymentName": "test-deployment", "SidecarContainerName": "sidecar-container"}`))

	controller.ClearSidecar(context)

	// Check that the HTTP response status code is 400
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Check that the HTTP response body contains the expected error message
	assert.JSONEq(t, `{"error": "Error binding JSON"}`, w.Body.String())
}

func TestClearSidecarErrorClearingSidecar(t *testing.T) {
	// You'll need to mock the kube.IKubeClient interface and its methods
	// to test this function. Here's a basic example:

	kubeClient := new(kube.KubeClientMock)
	kubeClient.On("ClearSidecar", mock.Anything).Return(nil, assert.AnError)

	logger := logging.New()

	controller := New(logger, kubeClient)

	// Create a new gin context and request
	// gin.Default()
	w, context := createPostRequestFor("/api/injector/ClearSidecar", strings.NewReader(`{"Namespace": "test-namespace", "DeploymentName": "test-deployment", "SidecarContainerName": "sidecar-container"}`))

	controller.ClearSidecar(context)

	// Check that the HTTP response status code is 400
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Check that the HTTP response body contains the expected error message
	assert.JSONEq(t, `{"error": "Error clearing sidecar: assert.AnError general error for testing"}`, w.Body.String())
}

func createPostRequestFor(url string, bodyReader io.Reader) (*httptest.ResponseRecorder, *gin.Context) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", url, bodyReader)
	req.Header.Set("Content-Type", "application/json")
	context, _ := gin.CreateTestContext(w)
	context.Request = req
	return w, context
}
