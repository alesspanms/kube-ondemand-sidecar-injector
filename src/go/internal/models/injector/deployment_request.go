package injectormodels

type GetDeploymentsPayload struct {
	Namespace                      string `json:"Namespace" binding:"required"`
	Filtered                       bool   `json:"Filtered"`
	DeploymentNameSubstringPattern string `json:"DeploymentNameSubstringPattern"`
}

type GetSingleDeploymentPayload struct {
	Namespace      string `json:"Namespace" binding:"required"`
	DeploymentName string `json:"DeploymentName" binding:"required"`
}
