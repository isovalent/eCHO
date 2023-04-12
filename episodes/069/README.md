# Episode 69: Evaluating cluster-api Distributions. 

[YouTube](https://youtu.be/oTCRZ-bt-Xo)

## Headlines

- [echo news 19](https://isogo.to/echo-news-19)

## Notes

[what even is cluster api](https://cluster-api.sigs.k8s.io/)
[what's kubeadm?](https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/create-cluster-kubeadm/)

[amazon eks anywhere](https://www.youtube.com/watch?v=yAoZi89CZ5o&list=PLDg_GiBbAx-mY3VFLPbLHcxo6wUjejAOC&index=47) 

[google anthos](https://cloud.google.com/anthos/clusters/docs/on-prem/latest/concepts/cluster-api)

[cluster-api-provider-xx](https://github.com/orgs/kubernetes-sigs/repositories?language=&q=cluster-api-provider&sort=&type=all)

[openshift](https://docs.openshift.com/container-platform/4.11/machine_management/capi-machine-management.html) 

[talos](https://github.com/siderolabs/sidero)

[cluster-api provider kubevirt](https://deploy-preview-7444--kubernetes-sigs-cluster-api.netlify.app/user/quick-start)


### Constraints.

I want to be able to programatically create and scale kubernetes clusters wherein each node has it's own kernel. So kind.sigs.k8s.io won't do it for me. 


### install script
```
# install cilium

cilium install

# install kubevip

kubectl apply -f https://raw.githubusercontent.com/kube-vip/kube-vip-cloud-provider/main/manifest/kube-vip-cloud-controller.yaml
kubectl create configmap --namespace kube-system kubevip --from-literal range-global=172.18.100.10-172.18.100.30

kubectl apply -f https://kube-vip.io/manifests/rbac.yaml
docker run --rm --net=host ghcr.io/kube-vip/kube-vip manifest daemonset --services --inCluster --arp --interface eth0 | kubectl apply -f -

# install kubevirt

# get KubeVirt version
KV_VER=$(curl "https://api.github.com/repos/kubevirt/kubevirt/releases/latest" | jq -r ".tag_name")
# deploy required CRDs
kubectl apply -f "https://github.com/kubevirt/kubevirt/releases/download/${KV_VER}/kubevirt-operator.yaml"
# deploy the KubeVirt custom resource
kubectl apply -f "https://github.com/kubevirt/kubevirt/releases/download/${KV_VER}/kubevirt-cr.yaml"
kubectl wait -n kubevirt kv kubevirt --for=condition=Available --timeout=10m


```