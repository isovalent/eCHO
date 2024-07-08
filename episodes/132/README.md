# Episode 132 - Pushing the limits of eBPF 🐝

![](https://www.youtube.com/watch?v=BrH087OK8hU)

Liz Rice & Dan Finneran

## Events

* [OSS Summit 2024 NA](https://events.linuxfoundation.org/open-source-summit-north-america/program/schedule/)
* [Devoxx 2024 France](https://www.devoxx.fr/schedule/talk/?id=5292) - Au cœur de la ruche eBPF!
* Upcoming CFPs
    * KubeCon NA 2024
    * Open Source Summit EU
    * Various KCDs (hopefully in your area)
* AWS Summit London next week


## Cloud Native news!

* [Kubernetes v1.30: Uwubernetes](https://kubernetes.io/blog/2024/04/17/kubernetes-v1-30-release/)
    * Kubernetes is LoadBalancer aware: `.status.loadBalancer.ingress.ipMode`
        * `VIP` - traffic comes FROM the loadbalancer
        * `Proxy` - traffic is sent to the pod directly from where it came
    * Traffic distribution for Kubernetes services:
        * `PreferClose` - means that treaffic will aim to be topologically close to the client, same rack etc..
    * Dual stack has graduated to **stable**
* [Elastic Open telemetry agent is opensource and using eBPF](https://www.elastic.co/blog/elastic-universal-profiling-agent-open-source)
* [Cisco announce **HyperShield** 🛡️](https://blogs.cisco.com/news/cisco-hypershield-security-reimagined-hyper-distributed-security-for-the-ai-scale-data-center) using eBPF and Cilium

## Pushing the limits of eBPF 🐝

* [Conways game of life](https://en.wikipedia.org/wiki/Conway%27s_Game_of_Life)
* Talk [slides](https://speakerdeck.com/lizrice/ebpfs-abilities-and-limitations-the-truth)
* eBPF [code](https://gist.github.com/lizrice/0e098e08fa09bc7e20af5cfee0777dfd)

{%youtube tClsqnZMN6I %}
