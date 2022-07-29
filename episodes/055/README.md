# Episode 55: Cilium OSS 1.12 Release

{%youtube oeh3u4V2--M%}



# Hosts:
Bruno and Quentin

# Guests:

- Name: Vincent Li F5
    - Feature:  [Add VXLAN Tunnel Endpoint (VTEP) Integration](https://github.com/cilium/cilium/pull/17370)
- Name: Hemanth Malla Datadog
    - Feature: [Adding support for AWS ENI prefix delegation - IPv4 Only]()
    - 
- Name: Will Daly Microsoft
    - Feature: [Delegated IPAM plugin]()

## News:
* [Echo News 11](isogo.to/echo-news-11)
    * Cilium 1.12 and Service Mesh out now. eBPF for storage, IoT, and continuous profiling


## List of Features to discuss
![](https://i.imgur.com/tqOvKDC.png)

### [Release Notes](https://github.com/cilium/cilium/releases/tag/v1.12.0)

### [Cilium Blog](https://isovalent.com/blog/post/cilium-release-112/)

## Cilium OSS 1.12 – New Features at a Glance

### External Workload Improvements:

* [**Vincent Li**] VTEP support: This new integration allows Cilium to peer with VXLAN Tunnel Endpoint devices in on-prem environments. [More details](https://isovalent.com/blog/post/cilium-release-112/#vtep-support)
* [**Bruno**] Egress Gateway promoted to Stable: Cilium enables users to route selected cluster-external connections through specific Gateway nodes, masquerading them with predictable IP addresses to allow integration with traditional firewalls that require static IP addresses. A new `CiliumEgressGatewayPolicy` CRD improves the selection of the Gateway node and its egress interface. Additional routing logic ensures compatibility with ENI interfaces. [More details](https://isovalent.com/blog/post/cilium-release-112/#egress-stable)
* [**Quentin**] Improved BGP control plane: IPv6 support has been added to the BGP control plane. By leveraging a new feature-rich BGP engine, Cilium can now set up IPv6 peering sessions and advertise BGP IPv6 Pod CIDRs. [More details](https://isovalent.com/blog/post/cilium-release-112/#improved-bgp)


### Service Mesh & Ingress:

* [**Bruno**] Integrated Ingress Controller: A fully compliant Kubernetes Ingress controller embedded into Cilium. Additional annotations are supported for more advanced use cases. [More details](https://isovalent.com/blog/post/cilium-release-112/#ingress)
* Sidecar-free datapath option: A new datapath option for Cilium Service Mesh as an alternative to the Istio integration allowing to run a service mesh without sidecars. [More details](https://isovalent.com/blog/post/cilium-release-112/#sidecar-free)
* Envoy CRD: A new Kubernetes CRD making the full power of Envoy available whereever Cilium runs. You can express your requirements in an Envoy Configuration and inject it anywhere in the network. [More details](https://isovalent.com/blog/post/cilium-release-112/#envoy-crd)
* Gateway API: The work on supporting Gateway API has started. If you are interested in contributing, reach out on [Slack](https://cilium.io/slack).

### Cluster Mesh:
[**Bruno**]
* Topology-aware routing and service affinity: With a single line of YAML annotation, services can be configured to prefer endpoints in the local or remote cluster, if endpoints in multiple clusters are available. [More details](https://isovalent.com/blog/post/cilium-release-112/#service-affinity)
* Simplified cluster connections with Helm: New simplified user experience to connect Kubernetes clusters together using Cilium and Helm. [More details](https://isovalent.com/blog/post/cilium-release-112/#clustermesh-helm)
* Lightweight Multi-Cluster for External Workloads: Cluster Mesh now supports special lightweight remote clusters. This allows running lightweight clusters for use by external workloads. [More details](https://isovalent.com/blog/post/cilium-release-112/#remote-kubernetes-api)

### Security:
[**Duffie**]
* Running Cilium as unprivileged Pod: You can now run Cilium as an unprivileged container/Pod to reduce the attack surface of a Cilium installation. [More details](https://isovalent.com/blog/post/cilium-release-112/#unprivileged-pod)
* Reduction in required Kubernetes privileges: The required Kubernetes privileges have been greatly reduced to the least needed for Cilium to operate. [More details](https://isovalent.com/blog/post/cilium-release-112/#kubernetes-privileges)
* Network policies for ICMP: Cilium users can now allow a subset of ICMP traffic on egress and ingress, with the usual CiliumNetworkPolicies and CiliumClusterwideNetworkPolicies. [More details](https://isovalent.com/blog/post/cilium-release-112/#icmp-cnp)

### Load-Balancing:

* [**Daniel**] NAT46/64 Support for Load Balancer: Cilium L4 load-balancer (L4LB) now supports NAT46 and NAT64 for services. This allows exposing an IPv6-only Pod via an IPv4 service IP or vice versa. This is particularly useful to load-balance IPv4 client traffic at the edge to IPv6-only clusters. [More details](https://isovalent.com/blog/post/cilium-release-112/#ingress)
* [**Daniel**] Quarantining Service backends: A new API to quarantine service backends that are unreachable, and unable to do load-balancing. The latter can be utilized by health checker components to drain traffic from unstable backends, or backends that should be placed into maintenance mode. [More details](https://isovalent.com/blog/post/cilium-release-112/#quarantining)
* [**Quentin**] Improved Multi-Homing for Load Balancer: Cilium’s datapath is extended to support multiple native network devices and multiple paths. [More details](https://isovalent.com/blog/post/cilium-release-112/#multihoming-lb)

### Networking:

* [**Daniel**] BBR congestion control for Pods: Cilium is now the first CNI to support TCP BBR congestion control for Pods in order to achieve significantly better throughput and lower latency for Pods exposed to lossy networks such as the Internet. [More details](https://isovalent.com/blog/post/cilium-release-112/#bbr)
* [**Daniel**] Bandwidth Manager promoted to Stable: The bandwidth manager used to rate-limit Pod traffic and optimize network utilization has been promoted to stable [More details](https://isovalent.com/blog/post/cilium-release-112/#bm-stable)
* [**Duffie**]Dynamic Allocation of Pod CIDRs (beta): A new IPAM mode that improves Pod IP address space utilization due to dynamic assignment of Pod CIDRs to nodes. The latter can now allocate additional Pod CIDR pools dynamically to each node based on their usage. [More details](https://isovalent.com/blog/post/cilium-release-112/#dynamic-allocation)
* [**Will Daly**] Delegated IPAM mode: This is another new IPAM mode that enables a feature in the CNI spec that supports delegating the handling of IPAM. [More details](https://www.cni.dev/plugins/current/ipam/host-local/)
* Send ICMP unreachable host on Pod deletion: When a Pod is deleted, Cilium can now install a route which informs other endpoints that the Pod is now unreachable. [More details](https://isovalent.com/blog/post/cilium-release-112/#send-icmp-unreachable)
* [**Hemanth Malla**] AWS ENI prefix delegation: Cilium now supports the AWS ENI prefix delegation feature, which effectively increases the allocatable IP capacity when running in ENI mode. [More details](https://isovalent.com/blog/post/cilium-release-112/#eni-prefix-delegation)
* AWS EC2 instance tag filter: A new Cilium Operator flag improves the scalability in large AWS subscriptions. [More details](https://isovalent.com/blog/post/cilium-release-112/#ec2-tag-filter)

### Tetragon:
[**Duffie**]
* Initial release: Initial release of Tetragon that provides security observability and runtime enforcement using eBPF. [More details](https://isovalent.com/blog/post/cilium-release-112/#tetragon)

### User Experience:
[**Duffie**]
* [**Nicolas**] Automatic Helm Value Discovery: Cilium CLI is now capable of automatically discovering cluster types and generating the ideal helm config values for the discovered cluster type. [More details](https://isovalent.com/blog/post/cilium-release-112/#cilium-cli-helm)
* [**Nicolas**] AKS BYOCNI support: Cilium and the Cilium CLI now support AKS clusters created with the new Bring-Your-Own Container Network Interface (BYOCNI) mode. [More details](https://isovalent.com/blog/post/cilium-release-112/#aks-byocni)
* [**Nicolas**] Improved chaining mode support: Improved integration between Cilium and cloud CNI plugins for Azure and AWS in chaining mode. [More details](https://isovalent.com/blog/post/cilium-release-112/#improved-chaining)
* [**Nicolas**] Better troubleshooting with Hubble CLI: Many improvements to the Hubble CLI including a better indication of whether a particular connection has been allowed or denied. [More details](https://isovalent.com/blog/post/cilium-release-112/#better-hubble-cli)


* [Vincent Li VTEP Support]()
![](https://i.imgur.com/4HCFGMI.png)


* [Hemanth Malla From Datadog adds support aws eni prefix delegation](https://github.com/cilium/cilium/pull/18463)

![](https://user-images.githubusercontent.com/3775612/149224657-0452658e-c0eb-42db-a261-8d773340780a.png)

* [Making operator aware of pending pod backlog on nodes for IP allocations](https://github.com/cilium/cilium/pull/19007)


* [Delegated IPAM plugin] [ref](https://github.com/cilium/cilium/pull/19219) @wedaly)

