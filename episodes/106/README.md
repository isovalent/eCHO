With [Duffie Cooley](https://twitter.com/mauilion)
{%youtube kurMo3r4Ol4 %}

## Headlines

* [eBPF Summit 2023](https://ebpf.io/summit-2023/)
* [Echo Newsletter](https://cilium.io/newsletter)


# Live Migration of an AWS EKS cluster. 
## Migrating to Cilium

This document covers a few different strategies for migrating to cilium in each cloud environment. The intent of this document is to bring up and discuss the different solutions and to choose one to focus on for each cloud environment so that Isovalent can help support the migration of clusters to cilium enterprise in these environments.

## High level Overview
The CNI layer is a very low level component in an operating Kubernetes cluster. Most CNI implementations take care of managing IPAM (IP Address Management) as well as the lifecycle of network namespaces and interfaces associated with pods. A good number of them also support network policy to some extent. These different aspects mean that any migration to another CNI. Should consider each of these elements to ensure the migration is successful. We often hear from customers that they would prefer a hitless migration and that all pods during the migration need to allow connectivity during the period of migration. It's ok to restart pods on a scheduled and managed basis but it's not okay to restart all pods on all nodes all at once.

When considering any migration or any other changes with regard to CNI implementations. It's important to ensure that you have a well defined and tested strategy for the migration. This includes all of the necessary steps to perform the migration according to your migration strategy as well as how to pause, roll back or test each step along the way. This is a considerable effort to get right and should require testing and determine success and failure metrics throughout the process.

### Cilium Migration Tooling.
#### Install
When installed by default Cilium will override the existing CNI in a cluster with itself. It does this by placing a file in `/etc/cni/net.d/05-cilium.conflist` and the necessary `cilium-cni` binary in `/opt/cni/bin/`.

This feature can be configured with the helm values:
:::spoiler
```
cni:
  # -- Install the CNI configuration and binary files into the filesystem.
  install: true

  # -- Remove the CNI configuration and binary files on agent shutdown. Enable this
  # if you're removing Cilium from the cluster. Disable this to prevent the CNI
  # configuration file from being removed during agent upgrade, which can cause
  # nodes to go unmanageable.
  uninstall: false

  # -- Configure chaining on top of other CNI plugins. Possible values:
  #  - none
  #  - aws-cni
  #  - flannel
  #  - generic-veth
  #  - portmap
  chainingMode: ~

  # -- A CNI network name in to which the Cilium plugin should be added as a chained plugin.
  # This will cause the agent to watch for a CNI network with this network name. When it is
  # found, this will be used as the basis for Cilium's CNI configuration file. If this is
  # set, it assumes a chaining mode of generic-veth. As a special case, a chaining mode
  # of aws-cni implies a chainingTarget of aws-cni.
  chainingTarget: ~

  # -- Make Cilium take ownership over the `/etc/cni/net.d` directory on the
  # node, renaming all non-Cilium CNI configurations to `*.cilium_bak`.
  # This ensures no Pods can be scheduled using other CNI plugins during Cilium
  # agent downtime.
  exclusive: true

  # -- Configure the log file for CNI logging with retention policy of 7 days.
  # Disable CNI file logging by setting this field to empty explicitly.
  logFile: /var/run/cilium/cilium-cni.log

  # -- Skip writing of the CNI configuration. This can be used if
  # writing of the CNI configuration is performed by external automation.
  customConf: false

  # -- Configure the path to the CNI configuration directory on the host.
  confPath: /etc/cni/net.d

  # -- Configure the path to the CNI binary directory on the host.
  binPath: /opt/cni/bin

  # -- Specify the path to a CNI config to read from on agent start.
  # This can be useful if you want to manage your CNI
  # configuration outside of a Kubernetes environment. This parameter is
  # mutually exclusive with the 'cni.configMap' parameter. The agent will
  # write this to 05-cilium.conflist on startup.
  # readCniConf: /host/etc/cni/net.d/05-sample.conflist.input

  # -- When defined, configMap will mount the provided value as ConfigMap and
  # interpret the cniConf variable as CNI configuration file and write it
  # when the agent starts up
  # configMap: cni-configuration

  # -- Configure the key in the CNI ConfigMap to read the contents of
  # the CNI configuration from.
  configMapKey: cni-config

  # -- Configure the path to where to mount the ConfigMap inside the agent pod.
  confFileMountPath: /tmp/cni-configuration

  # -- Configure the path to where the CNI configuration directory is mounted
  # inside the agent pod.
  hostConfDirMountPath: /host/etc/cni/net.d

```
:::

#### Operator

The cilium operator is deployed in addition to the cilium daemonset. It has a feature that is enabled by default that will restart all pods that are not managed by cilium on a schedule. When migrating to cilium it's recommended that we disable this feature during the migration so that the operator doesn't attempt to restart pods until the migration step supports it and then to reconfigure the operator enable it specifically. 

Not disabling this feature will cause all pods on the cluster to cycle repeatedly until they are managed by cilium. This action is taken from the operator pod as part of the operator deployment. It is not agent specific. 

Controlled by the following helm value:
:::spoiler
```
operator:
  unmanagedPodWatcher:
    # -- Restart any pod that are not managed by Cilium.
    restart: true
    # -- Interval, in seconds, to check if there are any pods that are not
    # managed by Cilium.
    intervalSeconds: 15

```
:::
#### IP Address Management
This is implemented differently in each of cloud environments. The simplest is host local ipam. This indicates that the cni provider is allocating ip addresses from a subnet allocated to the node. In Cilium we refer to this as `ipam.mode=Kubernetes` this is where the cni will determine from the annotation on the node the subnet to allocate ip addresses to pods and handle all of that locally.

Cilium by default will use ipam.mode cluster-pool. This will allocate a pod cidr per node from the default `clusterPoolIPv4CIDRList` `["10.0.0.0/8"]`

In the case of migration it's important to determine how ipam is being managed and how to configure cilium to manage it in a consistent way. It's likely that you would want to change the default ipam configuration to align with how the existing cni is managing ip addresses. 

In some cases when migrating in place we can use cilium ipam in cluster pool mode to support pods on the cilium network and pods on the existing cni network. As we are effectively standing up two routable network on the underlying host and migrating pods in place. For more information about this take a look at the migration strategy section. 

You can learn more about the different IPAM configurations in the different environments [here](https://docs.cilium.io/en/stable/network/concepts/ipam/)

The following helm chart values manage this configuration: 
:::spoiler
```
ipam:
  # -- Configure IP Address Management mode.
  # ref: https://docs.cilium.io/en/stable/network/concepts/ipam/
  mode: "cluster-pool"
  # -- Maximum rate at which the CiliumNode custom resource is updated.
  ciliumNodeUpdateRate: "15s"
  operator:
    # -- IPv4 CIDR list range to delegate to individual nodes for IPAM.
    clusterPoolIPv4PodCIDRList: ["10.0.0.0/8"]
    # -- IPv4 CIDR mask size to delegate to individual nodes for IPAM.
    clusterPoolIPv4MaskSize: 24
    # -- IPv6 CIDR list range to delegate to individual nodes for IPAM.
    clusterPoolIPv6PodCIDRList: ["fd00::/104"]
    # -- IPv6 CIDR mask size to delegate to individual nodes for IPAM.
    clusterPoolIPv6MaskSize: 120
    # -- IP pools to auto-create in multi-pool IPAM mode.
    autoCreateCiliumPodIPPools: {}
      # default:
      #   ipv4:
      #     cidrs:
      #       - 10.10.0.0/8
      #     maskSize: 24
      # other:
      #   ipv6:
      #     cidrs:
      #       - fd00:100::/80
      #     maskSize: 96
    # -- The maximum burst size when rate limiting access to external APIs.
    # Also known as the token bucket capacity.
    # @default -- `20`
    externalAPILimitBurstSize: ~
    # -- The maximum queries per second when rate limiting access to
    # external APIs. Also known as the bucket refill rate, which is used to
    # refill the bucket up to the burst size capacity.
    # @default -- `4.0`
    externalAPILimitQPS: ~
```
:::
#### Per Node Config
Cilium as of 1.13 has the ability to manage cilium cni configuration on a per node basis. This allows for overriding the cilium config configmap on a per node basis. This can be helpful in cases where we are migrating in place. With this mechanism we can override specific configurations of cilium by labeling nodes that we want the config overlay to apply to. Read more about this [here](https://docs.cilium.io/en/stable/configuration/per-node-config)

## Migration Strategies
Depending on the environment and the existing CNI implementation there are a number of general strategies for handling a CNI migration. Each has it's advantages. 


### Constraints

There are a number of constraints to consider when evaluating a migration strategy. In some cases it's possible to spin up new clusters with the new CNI and migrate pods to the new cluster. When considering stateful pods this can be challenging. There are tools like [Velero](https://velero.io) that can backup a set of resources including their state and restore them to a different cluster. These tools can be very useful in migrating pods where you have the time to take the pod down and bring it up on new infrastructure. 



#### Downtime.
Downtime refers to a state that will cause pods to become unreachable during the period of migration. At a high level downtime is a normal part of the lifecycle of pods in a Kubernetes evironment. When Changing the underlying CNI these pods need to be restarted so that they can be configured by the new CNI.

#### Connectivity.
Connectivity between pods on different CNI implementations is a complex topic. In some specific cases it's possible to support pod to pod communication across two CNI implementations. This is usually possible when both CNI implementations meet at the routing layer. In the case where some nodes are managed by Cilium and other nodes are managed by another CNI we have to determine how pod traffic will be routed. 

If all pod ip addresses are routed to nodes as is the case with kubenet we can use cilium to replace the function of managing the pod networks on each node and rely on the existing routing mechanisms outside the node to ensure that traffic continues to flow. 

In scenarios where we are changing the type of encapsulation or perhaps migrating from a direct routed scenario to an overlay configuration this can be a bit more difficult. Specifically we need some mechanism to ensure that routes exist on all nodes to be sure that traffic can be routed during the period of migration. 

Finally, if network policy is implemented in the existing CNI configuration all policies should be removed prior to migration. The reason for this is that Cilium implements network policy in eBPF by attaching a program to each pod. In most cases the existing CNI implementation may not support network policy and if it does it will have implemented this in iptables on the underlying host. This can lead to conditions where traffic is enforced by policy both in ebpf and in iptables and can be difficult to troubleshoot. 


### In Place

This stategy assumes that we will reuse the existing kubelets and deploy cilium in addition to the existing CNI. We will disable the unmanagedpods watcher functionality of the cilium operator until all pods have been migrated. 

We will deploy Cilium in addition to the existing CNI. 

Cilium will take precedence for any new pods deployed. 

All pods should be restarted in whatever order makes sense until all pods are on the managed by cilium. 

Once complete you can uninstall the old cni components and restart the nodes. 

restarting the nodes themselves is a step to ensure that any configurations of iptables or any underlying network interfaces or configurations are removed after the migration is complete. 



### Node Pool

This strategy assumes that we will deploy new nodes with cilium installed and migrate existing pods to the new nodes. Once completely migrated we will remove the old nodes with the existing CNI from service. 

## Validating connectivity and testing.

To validate connectivity throughout the process we can use a tool like [goldpinger](https://github.com/bloomberg/goldpinger) this tool will startup a pod and discover all the other pods that have been deployed as and establish connectivity and report on success or failure. You can read more in the documentation for the tool. This tool does a decent job of validating connectivity across CNI implementations. 

## Amazon

### Overview
This migration assumes that we will migrate applications from existing node groups to new ones with cilium installed per our discussion. 

Initially we are going to bring up an eks cluster with two node groups. 

`ng-1` which is meant to be the existing node group that will be serviced by aws-cni. 

`ng-cilium` a new node group that will be serviced by cilium. 

initially ng-cilium will be defined but have no nodes deployed. We will also taint this node group according to ensure that the nodes do not become ready until the cilium has been successfully deployed to it. 

We will modify the aws-node daemonset to ensure that it doesn't deploy to anything but the nodes that we know about before introducing the new ng-cilium nodes. 

We will then deploy cilium on the new ng-cilium node group and validate that cilium is installed and connectivity works across nodes and across cni implmentations. 

We then validate that we can migrate workloads.

Since cilium is configured and deployed in ENI mode. Traffic between the CNI implementations is supported.

Note that netowrk policy should not be applied until all applications have migrated to cilium. 



#### eks-config
The sample eks config file:
filename: `eks-config.yaml`
:::spoiler eks-config.yaml
```
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig
metadata:
  name: dc-migrate
  region: us-west-1

managedNodeGroups:
- name: ng-1
  instanceType: m5.large
  desiredCapacity: 3
  volumeSize: 80
  privateNetworking: true
- name: ng-cilium
  instanceType: m5.large
  desiredCapacity: 0
  volumeSize: 80
  privateNetworking: true
  # taint nodes so that application pods are
  # not scheduled/executed until Cilium is deployed.
  # Alternatively, see the note above regarding taint effects.
  taints:
   - key: "node.cilium.io/agent-not-ready"
     value: "true"
     effect: "NoSchedule"
```
:::

#### create cluster
eksctl create cluster --config-file 

#### Deploy Sample Application
Once the cluster is deployed. We are going to deploy our sample application to it. 
This consists of a few files that were templated from helm.

:::spoiler goldpinger-infra.yaml
```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: goldpinger
  labels:
    app.kubernetes.io/name: goldpinger
    app.kubernetes.io/instance: goldpinger
---
# Source: goldpinger/templates/clusterrole.yaml
kind: Role
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: goldpinger-role
  labels:
    app.kubernetes.io/name: goldpinger
    app.kubernetes.io/instance: goldpinger
rules:
- apiGroups: [""]
  resources: ["pods"]
  verbs: ["list"]
---
kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: goldpinger-rolebinding
  labels:
    app.kubernetes.io/name: goldpinger
    app.kubernetes.io/instance: goldpinger
subjects:
  - kind: ServiceAccount
    name: goldpinger
roleRef:
  kind: Role
  name: goldpinger-role
  apiGroup: rbac.authorization.k8s.io
---
apiVersion: v1
kind: Service
metadata:
  name: goldpinger
  labels:
    app.kubernetes.io/name: goldpinger
    app.kubernetes.io/instance: goldpinger
spec:
  type: LoadBalancer
  ports:
    - port: 8080
      targetPort: 8080
      protocol: TCP
      name: http
  selector:
    app.kubernetes.io/name: goldpinger
    app.kubernetes.io/instance: goldpinger
```
:::

:::spoiler goldpinger-ds.yaml
```yaml=
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: goldpinger
  labels:
    app.kubernetes.io/name: goldpinger
    app.kubernetes.io/instance: goldpinger
spec:
  selector:
    matchLabels:
      app.kubernetes.io/name: goldpinger
      app.kubernetes.io/instance: goldpinger
  template:
    metadata:
      labels:
        app.kubernetes.io/name: goldpinger
        app.kubernetes.io/instance: goldpinger
    spec:
      priorityClassName: 
      serviceAccountName: goldpinger
      containers:
      - name: goldpinger-daemon
        image: mauilion/goldpinger:v3.5.1
        imagePullPolicy: IfNotPresent
        env:
          - name: HOSTNAME
            valueFrom:
              fieldRef:
                fieldPath: spec.nodeName
          - name: HOST
            value: "0.0.0.0"
          - name: PORT
            value: "8080"
          - name: LABEL_SELECTOR
            value: "app.kubernetes.io/name=goldpinger"
          - name: REFRESH_INTERVAL
            value: "5"
        ports:
        - name: http
          containerPort: 8080
          protocol: TCP
        livenessProbe:
          httpGet:
            path: /
            port: http
        readinessProbe:
          httpGet:
            path: /
            port: http
        resources:
            {}

```
:::
:::spoiler goldpinger-deploy.yaml
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  generation: 1
  labels:
    app.kubernetes.io/instance: goldpinger
    app.kubernetes.io/name: goldpinger
  name: goldpinger
spec:
  replicas: 5
  selector:
    matchLabels:
      app.kubernetes.io/instance: goldpinger
      app.kubernetes.io/name: goldpinger
  template:
    metadata:
      creationTimestamp: null
      labels:
        app.kubernetes.io/instance: goldpinger
        app.kubernetes.io/name: goldpinger
    spec:
      containers:
      - env:
        - name: HOSTNAME
          valueFrom:
            fieldRef:
              apiVersion: v1
              fieldPath: spec.nodeName
        - name: HOST
          value: 0.0.0.0
        - name: PORT
          value: "8080"
        - name: LABEL_SELECTOR
          value: app.kubernetes.io/name=goldpinger
        - name: REFRESH_INTERVAL
          value: "5"
        image: mauilion/goldpinger:v3.5.1
        imagePullPolicy: IfNotPresent
        livenessProbe:
          failureThreshold: 3
          httpGet:
            path: /
            port: http
            scheme: HTTP
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 1
        name: goldpinger-daemon
        ports:
        - containerPort: 8080
          name: http
          protocol: TCP
        readinessProbe:
          failureThreshold: 3
          httpGet:
            path: /
            port: http
            scheme: HTTP
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 1
        resources: {}
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
      dnsPolicy: ClusterFirst
      restartPolicy: Always
      schedulerName: default-scheduler
      securityContext: {}
      serviceAccount: goldpinger
      serviceAccountName: goldpinger
      terminationGracePeriodSeconds: 30
```
:::
#### validate connectivity
Once these are deployed we can see a full mesh of connectivity between all pods by taking a look at the deployed loadBalancer service. 

https://yourawslb.elb.amazon.com:8080


We haven't yet deployed cilium. 

#### reconfigure aws-node daemonset 

The intention here is to leave the existing node-group alone and deploy cilium only on the new ng-cilium node-group. Normally we would remove the `aws-node` daemonset in the kube-system namespace and that would ensure that the aws-node agent is not run on any new nodes introduced. Instead we are going to reconfigure the daemonset to only deploy to nodes that we already know about. The first step is to label all existing nodes with cni=aws-cni



label all existing nodes with the label: 
`kubectl label nodes --all cni=aws-cni`


Next we will modify the aws-node daemonset to only target nodes with the cni=aws-cni node label.

:::spoiler before:
```
    spec:
      affinity:
        nodeAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            nodeSelectorTerms:
            - matchExpressions:
              - key: kubernetes.io/os
                operator: In
                values:
                - linux
              - key: kubernetes.io/arch
                operator: In
                values:
                - amd64
                - arm64
              - key: eks.amazonaws.com/compute-type
                operator: NotIn
                values:
                - fargate

```
:::
:::spoiler after:
```
    spec:
      affinity:
        nodeAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            nodeSelectorTerms:
            - matchExpressions:
              - key: kubernetes.io/os
                operator: In
                values:
                - linux
              - key: kubernetes.io/arch
                operator: In
                values:
                - amd64
                - arm64
              - key: eks.amazonaws.com/compute-type
                operator: NotIn
                values:
                - fargate
              - key: cni
                operator: In
                values: 
                - aws-cni
```
:::

With this change all of the aws-node pods will restart. 

Once complete we can make sure that all things in the existing node group continue to work. Let's restart the goldpinger deployment and make sure all pods come back up and in service. 

#### validate connectivity 
`kubectl rollout restart deploy/goldpinger`

Once satisfied that things are working well let's move on to bringing up the ng-cilium node group and adding cilium to the cluster. 

```
eksctl get nodegroup --cluster dc-migrate
CLUSTER		NODEGROUP	STATUS	CREATED			MIN SIZE	MAX SIZE	DESIRED CAPACITY	INSTANCE TYPE	IMAGE ID	ASG NAME					TYPE
dc-migrate	ng-1		ACTIVE	2023-08-29T04:03:53Z	3		3		3			m5.large	AL2_x86_64	eks-ng-1-7ac51fbb-24bb-7be9-4938-10959e334240		managed
dc-migrate	ng-cilium	ACTIVE	2023-08-29T04:03:57Z	0		1		0			m5.large	AL2_x86_64	eks-ng-cilium-d0c51fbb-2be4-13e1-218c-a60067c85215	managed
```
#### scale up ng-cilium

To scale up the `ng-cilium` node group use eksctl

`eksctl scale nodegroup --cluster dc-migrate ng-cilium --nodes 3 --nodes-max=3`

These nodes will come up in an NotReady state. You should also verify that the `aws-node` daemonset is not deployed to these nodes based on our affinity rules. 

#### install cilium

Let's install cilium. 

In the values below I have configured cilium to deploy in 'eni' mode. This means that cilium will handle the eni's allocated to a node and configure pods accordingly. 

We have also disabled the unmanagedPodWatcher to ensure that pods aren't restarted until we are ready for them to be.

Finally, We have also configured affinity for cilium to install only on nodes that don't match the `cni=aws-cni` label. Ideally we want only nodes that we did not know about before scaling up the `ng-cilium` node group to have cilium installed on them. 


cilium install --dry-run-helm-values


:::spoiler cilium-values.yaml
```
cluster:
  name: dc-migrate-us-west-1-eksctl-io
egressMasqueradeInterfaces: eth0
eni:
  enabled: true
bpf:
  hostLegacyRouting: true
ipam:
  mode: eni
operator:
  replicas: 1
  unmanagedPodWatcher:
    restart: false
serviceAccounts:
  cilium:
    name: cilium
  operator:
    name: cilium-operator
tunnel: disabled
hubble:
  enabled: true
  relay:
    enabled: true
affinity:
  nodeAffinity:
    requiredDuringSchedulingIgnoredDuringExecution:
      nodeSelectorTerms:
      - matchExpressions:
        - key: cni
          operator: NotIn
          values:
          - aws-cni
```
:::




#### Verify that Cilium is up. 
Once cilium comes up we will see the new nodes in a `Ready` state. 

`kubectl get nodes -l cni!=aws-cni`

We will also see that the goldpinger daemonset has deployed a pod on each of our new nodes. 

```
for node in $(kubectl get nodes -l cni!=aws-cni -ojsonpath='{.items[*].metadata.name}'); do kubectl get pods -owide --field-selector spec.nodeName=$node --no-headers; done
```
#### Verify that all goldpinger pods are reachable. 

https://yourawslb.elb.amazon.com:8080 

Cordon existing nodes and restart the goldpinger daemonset. 

kubectl cordon -l cni=aws-cni

restart the goldpinger deployment. 

kubectl rollout restart deploy goldpinger. 

We should see all the goldpinger deployment pods come up on the new nodes. 

Validate with any other tests you are concerned with. 

#### Migrate all workloads over to the cilium nodes. 
Finish migrating all pods over to the new node group. 

#### delete the aws-node daemonset. 
Once all migration is completed we can remove the aws-node daemonset. Note that doing this before the migration is complete will break the cni configuration on nodes where cilium isn't installed. This will keep new pods from being created on those nodes. Any existing pods that are still up will continue to operate. 


#### scale down and remove the old node groups. 
scale down the none cilium enabled node groups. 


