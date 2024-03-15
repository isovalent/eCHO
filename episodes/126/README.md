# Episode 126: Ideas on Migrating Calico to Cilium with Network Policy.

with Duffie Cooley

## News

## Session


:::spoiler cluster api bring up



```bash
#create a kind cluster with no cni installed.
cat <<EOF | kind create cluster --name=capi --config -
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
containerdConfigPatches:
- |-
  [plugins."io.containerd.grpc.v1.cri".registry.mirrors."docker.io"]
    endpoint = ["http://docker-cache:5000"]
- |-
  [plugins."io.containerd.grpc.v1.cri".registry.mirrors."quay.io"]
    endpoint = ["http://quay-cache:5000"]
- |-
  [plugins."io.containerd.grpc.v1.cri".registry.mirrors."ghcr.io"]
    endpoint = ["http://ghcr-cache:5000"]
networking:
  kubeProxyMode: "none"
  disableDefaultCNI: true
  serviceSubnet: "10.96.0.0/16"
nodes:
- role: control-plane
  extraMounts:
    - hostPath: /var/run/docker.sock
      containerPath: /var/run/docker.sock
EOF

# install cilium and hubble
cilium install --set l2announcements.enabled=true
cilium hubble enable --ui
#virtink needs certmanager so we will init clusterctl first
export CLUSTER_TOPOLOGY=true
export EXP_MACHINE_POOL=true

clusterctl init --infrastructure docker

# define load balancing for cilium
cat <<EOF | kubectl apply -f -
---
apiVersion: "cilium.io/v2alpha1"
kind: CiliumL2AnnouncementPolicy
metadata:
  name: l2-default
spec:
  serviceSelector:
    matchLabels: {} # match all
  interfaces:
  - ^eth[0-9]+
  externalIPs: false
  loadBalancerIPs: true
---
apiVersion: "cilium.io/v2alpha1"
kind: CiliumLoadBalancerIPPool
metadata:
  name: "l2-default"
spec:
  cidrs:
  - cidr: "172.18.200.0/24"
EOF
```
:::

:::spoiler cluster definition
```yaml
apiVersion: cluster.x-k8s.io/v1beta1
kind: ClusterClass
metadata:
  name: quick-start
  namespace: c1
spec:
  controlPlane:
    machineInfrastructure:
      ref:
        apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
        kind: DockerMachineTemplate
        name: quick-start-control-plane
    ref:
      apiVersion: controlplane.cluster.x-k8s.io/v1beta1
      kind: KubeadmControlPlaneTemplate
      name: quick-start-control-plane
  infrastructure:
    ref:
      apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
      kind: DockerClusterTemplate
      name: quick-start-cluster
  patches:
  - definitions:
    - jsonPatches:
      - op: add
        path: /spec/template/spec/kubeadmConfigSpec/clusterConfiguration/imageRepository
        valueFrom:
          variable: imageRepository
      selector:
        apiVersion: controlplane.cluster.x-k8s.io/v1beta1
        kind: KubeadmControlPlaneTemplate
        matchResources:
          controlPlane: true
    description: Sets the imageRepository used for the KubeadmControlPlane.
    enabledIf: '{{ ne .imageRepository "" }}'
    name: imageRepository
  - definitions:
    - jsonPatches:
      - op: add
        path: /spec/template/spec/kubeadmConfigSpec/clusterConfiguration/etcd
        valueFrom:
          template: |
            local:
              imageTag: {{ .etcdImageTag }}
      selector:
        apiVersion: controlplane.cluster.x-k8s.io/v1beta1
        kind: KubeadmControlPlaneTemplate
        matchResources:
          controlPlane: true
    description: Sets tag to use for the etcd image in the KubeadmControlPlane.
    name: etcdImageTag
  - definitions:
    - jsonPatches:
      - op: add
        path: /spec/template/spec/kubeadmConfigSpec/clusterConfiguration/dns
        valueFrom:
          template: |
            imageTag: {{ .coreDNSImageTag }}
      selector:
        apiVersion: controlplane.cluster.x-k8s.io/v1beta1
        kind: KubeadmControlPlaneTemplate
        matchResources:
          controlPlane: true
    description: Sets tag to use for the etcd image in the KubeadmControlPlane.
    name: coreDNSImageTag
  - definitions:
    - jsonPatches:
      - op: add
        path: /spec/template/spec/customImage
        valueFrom:
          template: |
            kindest/node:{{ .builtin.machineDeployment.version | replace "+" "_" }}
      selector:
        apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
        kind: DockerMachineTemplate
        matchResources:
          machineDeploymentClass:
            names:
            - default-worker
    - jsonPatches:
      - op: add
        path: /spec/template/spec/template/customImage
        valueFrom:
          template: |
            kindest/node:{{ .builtin.machinePool.version | replace "+" "_" }}
      selector:
        apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
        kind: DockerMachinePoolTemplate
        matchResources:
          machinePoolClass:
            names:
            - default-worker
    - jsonPatches:
      - op: add
        path: /spec/template/spec/customImage
        valueFrom:
          template: |
            kindest/node:{{ .builtin.controlPlane.version | replace "+" "_" }}
      selector:
        apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
        kind: DockerMachineTemplate
        matchResources:
          controlPlane: true
    description: Sets the container image that is used for running dockerMachines
      for the controlPlane and default-worker machineDeployments.
    name: customImage
  - definitions:
    - jsonPatches:
      - op: add
        path: /spec/template/spec/kubeadmConfigSpec/clusterConfiguration/apiServer/extraArgs
        value:
          admission-control-config-file: /etc/kubernetes/kube-apiserver-admission-pss.yaml
      - op: add
        path: /spec/template/spec/kubeadmConfigSpec/clusterConfiguration/apiServer/extraVolumes
        value:
        - hostPath: /etc/kubernetes/kube-apiserver-admission-pss.yaml
          mountPath: /etc/kubernetes/kube-apiserver-admission-pss.yaml
          name: admission-pss
          pathType: File
          readOnly: true
      - op: add
        path: /spec/template/spec/kubeadmConfigSpec/files
        valueFrom:
          template: |
            - content: |
                apiVersion: apiserver.config.k8s.io/v1
                kind: AdmissionConfiguration
                plugins:
                - name: PodSecurity
                  configuration:
                    apiVersion: pod-security.admission.config.k8s.io/v1{{ if semverCompare "< v1.25" .builtin.controlPlane.version }}beta1{{ end }}
                    kind: PodSecurityConfiguration
                    defaults:
                      enforce: "{{ .podSecurityStandard.enforce }}"
                      enforce-version: "latest"
                      audit: "{{ .podSecurityStandard.audit }}"
                      audit-version: "latest"
                      warn: "{{ .podSecurityStandard.warn }}"
                      warn-version: "latest"
                    exemptions:
                      usernames: []
                      runtimeClasses: []
                      namespaces: [kube-system]
              path: /etc/kubernetes/kube-apiserver-admission-pss.yaml
      selector:
        apiVersion: controlplane.cluster.x-k8s.io/v1beta1
        kind: KubeadmControlPlaneTemplate
        matchResources:
          controlPlane: true
    description: Adds an admission configuration for PodSecurity to the kube-apiserver.
    enabledIf: '{{ .podSecurityStandard.enabled }}'
    name: podSecurityStandard
  variables:
  - name: imageRepository
    required: true
    schema:
      openAPIV3Schema:
        default: ""
        description: imageRepository sets the container registry to pull images from.
          If empty, nothing will be set and the from of kubeadm will be used.
        example: registry.k8s.io
        type: string
  - name: etcdImageTag
    required: true
    schema:
      openAPIV3Schema:
        default: ""
        description: etcdImageTag sets the tag for the etcd image.
        example: 3.5.3-0
        type: string
  - name: coreDNSImageTag
    required: true
    schema:
      openAPIV3Schema:
        default: ""
        description: coreDNSImageTag sets the tag for the coreDNS image.
        example: v1.8.5
        type: string
  - name: podSecurityStandard
    required: false
    schema:
      openAPIV3Schema:
        properties:
          audit:
            default: restricted
            description: audit sets the level for the audit PodSecurityConfiguration
              mode. One of privileged, baseline, restricted.
            type: string
          enabled:
            default: false
            description: enabled enables the patches to enable Pod Security Standard
              via AdmissionConfiguration.
            type: boolean
          enforce:
            default: baseline
            description: enforce sets the level for the enforce PodSecurityConfiguration
              mode. One of privileged, baseline, restricted.
            type: string
          warn:
            default: restricted
            description: warn sets the level for the warn PodSecurityConfiguration
              mode. One of privileged, baseline, restricted.
            type: string
        type: object
  workers:
    machineDeployments:
    - class: default-worker
      template:
        bootstrap:
          ref:
            apiVersion: bootstrap.cluster.x-k8s.io/v1beta1
            kind: KubeadmConfigTemplate
            name: quick-start-default-worker-bootstraptemplate
        infrastructure:
          ref:
            apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
            kind: DockerMachineTemplate
            name: quick-start-default-worker-machinetemplate
    machinePools:
    - class: default-worker
      template:
        bootstrap:
          ref:
            apiVersion: bootstrap.cluster.x-k8s.io/v1beta1
            kind: KubeadmConfigTemplate
            name: quick-start-default-worker-bootstraptemplate
        infrastructure:
          ref:
            apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
            kind: DockerMachinePoolTemplate
            name: quick-start-default-worker-machinepooltemplate
---
apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
kind: DockerClusterTemplate
metadata:
  name: quick-start-cluster
  namespace: c1
spec:
  template:
    spec: {}
---
apiVersion: controlplane.cluster.x-k8s.io/v1beta1
kind: KubeadmControlPlaneTemplate
metadata:
  name: quick-start-control-plane
  namespace: c1
spec:
  template:
    spec:
      kubeadmConfigSpec:
        clusterConfiguration:
          apiServer:
            certSANs:
            - localhost
            - 127.0.0.1
            - 0.0.0.0
            - host.docker.internal
          controllerManager:
            extraArgs:
              enable-hostpath-provisioner: "true"
        initConfiguration:
          nodeRegistration: {}
        joinConfiguration:
          nodeRegistration: {}
---
apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
kind: DockerMachineTemplate
metadata:
  name: quick-start-control-plane
  namespace: c1
spec:
  template:
    spec:
      extraMounts:
      - containerPath: /var/run/docker.sock
        hostPath: /var/run/docker.sock
---
apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
kind: DockerMachineTemplate
metadata:
  name: quick-start-default-worker-machinetemplate
  namespace: c1
spec:
  template:
    spec:
      extraMounts:
      - containerPath: /var/run/docker.sock
        hostPath: /var/run/docker.sock
---
apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
kind: DockerMachinePoolTemplate
metadata:
  name: quick-start-default-worker-machinepooltemplate
  namespace: c1
spec:
  template:
    spec:
      template: {}
---
apiVersion: bootstrap.cluster.x-k8s.io/v1beta1
kind: KubeadmConfigTemplate
metadata:
  name: quick-start-default-worker-bootstraptemplate
  namespace: c1
spec:
  template:
    spec:
      joinConfiguration:
        nodeRegistration: {}
---
apiVersion: cluster.x-k8s.io/v1beta1
kind: Cluster
metadata:
  name: c1
  namespace: c1
spec:
  clusterNetwork:
    pods:
      cidrBlocks:
      - 10.0.0.0/8
    serviceDomain: cluster.local
    services:
      cidrBlocks:
      - 192.168.0.0/16
  topology:
    class: quick-start
    controlPlane:
      metadata: {}
      replicas: 3
    variables:
    - name: imageRepository
      value: ""
    - name: etcdImageTag
      value: ""
    - name: coreDNSImageTag
      value: ""
    - name: podSecurityStandard
      value:
        audit: restricted
        enabled: false
        enforce: baseline
        warn: restricted
    version: v1.29.0
    workers:
      machineDeployments:
      - class: default-worker
        name: calico
        replicas: 3
      - class: default-worker
        name: cilium
        replicas: 0
```
:::

:::spoiler calico custom resource
```yaml
# This section includes base Calico installation configuration.
# For more information, see: https://docs.tigera.io/calico/latest/reference/installation/api#operator.tigera.io/v1.Installation
apiVersion: operator.tigera.io/v1
kind: Installation
metadata:
  name: default
spec:
  # Configures Calico networking.
  calicoNetwork:
    # Note: The ipPools section cannot be modified post-install.
    ipPools:
    - blockSize: 26
      cidr: 10.100.0.0/16
      encapsulation: VXLAN
      natOutgoing: Enabled
      nodeSelector: cni=="calico"
  controlPlaneNodeSelector:
    cni: "calico"
  calicoNodeDaemonSet:
    spec:
      template:
        spec:
          nodeSelector:
            cni: "calico"
  csiNodeDriverDaemonSet:
    spec:
      template:
        spec:
          nodeSelector:
            cni: "calico"
  calicoKubeControllersDeployment:
    spec:
      template:
        spec:
          nodeSelector:
            cni: "calico"
  typhaDeployment:
    spec:
      template:
        spec:
          nodeSelector:
            cni: "calico"


```
:::

::: spoiler images to preload
```yaml
calico/cni:v3.27.2
calico/kube-controllers:v3.27.2
calico/node-driver-registrar:v3.27.2
calico/node:v3.27.2
calico/pod2daemon-flexvol:v3.27.2
calico/typha:v3.27.2
calico/csi:v3.27.2
mauilion/goldpinger:v3.5.1

```
:::
