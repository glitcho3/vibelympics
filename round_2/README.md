# Helm Auditor  
Security and provenance auditor for Helm charts and Kubernetes workloads

## Overview
Helm Auditor is a supply chain security tool that analyzes a Helm chart and produces a multi-layer audit report covering configuration risks, container image vulnerabilities, SBOM completeness, and provenance attestations.  
It is designed for auditing Helm charts as packages, providing insight into both the chart itself and the images it deploys. Auditing charts gives a broader perspective, from chart structure to container images and their dependencies.

Helm Auditor runs as a single container that orchestrates the full workflow: fetching the chart, templating manifests, scanning configurations and images, verifying signatures, collecting provenance data, and exporting structured reports.

## What it evaluates
Helm Auditor focuses on supply chain and security risk within a chart and the container images referenced by that chart. The audit includes:

### Chart and configuration analysis
- Misconfigurations using Trivy config scanning  
- Privilege escalation risks  
- Workload security posture  
- Unsafe defaults and insecure capabilities  

### Container image analysis
- CVEs from vulnerability databases  
- Severity breakdown and affected packages  
- Base image identification  
- Rebuild recommendation hints  

### Supply chain and provenance
- SBOM presence and completeness  
- Cosign signature verification  
- Attestation inspection  

### Chart level metadata
- Chart version and repository integrity  
- Chart source URL  

Reports are exported as JSON to a persistent volume for later inspection.

## Why Helm
A Helm chart functions as a package containing:
- Declarative configuration  
- Application topology  
- Container image references  
- Security relevant metadata  

This allows for auditing both the chart and its underlying images, giving a broader view of supply chain risks compared to auditing only PyPI or npm packages.

## Architecture
The container includes a pipeline of cooperating components:

1. **Fetcher** – Pulls the chart from an OCI repository.  
2. **Template engine** – Renders manifests to produce static YAML for analysis.  
3. **Trivy scanner** – Performs configuration analysis on Kubernetes manifests.  
4. **Image auditor** – Extracts image references and performs vulnerability scanning, SBOM extraction, and signature verification via Cosign.  
5. **Aggregator** – Normalizes all results into structured JSON reports.  
6. **Reporter** – Writes JSON output to the reports volume.

The Pod runs init containers for chart preparation and scanning plus a main container running the Helm Auditor binary.

## Prerequisites
- Docker  
- Minikube  

Start Minikube with systemd support and containerd runtime:
```bash
minikube start --driver=docker --force-systemd=true --container-runtime=containerd
```

## Using the `make.sh` helper
A `make.sh` script automates build, deployment, and optional reset operations. Usage:

```bash
./make.sh [--reset-jobs] [--reset-pvc] [--debug]
```

Flags:
- `--reset-jobs` – delete old Trivy or provenance jobs before starting  
- `--reset-pvc` – delete the reports PVC before starting  
- `--debug` – tail logs of all containers in the pod  

`make.sh` performs:
- Recreating ConfigMap  
- Applying PVC and RBAC manifests  
- Building and loading the Helm Auditor image  
- Deploying the pod  
- Waiting for pod readiness  

Reports remain in the PVC; you can copy them manually if needed.

## How to run manually

### 1. Build the container
```bash
docker build -t helm-auditor:latest .
```
### Building the Docker Image

This project uses **Docker Buildx** for building images. Buildx enables advanced features like multi-platform builds and efficient cache usage. To build the auditor image:

```bash
docker buildx build --load -t helm-auditor:latest .
```

### 2. Load image into Minikube
```bash
minikube image load helm-auditor:latest
```

### 3. Configure chart parameters
Edit `k8s/chart-config.yaml` and set:

```yaml
PROM_CHART: kube-prometheus-stack
PROM_REPO: oci://ghcr.io/prometheus-community/charts/
PROM_VERSION: 80.0.0
OUTPUT_FOLDER: "/reports/kube-prometheus-stack/"
```

### 4. Deploy the auditor Pod
```bash
kubectl apply -f k8s/helm-auditor.yaml
```

### 5. Wait for completion
```bash
kubectl logs -f helm-auditor
```

### 6. Retrieve reports
```bash
kubectl cp default/helm-auditor:/reports ./reports-local
```

Reports include configuration scans, image vulnerabilities, SBOM, attestations, and aggregated summaries in JSON.

## CLI usage
The container is also usable as a CLI tool:

```bash
docker run --rm \
  -e PROM_REPO="oci://ghcr.io/prometheus-community/charts/" \
  -e PROM_CHART="kube-prometheus-stack" \
  -e PROM_VERSION="80.0.0" \
  helm-auditor:latest
```

Reports will be printed or written to a mounted folder.

## Output format
Structured JSON with fields like:
- chart  
- images  
- CVE list  
- signature status  
- attestation type  
- SBOM presence  
- high risk configuration patterns  

## Purpose and advantages
Helm Auditor provides a systematic, automated approach to analyzing supply chain risks in Helm charts and container images.  
It gives actionable insights into misconfigurations, vulnerabilities, and provenance issues, helping teams ensure software integrity before deployment.

