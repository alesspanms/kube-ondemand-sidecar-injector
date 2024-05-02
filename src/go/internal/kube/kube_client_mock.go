package kube

import (
	"github.com/stretchr/testify/mock"

	injectormodels "github.com/alesspanms/kube-ondemand-sidecar-injector/internal/models/injector"
)

type KubeClientMock struct {
	mock.Mock
}

func (m *KubeClientMock) GetDeployments(namespace string) ([]injectormodels.Deployment, error) {
	args := m.Called(namespace)
	res := args.Get(0)
	if res == nil {
		return nil, args.Error(1)
	}
	return res.([]injectormodels.Deployment), args.Error(1)
}

func (m *KubeClientMock) GetSingleDeployment(namespace string, name string) (injectormodels.Deployment, error) {
	args := m.Called(namespace, name)
	res := args.Get(0)
	if res == nil {
		return injectormodels.Deployment{}, args.Error(1)
	}
	return res.(injectormodels.Deployment), args.Error(1)
}

func (m *KubeClientMock) SetSidecar(payload *injectormodels.SetSidecarPayload) (injectormodels.Deployment, error) {
	args := m.Called(payload)
	res := args.Get(0)
	if res == nil {
		return injectormodels.Deployment{}, args.Error(1)
	}
	return args.Get(0).(injectormodels.Deployment), args.Error(1)
}

func (m *KubeClientMock) ClearSidecar(payload *injectormodels.ClearSidecarPayload) (injectormodels.Deployment, error) {
	args := m.Called(payload)
	res := args.Get(0)
	if res == nil {
		return injectormodels.Deployment{}, args.Error(1)
	}
	return args.Get(0).(injectormodels.Deployment), args.Error(1)
}
