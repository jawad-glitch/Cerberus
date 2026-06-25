#!/bin/bash
set -e

echo "Starting Cerberus Loacl Environment..."

kind delete cluster --name cerberus 2>/dev/null || true
kind create cluster --name cerberus
kubectl create namespace cerberus-system

helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo add neo4j https://helm.neo4j.com/neo4j
helm repo update

helm install prometheus prometheus-community/prometheus \
  --namespace cerberus-system \
  -f infra/prometheus-values.yaml

helm install neo4j neo4j/neo4j \
  --namespace cerberus-system \
  -f infra/neo4j-values.yaml

echo "Cerberus environment ready"
echo "Run: kubectl get pods -n cerberus-system"
kubectl apply -f operator/config/crd/bases/ml.cerberus.io_mlmodels.yaml
echo "✅ CRDs installed"
