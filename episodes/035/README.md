# Episode 35: kubernetes Topology Aware routing and cilium 


**Join us on [YouTube](https://youtu.be/7Clr3rY02NQ) on Friday 4th Febuary, 7pm GMT**

## Headlines


### Resources:
[Docs](https://kubernetes.io/docs/concepts/services-networking/topology-aware-hints/)
[tests](https://github.com/kubernetes/kubernetes/commit/1dcf09c1bf899b3f5bf089e31942574e6663ceaf)
[test infra](https://github.com/kubernetes/test-infra/blob/master/experiment/kind-multizone-e2e.sh)

### Talks 


### Podcasts


### Articles 


### Tools

kind config.
`kind.yaml`
```yaml=
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
networking:
  disableDefaultCNI: true #disable the default cni.
  kubeProxyMode: none #do not deploy kube-proxy
featureGates:
  TopologyAwareHints: true
nodes:
- role: control-plane
- role: worker
  kubeadmConfigPatches:
  - |
    kind: JoinConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "topology.kubernetes.io/zone=a"
- role: worker
  kubeadmConfigPatches:
  - |
    kind: JoinConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "topology.kubernetes.io/zone=a"
- role: worker
  kubeadmConfigPatches:
  - |
    kind: JoinConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "topology.kubernetes.io/zone=b"
- role: worker
  kubeadmConfigPatches:
  - |
    kind: JoinConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "topology.kubernetes.io/zone=b"
- role: worker
  kubeadmConfigPatches:
  - |
    kind: JoinConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "topology.kubernetes.io/zone=c"
- role: worker
  kubeadmConfigPatches:
  - |
    kind: JoinConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "topology.kubernetes.io/zone=c"

```
kind start command:

```bash=
kind create cluster --config kind.yaml --image=kindest/node:v1.23.3 --name=cilium

```

cilium install
`cilium-install.sh`
```bash=
helm upgrade --install cilium cilium/cilium --version 1.11.1 \
    --namespace kube-system \
    --set kubeProxyReplacement=strict \
    --set k8sServiceHost="cilium-control-plane" \
    --set k8sServicePort=6443 \
    --set loadBalancer.serviceTopology=true \
    --set ipam.mode="kubernetes"
```

Create the example pods:
`pods.sh`
```bash=
#!/bin/bash

# The purpose of this script is to deploy to each node in the cluster 2 pods. 
# Each pod will have an env var that shows it's zone.

function echopod () {
  ZONE=""
  case $1 in
    cilium-worker)
      ZONE=a
      ;;
    cilium-worker2)
      ZONE=a
      ;;
    cilium-worker3)
      ZONE=b
      ;;
    cilium-worker4)
      ZONE=b
      ;;
    cilium-worker5)
      ZONE=c
      ;;
    cilium-worker6)
      ZONE=c
      ;;
  esac
  kubectl run echo${2}-${1} \
     --image overridden  --labels app=echo,pod=echo${2}-${1},node=${1},zone=$ZONE  --overrides \
    '{
      "spec":{
        "hostname": "echo'${2}-${1}'",
	      "subdomain": "test",
        "nodeName": "'$1'",
        "containers":[{
          "name":"echo",
          "image":"inanimate/echo-server",
          "env":[{
            "name":"ZONE",
            "value":"'$ZONE'"
          }]
        }]
      }
    }'
}


for worker in $(kind get nodes --name=cilium | grep worker)
  do 
    for i in {1..2}
      do echopod $worker $i
    done
  done

kubectl create service clusterip echo --tcp 8080 

echo "to mark the service eligible for topology aware hints."
echo "kubectl annotate svc echo service.kubernetes.io/topology-aware-hints=auto"
echo "to remove the annotation"
echo "kubectl annotate svc echo service.kubernetes.io/topology-aware-hints-"
echo "to check it's working"
echo "kubectl get endpointslices -o yaml | grep -i for"
```
