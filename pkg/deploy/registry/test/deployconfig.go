package test

import (
	"sync"

	"github.com/openshift/origin/pkg/deploy/api"
	metainternal "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/watch"
	apirequest "k8s.io/apiserver/pkg/endpoints/request"
)

type DeploymentConfigRegistry struct {
	Err               error
	DeploymentConfig  *api.DeploymentConfig
	DeploymentConfigs *api.DeploymentConfigList
	sync.Mutex
}

func NewDeploymentConfigRegistry() *DeploymentConfigRegistry {
	return &DeploymentConfigRegistry{}
}

func (r *DeploymentConfigRegistry) ListDeploymentConfigs(ctx apirequest.Context, label labels.Selector, field fields.Selector) (*api.DeploymentConfigList, error) {
	r.Lock()
	defer r.Unlock()

	return r.DeploymentConfigs, r.Err
}

func (r *DeploymentConfigRegistry) GetDeploymentConfig(ctx apirequest.Context, id string) (*api.DeploymentConfig, error) {
	r.Lock()
	defer r.Unlock()

	return r.DeploymentConfig, r.Err
}

func (r *DeploymentConfigRegistry) CreateDeploymentConfig(ctx apirequest.Context, image *api.DeploymentConfig) error {
	r.Lock()
	defer r.Unlock()

	r.DeploymentConfig = image
	return r.Err
}

func (r *DeploymentConfigRegistry) UpdateDeploymentConfig(ctx apirequest.Context, image *api.DeploymentConfig) error {
	r.Lock()
	defer r.Unlock()

	r.DeploymentConfig = image
	return r.Err
}

func (r *DeploymentConfigRegistry) DeleteDeploymentConfig(ctx apirequest.Context, id string) error {
	r.Lock()
	defer r.Unlock()

	return r.Err
}

func (r *DeploymentConfigRegistry) WatchDeploymentConfigs(ctx apirequest.Context, options *metainternal.ListOptions) (watch.Interface, error) {
	return nil, r.Err
}
