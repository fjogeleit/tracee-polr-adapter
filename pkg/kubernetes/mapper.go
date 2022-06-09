package kubernetes

import (
	"crypto/sha1"
	"fmt"
	"strings"

	"github.com/fjogeleit/tracee-polr-adapter/pkg/tracee"
	"github.com/kyverno/kyverno/api/policyreport/v1alpha2"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

type Severity = int

const (
	informative Severity = iota
	low
	medium
	high
	critical
)

const (
	source            = "Tracee"
	clusterReportName = "tracee-cluster-policy-report"
)

var (
	reportLabels = map[string]string{
		"managed-by": "tracee-polr-adapter",
	}
)

func MapEvent(event tracee.Event) *v1alpha2.PolicyReportResult {
	var category, rule string

	if c, ok := event.SigMetadata.Properties[tracee.CategoryKey]; ok {
		if cString, ok := c.(string); ok {
			parts := strings.Split(cString, ":")
			category = strings.TrimSpace(parts[0])
			rule = strings.TrimSpace(parts[1])
		}

		delete(event.SigMetadata.Properties, tracee.CategoryKey)
	}

	props := map[string]string{
		"version":     event.SigMetadata.Version,
		"eventName":   event.Context.EventName,
		"eventID":     fmt.Sprint(event.Context.EventID),
		"processName": event.Context.ProcessName,
		"returnValue": fmt.Sprint(event.Context.ReturnValue),
		"resultID":    GeneratePolicyReportResultID(event.Context.ContainerID, event.Context.EventID, event.Context.Timestamp),
	}

	if event.Context.ContainerID != "" {
		props["containerID"] = event.Context.ContainerID
	}
	if event.Context.ContainerName != "" {
		props["containerName"] = event.Context.ContainerName
	}
	if event.Context.ContainerImage != "" {
		props["containerImage"] = event.Context.ContainerImage
	}

	for k, v := range event.SigMetadata.Properties {
		props[k] = fmt.Sprint(v)
	}

	var resources = make([]*corev1.ObjectReference, 0)
	if event.Context.PodName != "" {
		resources = append(resources, &corev1.ObjectReference{
			Namespace:  event.Context.PodNamespace,
			Name:       event.Context.PodName,
			UID:        types.UID(event.Context.PodUID),
			Kind:       "Pod",
			APIVersion: "v1",
		})
	}

	return &v1alpha2.PolicyReportResult{
		Category:   category,
		Policy:     event.SigMetadata.Name,
		Rule:       rule,
		Result:     MapResult(event.SigMetadata.Severity),
		Severity:   MapServerity(event.SigMetadata.Severity),
		Message:    event.SigMetadata.Description,
		Resources:  resources,
		Source:     source,
		Properties: props,
		Timestamp:  v1.Timestamp{Seconds: int64(event.Context.Timestamp)},
	}
}

func MapServerity(severity Severity) v1alpha2.PolicySeverity {
	if severity == informative {
		return ""
	} else if severity == low {
		return v1alpha2.SeverityLow
	} else if severity == medium {
		return v1alpha2.SeverityMedium
	}

	return v1alpha2.SeverityHigh
}

func MapResult(severity Severity) v1alpha2.PolicyResult {
	if severity == informative {
		return v1alpha2.StatusSkip
	} else if severity == low {
		return v1alpha2.StatusWarn
	} else if severity == medium {
		return v1alpha2.StatusWarn
	} else if severity == high {
		return v1alpha2.StatusFail
	}

	return v1alpha2.StatusError
}

func CreatePolicyReport(ns string) *v1alpha2.PolicyReport {
	return &v1alpha2.PolicyReport{
		ObjectMeta: v1.ObjectMeta{
			Name:      GeneratePolicyReportName(ns),
			Namespace: ns,
			Labels:    reportLabels,
		},
		Summary: v1alpha2.PolicyReportSummary{},
		Results: []*v1alpha2.PolicyReportResult{},
	}
}

func CreateClusterPolicyReport() *v1alpha2.ClusterPolicyReport {
	return &v1alpha2.ClusterPolicyReport{
		ObjectMeta: v1.ObjectMeta{
			Name:   clusterReportName,
			Labels: reportLabels,
		},
		Summary: v1alpha2.PolicyReportSummary{},
		Results: []*v1alpha2.PolicyReportResult{},
	}
}

func GeneratePolicyReportName(ns string) string {
	return fmt.Sprintf("tracee-polr-ns-%s", ns)
}

func GeneratePolicyReportResultID(containerID string, eventID int, timestamp int) string {
	id := fmt.Sprintf("%s_%d_%d", containerID, eventID, timestamp)

	h := sha1.New()
	h.Write([]byte(id))

	return fmt.Sprintf("%x", h.Sum(nil))
}
