#!/bin/bash

create_kind_cluster() {
  kind create cluster -n "$1"
}

install_kubefed_helm_chart() {
  kubectl config use-context "kind-$1"
  helm --namespace kube-federation-system upgrade -i kubefed kubefed-charts/kubefed --create-namespace
}

fix_networking_for_cluster() {
  IP_ADDR=$(docker inspect -f "{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}" "$1-control-plane")
  kubectl patch kubefedclusters -n kube-federation-system "$2" --type merge --patch "{\"spec\":{\"apiEndpoint\":\"https://${IP_ADDR}:6443\"}}"
}

connect_cluster() {
  kubefedctl join "$1" --cluster-context "$1" --host-cluster-context kind-cluster1 --v=2
}

# Create clusters
create_kind_cluster cluster1
create_kind_cluster cluster2
create_kind_cluster cluster3

# Add the kubefed helm chart
helm repo add kubefed-charts https://raw.githubusercontent.com/kubernetes-sigs/kubefed/master/charts

# Install the kubefed helm chart to each cluster
install_kubefed_helm_chart cluster1
install_kubefed_helm_chart cluster2
install_kubefed_helm_chart cluster3

# Fix networking issue on each of the clusters
fix_networking_issue cluster1 kind-cluster1
fix_networking_issue cluster2 kind-cluster2
fix_networking_issue cluster3 kind-cluster3

connect_cluster kind-cluster1
connect_cluster kind-cluster2
connect_cluster kind-cluster3
