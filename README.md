# Tracee PolicyReport Adapter

Install with Helm

```bash
helm repo add trivy-operator-polr-adapter https://fjogeleit.github.io/tracee-polr-adapter
helm install tracee-polr-adapter tracee-polr-adapter/tracee-polr-adapter -n tracee-adapter --create-namespace
```