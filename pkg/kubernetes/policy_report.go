package kubernetes

import (
	"fmt"

	"github.com/fjogeleit/tracee-polr-adapter/pkg/tracee"

	"github.com/kyverno/kyverno/api/policyreport/v1alpha2"
	pr "github.com/kyverno/kyverno/pkg/client/clientset/versioned/typed/policyreport/v1alpha2"
	"golang.org/x/net/context"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
)

type Client struct {
	k8sClient  pr.Wgpolicyk8sV1alpha2Interface
	maxResults int
}

func (p *Client) ProcessEvent(ctx context.Context, event tracee.Event) error {

	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		if event.Context.PodNamespace == "" {
			return p.handleClusterScoped(ctx, event)
		}

		return p.handleNamespaced(ctx, event, event.Context.PodNamespace)
	})
}

func (p *Client) handleNamespaced(ctx context.Context, event tracee.Event, ns string) error {
	polr, err := p.k8sClient.PolicyReports(ns).Get(ctx, GeneratePolicyReportName(ns), v1.GetOptions{})
	if errors.IsNotFound(err) {
		polr = CreatePolicyReport(ns)

		polr, err = p.k8sClient.PolicyReports(ns).Create(ctx, polr, v1.CreateOptions{})
		if err != nil {
			return fmt.Errorf("failed to create PolicyReport in namespace %s: %s", ns, err)
		}
	} else if err != nil {
		return err
	}

	polr.Results = p.updateResults(event, polr.Results, &polr.Summary)

	_, err = p.k8sClient.PolicyReports(ns).Update(ctx, polr, v1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update PolicyReport in namespace %s: %s", ns, err)
	}

	return nil
}

func (p *Client) handleClusterScoped(ctx context.Context, event tracee.Event) error {
	polr, err := p.k8sClient.ClusterPolicyReports().Get(ctx, clusterReportName, v1.GetOptions{})
	if errors.IsNotFound(err) {
		polr = CreateClusterPolicyReport()

		polr, err = p.k8sClient.ClusterPolicyReports().Create(ctx, polr, v1.CreateOptions{})
		if err != nil {
			return fmt.Errorf("failed to create ClusterPolicyReport: %s", err)
		}
	} else if err != nil {
		return err
	}

	polr.Results = p.updateResults(event, polr.Results, &polr.Summary)

	_, err = p.k8sClient.ClusterPolicyReports().Update(ctx, polr, v1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update ClusterPolicyReport: %s", err)
	}

	return nil
}

func (p *Client) updateResults(event tracee.Event, results []*v1alpha2.PolicyReportResult, summary *v1alpha2.PolicyReportSummary) []*v1alpha2.PolicyReportResult {
	if p.maxResults > 0 && len(results) >= p.maxResults {
		index := len(results) - p.maxResults
		removed := results[index]

		removeResultFromSummary(summary, removed)
		results = results[index+1:]
	}

	result := MapEvent(event)
	addResultToSummary(summary, result)

	return append(results, result)
}

func removeResultFromSummary(sum *v1alpha2.PolicyReportSummary, result *v1alpha2.PolicyReportResult) {
	if result.Result == v1alpha2.StatusSkip {
		sum.Skip--
	}
	if result.Result == v1alpha2.StatusPass {
		sum.Pass--
	}
	if result.Result == v1alpha2.StatusWarn {
		sum.Warn--
	}
	if result.Result == v1alpha2.StatusFail {
		sum.Fail--
	}
	if result.Result == v1alpha2.StatusError {
		sum.Error--
	}
}

func addResultToSummary(sum *v1alpha2.PolicyReportSummary, result *v1alpha2.PolicyReportResult) {
	if result.Result == v1alpha2.StatusSkip {
		sum.Skip++
	}
	if result.Result == v1alpha2.StatusPass {
		sum.Pass++
	}
	if result.Result == v1alpha2.StatusWarn {
		sum.Warn++
	}
	if result.Result == v1alpha2.StatusFail {
		sum.Fail++
	}
	if result.Result == v1alpha2.StatusError {
		sum.Error++
	}
}

func NewPolicyReportClient(client pr.Wgpolicyk8sV1alpha2Interface, maxResults int) *Client {
	return &Client{
		k8sClient:  client,
		maxResults: maxResults,
	}
}
