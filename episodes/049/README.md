# Episode 49: Graceful Termination

{%youtube 9GBxJMp6UkI%}

# Headlines 

- [Echo News!](https://isovalent-9197153.hs-sites.com/echo-news-episode-8-ebpf-for-traffic-shaping-database-management-and-service-meshes.-egress-gateways-dual-stack-clusters-and-multus-with-cilium)

# Graceful termination



[Tracking terminating endpoints]
* [docs](https://kubernetes.io/docs/concepts/services-networking/endpoint-slices/#terminating)
* [kep](https://github.com/kubernetes/enhancements/blob/master/keps/sig-network/1672-tracking-terminating-endpoints)

[Node Graceful Shutdown]
* [docs](https://kubernetes.io/docs/concepts/architecture/nodes/)
* [kep](https://github.com/kubernetes/enhancements/blob/0e4d5df19d396511fe41ed0860b0ab9b96f46a2d/keps/sig-node/2000-graceful-node-shutdown)

[Cilium]
* [Blog](https://isovalent.com/blog/post/2021-12-release-111#graceful-termination)
* [Docs](https://docs.cilium.io/en/latest/gettingstarted/kubeproxy-free/#graceful-termination)
* [tests](https://github.com/cilium/cilium/pull/20163/files#diff-f3a7f606d5222450e169c8c12caa36f9b2bfbd962a762830a4799150fb0685a4)
* [demo app](https://github.com/cilium/graceful-termination-test-apps)