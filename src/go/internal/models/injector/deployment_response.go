package injectormodels

type Deployment struct {
	Namespace   string   `json:"Namespace"`
	Name        string   `json:"Name"`
	VolumeNames []string `json:"VolumeNames"`
}

// list of deployment
// type ListDeployments struct {
// 	Deployments []Deployment `json:"Deployments"`
// }
