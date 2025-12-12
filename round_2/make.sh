#!/usr/bin/env bash
set -euo pipefail

POD_NAME="helm-auditor"
IMAGE_NAME="helm-auditor:latest"
YAML_FILE="k8s/helm-auditor.yaml"
NAMESPACE="default"

RESET_JOBS=false
RESET_PVC=false
DEBUG=false

# Parse flags
while [[ $# -gt 0 ]]; do
    case "$1" in
        --reset-jobs) RESET_JOBS=true ;;
        --reset-pvc)  RESET_PVC=true ;;
        --debug)      DEBUG=true ;;
        *) echo "Unknown flag $1"; exit 1 ;;
    esac
    shift
done

# Colors
green() { printf "\033[0;32m%s\033[0m\n" "$1"; }
yellow() { printf "\033[1;33m%s\033[0m\n" "$1"; }

# Reset jobs
if $RESET_JOBS; then
    yellow "==> Deleting old trivy/provenance jobs..."
    for j in $(kubectl get jobs | grep -E 'trivy|prov' | awk '{print $1}'); do kubectl delete job $j; done
fi

# Reset PVC
if $RESET_PVC; then
    yellow "==> Resetting reports PVC..."
    kubectl patch pvc reports-pvc -p '{"metadata":{"finalizers":[]}}' --type=merge || true
    kubectl delete pvc reports-pvc --grace-period=0 --force --wait=false || true

fi

# Delete old pod
yellow "==> Deleting old pod..."
kubectl delete pod "$POD_NAME" --ignore-not-found || true

# ConfigMap
yellow "==> Recreating ConfigMap..."
kubectl apply -f k8s/chart-config.yaml

# PVC & RBAC
yellow "==> Applying PVC..."
kubectl apply -f k8s/reports-pvc.yaml
yellow "==> Applying RBAC..."
kubectl apply -f k8s/trivy-rbac.yaml

# Build & load image
yellow "==> Building helm-auditor image..."
docker build --cache-from=type=local,src=.build-cache -t "$IMAGE_NAME" .

yellow "==> Loading image into Minikube..."
minikube image load "$IMAGE_NAME"

# Deploy pod
yellow "==> Applying helm-auditor pod..."
kubectl apply -f "$YAML_FILE"

# Wait pod ready
yellow "==> Waiting for pod to be ready..."
# kubectl wait --for=condition=Ready pod/"$POD_NAME" --timeout=240s
for i in {1..30}; do
    status=$(kubectl get pod "$POD_NAME" -o jsonpath='{.status.phase}')
    echo "Pod status: $status"
    [[ "$status" == "Running" ]] && break
    sleep 2
done

# Optional debug
if $DEBUG; then
    green "==> DEBUG MODE: Tailing all container logs"
    for c in runner aggregator auditor reporter; do
        echo "--- logs $c ---"
        kubectl logs "$POD_NAME" -c "$c"
    done
else
    green "==> Done. Use --debug to see logs."
fi

#echo "[+] Copiando reportes desde el PVC..."
#kubectl cp default/helm-auditor:/reports ./reports-local
#echo "[+] Reportes guardados en ./reports-local"
