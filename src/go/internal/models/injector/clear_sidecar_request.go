package injectormodels

type ClearSidecarPayload struct {
	Namespace            string `json:"Namespace" binding:"required"`
	DeploymentName       string `json:"DeploymentName" binding:"required"`
	SidecarContainerName string `json:"SidecarContainerName" binding:"required"`
}
