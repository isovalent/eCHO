# Episode 63: cilium bpf recorder

{%youtube zh1y155aeJM %} 

## Headlines

- [echo news 15](https://isogo.to/echo-news-15)
- [eBPF Summit 2022](https://ebpf.io/summit-2022) is still open for registration

## Notes

- [Cilium BPF Recorder](https://docs.cilium.io/en/v1.12/cmdref/cilium_bpf_recorder/)
- [XDP](https://xdp-project.net)
- [BPF and XDP Reference Guide](https://docs.cilium.io/en/v1.12/bpf/#bpf-guide)


helm install 
```
helm upgrade cilium cilium/cilium --version 1.12.2     --namespace kube-system     --set tunnel=disabled     --set autoDirectNodeRoutes=true     --set kubeProxyReplacement=strict     --set loadBalancer.acceleration=native     --set loadBalancer.mode=hybrid     --set k8sServiceHost=api.c1.k8s.work     --set k8sServicePort=6443 --set ipv4NativeRoutingCIDR=10.0.0.0/8 
```


```
1       reserved:host
2       reserved:world
3       reserved:unmanaged
4       reserved:health
5       reserved:init
6       reserved:remote-node
7       reserved:kube-apiserver
        reserved:remote-node
8       reserved:ingress
```
