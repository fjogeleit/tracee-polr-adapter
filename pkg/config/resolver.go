package config

import (
	"github.com/fjogeleit/tracee-polr-adapter/pkg/api"
	k8s "github.com/fjogeleit/tracee-polr-adapter/pkg/kubernetes"
	"github.com/fjogeleit/tracee-polr-adapter/pkg/tracee"

	"github.com/kyverno/kyverno/pkg/client/clientset/versioned/typed/policyreport/v1alpha2"
	"k8s.io/client-go/rest"
)

// Resolver manages dependencies
type Resolver struct {
	config     *Config
	k8sConfig  *rest.Config
	polrClient *k8s.Client
}

// PolicyReportClient resolver method
func (r *Resolver) PolicyReportClient() (*k8s.Client, error) {
	if r.polrClient != nil {
		return r.polrClient, nil
	}

	client, err := v1alpha2.NewForConfig(r.k8sConfig)
	if err != nil {
		return nil, err
	}

	policyreportClient := k8s.NewPolicyReportClient(
		client,
		r.config.Results.MaxPerReport,
	)

	r.polrClient = policyreportClient

	return policyreportClient, nil
}

// Server resolver method
func (r *Resolver) Server() (api.Server, error) {
	client, err := r.PolicyReportClient()
	if err != nil {
		return nil, err
	}
	return api.NewServer(r.config.Webhook.Port, client, r.Filter()), nil
}

// Filter resolver method
func (r *Resolver) Filter() *tracee.Filter {
	return tracee.NewFilter(r.config.Results.MinimumSeverity)
}

// NewResolver constructor function
func NewResolver(config *Config, k8sConfig *rest.Config) Resolver {
	return Resolver{
		config:    config,
		k8sConfig: k8sConfig,
	}
}
