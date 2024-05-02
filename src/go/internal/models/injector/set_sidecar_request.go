package injectormodels

type SetSidecarPayload struct {
	Namespace            string   `json:"Namespace" binding:"required"`
	DeploymentName       string   `json:"DeploymentName" binding:"required"`
	SidecarContainerName string   `json:"SidecarContainerName" binding:"required"`
	SidecarImage         string   `json:"SidecarImage" binding:"required"`
	Command              []string `json:"Command"`
	VolumeMounts         []Volume `json:"VolumeMounts"`
}

type Volume struct {
	Name      string `json:"Name" binding:"required"`
	MountPath string `json:"MountPath" binding:"required"`
}
