## Kubernetes Clusters

Either k3s or a regular kubernetes cluster installed via kubeadm can be used.

Both clusters should have one node which should be on a local network. Both should be able to access each other over the network.

### Submariner
Submariner needs to be installed on both clusters. Following [the documentation](https://github.com/shaejaz/submariner-website/blob/faee030c2025150dd607796c52cd633442585cc9/src/content/getting-started/quickstart/k3s/_index.md) should be sufficient to have a running setup.

Note: Submariner has removed official support for k3s due to insufficient dev resources. The documnet above should work for the current Submariner version. Latest install guides can be used however, if k3s is not being used: [https://submariner.io/getting-started/quickstart/](https://submariner.io/getting-started/quickstart/)

### Operator

Once Submariner is working on both clusters, the operator needs to be installed on both as well. On each node run the following commands:

```
cd operator
make manifest
make install

make run
```

This would get the operator running in the clusters. It should start listening to nay new `DependencyList` resources being created in their respective clusters.

Do note, the `DependencyList` resources should be the same across the clusters in order to avoid untested behavior.
