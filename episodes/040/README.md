# Episode 40: Cilium host firewall 

## Headlines


* [BTFGen](https://kinvolk.io/blog/2022/03/btfgen-one-step-closer-to-truly-portable-ebpf-programs/)
* [Merbridge](https://istio.io/latest/blog/2022/merbridge/) shortening the data path in Istio (like [Cilium has been doing since 2018](https://cilium.io/blog/2018/08/07/istio-10-cilium) 😉)
* [eBPF — kProbes failure on AWS Amazon Linux 2 image](https://medium.com/@Amit_Sides/ebpf-kprobes-failure-on-aws-amazon-linux-2-image-522914639552)

## Cilium host firewall 
* This episode inspired by [this blog post](https://medium.com/@charled.breteche/kubernetes-security-explore-cilium-host-firewall-and-host-policies-de93ea9da38c) by Charles-Eduoard Brétéché
* Cilium host firewall [docs](https://docs.cilium.io/en/stable/gettingstarted/host-firewall/)

```
cilium install 

ks edit cm cilium-config
#   allow-localhost: "policy"
#   enable-host-firewall: "true"
#   enable-policy: always
#   policy-audit-mode: "true"

# restart with new config
ks delete pods -l k8s-app=cilium

# Or just enable hubble, this restarts Cilium anyway
cilium hubble enable

k get nodes
# pick control plane node and set NODE_NAME=<name>

ks get pods -l "k8s-app=cilium" -o wide
# CILIUM_CP=<pod name on CP node>
# CILIUM_WORKER=<pod name on worker node>

# Listing endpoints
ks exec $CILIUM_POD -- cilium endpoint list

ks exec $CILIUM_POD -- cilium endpoint get <ID>

# Monitoring policy verdicts
ks exec $CILIUM_POD -- cilium monitor -t policy-verdict

# Observe flows that got caught by audit
cilium hubble port-forward &
hubble observe --verdict AUDIT

# curl.yaml pod in host network namespace, sleeps so we can exec into it
k exec -it curl -- sh 

# Check node that pod's running on and look at endpoint list

nslookup example.com
ks exec $CILIUM_WORKER -- cilium monitor -t policy-verdict | grep <IP address>

# Label the node, look at it again in endpoint list 
k label node kind-worker policy=example

# Apply ccnp-example.yaml policy to allow port 80 egress traffic from labelled node
```