package injector

import (
	"net/http"

	"github.com/alesspanms/kube-ondemand-sidecar-injector/internal/kube"
	"github.com/alesspanms/kube-ondemand-sidecar-injector/internal/logging"
	injectormodels "github.com/alesspanms/kube-ondemand-sidecar-injector/internal/models/injector"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type InjectorController struct {
	logger     *logging.Logger
	kubeClient kube.IKubeClient
}

// NewInjectorController creates a new instance of the InjectorController
func New(logger *logging.Logger, kubeClient kube.IKubeClient) *InjectorController {
	return &InjectorController{
		logger:     logger,
		kubeClient: kubeClient,
	}
}

// GetDeployments godoc
// @Summary      Obtain a list of Deployment objects
// @Description  Get deployments for a given namespace or all namespaces
// @Tags         injector
// @Accept       json
// @Produce      json
// @Param        payload   body      injectormodels.GetDeploymentsPayload  true  "GetDeploymentsPayload type"
// @Success      200  {object}  []injectormodels.Deployment
// Failure      400  {object}  httputil.HTTPError
// Failure      404  {object}  httputil.HTTPError
// Failure      500  {object}  httputil.HTTPError
// @Router       /api/injector/GetDeployments [post]
// @Security ApiKeyAuth
func (ic *InjectorController) GetDeployments(c *gin.Context) {

	var payload injectormodels.GetDeploymentsPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		ic.logger.Log().Error("Error binding JSON", zap.Error(err))

		c.JSON(http.StatusBadRequest, gin.H{"error": "Error binding JSON"})
		return
	}

	ic.logger.Log().Info("GetDeployments - Received request", zap.Any("payload", payload))

	deployments, err := ic.kubeClient.GetDeployments(payload.Namespace)

	if err != nil {
		ic.logger.Log().Error("Error getting deployments", zap.Error(err))

		c.JSON(http.StatusBadRequest, gin.H{"error": "Error getting deployments: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, deployments)
}

// GetSingleDeployment godoc
// @Summary      Obtain a specific Deployment objects
// @Description  Get deployments for a given namespace and name
// @Tags         injector
// @Accept       json
// @Produce      json
// @Param        payload   body      injectormodels.GetSingleDeploymentPayload  true  "GetSingleDeploymentPayload type"
// @Success      200  {object}  injectormodels.Deployment
// Failure      400  {object}  httputil.HTTPError
// Failure      404  {object}  httputil.HTTPError
// Failure      500  {object}  httputil.HTTPError
// @Router       /api/injector/GetSingleDeployment [post]
// @Security ApiKeyAuth
func (ic *InjectorController) GetSingleDeployment(c *gin.Context) {

	var payload injectormodels.GetSingleDeploymentPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		ic.logger.Log().Error("Error binding JSON", zap.Error(err))

		c.JSON(http.StatusBadRequest, gin.H{"error": "Error binding JSON"})
		return
	}

	ic.logger.Log().Info("GetSingleDeployment - Received request", zap.Any("payload", payload))

	deployment, err := ic.kubeClient.GetSingleDeployment(payload.Namespace, payload.DeploymentName)

	if err != nil {
		ic.logger.Log().Error("Error getting single deployment", zap.Error(err))

		c.JSON(http.StatusBadRequest, gin.H{"error": "Error getting single deployment: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, deployment)
}

// SetSidecar godoc
// @Summary      Activate the sidecar
// @Description  Set the sidecar for a given deployment
// @Tags         injector
// @Accept       json
// @Produce      json
// @Param        payload   body      injectormodels.SetSidecarPayload  true  "SetSidecarPayload type"
// @Success      200  {object}  injectormodels.Deployment
// Failure      400  {object}  httputil.HTTPError
// Failure      404  {object}  httputil.HTTPError
// Failure      500  {object}  httputil.HTTPError
// @Router       /api/injector/SetSidecar [post]
// @Security ApiKeyAuth
func (ic *InjectorController) SetSidecar(c *gin.Context) {

	var payload injectormodels.SetSidecarPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		ic.logger.Log().Error("Error binding JSON", zap.Error(err))

		c.JSON(http.StatusBadRequest, gin.H{"error": "Error binding JSON"})
		return
	}

	ic.logger.Log().Info("SetSidecar - Received request", zap.Any("payload", payload))

	deployment, err := ic.kubeClient.SetSidecar(&payload)

	if err != nil {
		ic.logger.Log().Error("Error on setting sidecar", zap.Error(err))

		c.JSON(http.StatusBadRequest, gin.H{"error": "Error on setting sidecar: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, deployment)
}

// ClearSidecar godoc
// @Summary      Remove the sidecar
// @Description  Remove the sidecar from a given deployment
// @Tags         injector
// @Accept       json
// @Produce      json
// @Param        payload   body      injectormodels.ClearSidecarPayload  true  "ClearSidecarPayload type"
// @Success      200  {object}  injectormodels.Deployment
// Failure      400  {object}  httputil.HTTPError
// Failure      404  {object}  httputil.HTTPError
// Failure      500  {object}  httputil.HTTPError
// @Router       /api/injector/ClearSidecar [post]
// @Security ApiKeyAuth
func (ic *InjectorController) ClearSidecar(c *gin.Context) {

	var payload injectormodels.ClearSidecarPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		ic.logger.Log().Error("Error binding JSON", zap.Error(err))

		c.JSON(http.StatusBadRequest, gin.H{"error": "Error binding JSON"})
		return
	}

	ic.logger.Log().Info("ClearSidecar - Received request", zap.Any("payload", payload))

	deployment, err := ic.kubeClient.ClearSidecar(&payload)

	if err != nil {
		ic.logger.Log().Error("Error clearing sidecar", zap.Error(err))

		c.JSON(http.StatusBadRequest, gin.H{"error": "Error clearing sidecar: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, deployment)
}
