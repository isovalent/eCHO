# Episode 90: Cassandra and cilium clustermesh
[News](cilium.io/newsletter)

Background!

[Docs](https://docs.k8ssandra.io/install/local/multi-cluster-helm/)

``` bash
#!/usr/bin/env bash
#
# This script requires the following to be installed and available on your path:
#
#    - jq
#    - yq
#    - kustomize
#    - kind

set -e

getopt_version=$(getopt -V)
if [[ "$getopt_version" == " --" ]]; then
  echo "gnu-getopt doesn't seem to be installed. Install it using: brew install gnu-getopt"
  exit 1
fi

OPTS=$(getopt -o ho --long clusters:,cluster-prefix:,kind-node-version:,kind-worker-nodes:,output-file:,overwrite,help -n 'setup-kind-multicluster' -- "$@")
eval set -- "$OPTS"

default_kind_node_version=v1.25.3

function help() {
cat << EOF
Syntax: setup-kind-multicluster.sh [options]
Options:
  --clusters <clusters>          The number of clusters to create.
                                 Defaults to 1.
  --cluster-prefix <prefix>      The prefix to use to name clusters.
                                 Defaults to "k8ssandra-".
  --kind-node-version <version>  The image version of the kind nodes.
                                 Defaults to "$default_kind_node_version".
  --kind-worker-nodes <nodes>    The number of worker nodes to deploy.
                                 Can be a single number or a comma-separated list of numbers, one per cluster.
                                 Defaults to 3.
  --output-file <path>           The location of the file where the generated kubeconfig will be written to.
                                 Defaults to "./build/kind-kubeconfig". Existing content will be overwritten.
  -o|--overwrite                 Whether to delete existing clusters before re-creating them.
                                 Defaults to false.
  --help                         Displays this help message.
EOF
}

num_clusters=3
cluster_prefix="cluster-"
kind_node_version="$default_kind_node_version"
kind_worker_nodes=4
overwrite_clusters="no"
output_file="./build/kind-kubeconfig"
docker_network="172.18"
while true; do
  case "$1" in
    --clusters ) num_clusters=$2; shift 2 ;;
    --cluster-prefix ) cluster_prefix="$2"; shift 2 ;;
    --kind-node-version ) kind_node_version="$2"; shift 2 ;;
    --kind-worker-nodes ) kind_worker_nodes="$2"; shift 2 ;;
    --output-file ) output_file="$2"; shift 2 ;;
    -o | --overwrite) overwrite_clusters="yes"; shift ;;
    -h | --help ) help; exit;;
    -- ) shift; break ;;
    * ) break ;;
  esac
done


function create_cluster() {
  cluster_id=$1
  cluster_name=$2
  num_workers=$3
  node_version=$4

cat <<EOF | kind create cluster --name $cluster_name --image kindest/node:$node_version --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
networking:
  podSubnet: "10.20$cluster_id.0.0/16"
  disableDefaultCNI: true
nodes:
- role: control-plane
  extraPortMappings:
  - containerPort: 30080
    hostPort: 3${cluster_id}080
    protocol: TCP
  - containerPort: 30443
    hostPort: 3${cluster_id}443
    protocol: TCP
  - containerPort: 30942
    hostPort: 3${cluster_id}942
    protocol: TCP
$(for ((i=0; i<num_workers; i++)); do
cat << EOF2
- role: worker
  kubeadmConfigPatches:
    - |
      kind: JoinConfiguration
      nodeRegistration:
        kubeletExtraArgs:
          node-labels: "topology.kubernetes.io/zone=region$((${cluster_id}+1))-zone$(( (${i} % 3) +1))"
EOF2
done)
EOF

}

function install_cilium() {
  cluster_id=$1
  cluster_name=$2
  if [[ "$cluster_id" == "1" ]]; then
    cilium install --helm-set cluster.name="$cluster_name",cluster.id="$cluster_id",ingressController.enabled=true,ingressController.loadbalancerMode=dedicated,kubeProxyReplacement=partial --context "kind-$cluster_name"
    cilium hubble enable --ui --context "kind-$cluster_name"
  else
    cilium install --helm-set cluster.name="$cluster_name",cluster.id="$cluster_id",ingressController.enabled=true,ingressController.loadbalancerMode=dedicated,kubeProxyReplacement=partial --context "kind-$cluster_name" --inherit-ca "kind-${cluster_prefix}1"
  fi
  cilium clustermesh enable --service-type=LoadBalancer --context kind-$cluster_name
}


function cilium_clustermesh_connect() {
  for ((i=1; i<=num_clusters; i++))
    do for ((v=1; v<=num_clusters; v++))
      do if [[ $i -ne $v ]] ; then
        cilium clustermesh status --context kind-$cluster_prefix$i --wait
        cilium clustermesh connect --context kind-$cluster_prefix$i --destination-context kind-$cluster_prefix$v
      fi
    done
  done
}

function install_certmanager() {
  cluster_id=$1
  cluster_name=$2
  helm upgrade --install cert-manager -n cert-manager jetstack/cert-manager --set installCRDs=true --kube-context kind-$cluster_name --create-namespace
}

function install_k8ssandra_operator() {
  cluster_id=$1
  cluster_name=$2
  if [[ "cluster_id" == "1" ]] ; then
    helm upgrade --install k8ssandra-operator k8ssandra/k8ssandra-operator -n k8ssandra-operator --create-namespace --kube-context="kind-$cluster_name"
  else
    helm upgrade --install k8ssandra-operator k8ssandra/k8ssandra-operator -n k8ssandra-operator --create-namespace --kube-context="kind-$cluster_name" --set controlPlane=false
  fi
}

function install_kubevip() {
  cluster_id=$1
  cluster_name=$2
  kvversion="v0.5.12"
  kubectl --context kind-$cluster_name apply -f https://raw.githubusercontent.com/kube-vip/kube-vip/main/docs/manifests/rbac.yaml
  kubectl --context kind-$cluster_name apply -f https://raw.githubusercontent.com/kube-vip/kube-vip-cloud-provider/main/manifest/kube-vip-cloud-controller.yaml
  kubectl --context kind-$cluster_name create configmap --namespace kube-system kubevip --from-literal range-global=172.18.20$cluster_id.10-172.18.20$cluster_id.50
  cat <<EOF | kubectl --context kind-$cluster_name apply -f-
apiVersion: apps/v1
kind: DaemonSet
metadata:
  creationTimestamp: null
  labels:
    app.kubernetes.io/name: kube-vip-ds
    app.kubernetes.io/version: v0.5.12
  name: kube-vip-ds
  namespace: kube-system
spec:
  selector:
    matchLabels:
      app.kubernetes.io/name: kube-vip-ds
  template:
    metadata:
      creationTimestamp: null
      labels:
        app.kubernetes.io/name: kube-vip-ds
        app.kubernetes.io/version: v0.5.12
    spec:
      containers:
      - args:
        - manager
        env:
        - name: vip_arp
          value: "true"
        - name: port
          value: "6443"
        - name: vip_interface
          value: eth0
        - name: vip_cidr
          value: "32"
        - name: svc_enable
          value: "true"
        - name: vip_address
        - name: prometheus_server
          value: :2112
        image: ghcr.io/kube-vip/kube-vip:v0.5.12
        imagePullPolicy: Always
        name: kube-vip
        resources: {}
        securityContext:
          capabilities:
            add:
            - NET_ADMIN
            - NET_RAW
      hostNetwork: true
      serviceAccountName: kube-vip
  updateStrategy: {}
status:
  currentNumberScheduled: 0
  desiredNumberScheduled: 0
  numberMisscheduled: 0
  numberReady: 0
EOF
}

function delete_clusters() {
  echo "Deleting existing clusters..."

  for ((i=1; i<num_clusters; i++))
  do
    echo "Deleting cluster $((i+1)) out of $num_clusters"
    kind delete cluster --name "$cluster_prefix$i" || echo "Cluster $cluster_prefix$i doesn't exist yet"
  done
  echo
}

function create_clusters() {
  echo "Creating $num_clusters clusters..."

  for ((i=1; i<=num_clusters; i++))
  do
    echo "Creating cluster $i out of $num_clusters"
    if [[ "$kind_worker_nodes" == *,* ]]; then
      IFS=',' read -r -a nodes_array <<< "$kind_worker_nodes"
      nodes="${nodes_array[i]}"
    else
      nodes=$kind_worker_nodes
    fi
    create_cluster "$i" "$cluster_prefix$i" "$nodes" "$kind_node_version"
    install_cilium "$i" "$cluster_prefix$i"
    install_kubevip "$i" "$cluster_prefix$i"
    install_certmanager "$i" "$cluster_prefix$i"
    install_k8ssandra_operator "$i" "$cluster_prefix$i"
  done


}


# Creates a kubeconfig file that has entries for each of the clusters created.
# The file created is <project-root>/build/kind-kubeconfig and is intended for use
# primarily by tests running out of cluster.
function create_kubeconfig() {
  echo "Generating $output_file"

  temp_dir=$(mktemp -d)
  for ((i=0; i<num_clusters; i++))
  do
    kubeconfig_base="$temp_dir/$cluster_prefix$i.yaml"
    kind get kubeconfig --name "$cluster_prefix$i" --internal > "$kubeconfig_base"
  done

  basedir=$(dirname "$output_file")
  mkdir -p "$basedir"
  yq ea '. as $item ireduce({}; . *+ $item)' "$temp_dir"/*.yaml > "$output_file"
  # remove current-context for security
  yq e 'del(.current-context)' "$output_file" -i
}


if [[ "$overwrite_clusters" == "yes" ]]; then
  delete_clusters
fi
create_clusters
cilium_clustermesh_connect
kubectl config view --raw > ./build/kind-kubeconfig
#create_kubeconfig

# Set current context to the first cluster
kubectl config use "kind-${cluster_prefix}1"

```
