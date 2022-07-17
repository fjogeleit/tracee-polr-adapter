# Tracee PolicyReport Adapter

Can be used as webhook for [tracee](https://github.com/aquasecurity/tracee), to convert events into the unified PolicyReport and ClusterPolicyReport from the Kubernetes Policy Working Group. This makes it possible to use tooling like [Policy Reporter](https://github.com/kyverno/policy-reporter) for tracee as well.

## Pre Requirements

[PolicyReport CRDs](https://github.com/kubernetes-sigs/wg-policy-prototypes/tree/master/policy-report/crd/v1alpha2) are installed in your Cluster

## Installation via Helm

```bash
helm repo add tracee-polr-adapter https://fjogeleit.github.io/tracee-polr-adapter
helm install tracee-polr-adapter tracee-polr-adapter/tracee-polr-adapter -n tracee --create-namespace
```

## Configuration

### Results

* By default each report contains a maximum of 200 results, it can be changed by setting the value `results.maxPerReport` to any number > 0
* By default only violations with a severity > 0 (information) will processed, it can be changed by setting the value `results.minumumSeverity` to a number between 0 and 5

### Rules

It is possible to exclude rules by ID

```yaml (values.yaml)
rules:
    exclude: ["TCR-1"]
```

### Tracee

Currently the adapter requires the `v0.8.0` of `tracee` or above, with the experimental containers enrichment enabled, to have the Pod information into tracee events.

It is possible to also install __tracee__ with this Helm Chart with `--set tracee.enabled=true`. In this case the adapter will be preconfigured as tracee webhook.

Otherwise you need to configure the adapter as webhook yourself:

```yaml
...
    spec:
      containers:
        - name: tracee
          args:
            - --webhook http://tracee-polr-adapter:8080/webhook
            - --webhook-template ./templates/rawjson.tmpl
            - --webhook-content-type application/json
...
```

### Screenshots

![Policy Reporter UI - Tracee Screenshot](https://github.com/fjogeleit/tracee-polr-adapter/blob/main/images/tracee.png?raw=true)
