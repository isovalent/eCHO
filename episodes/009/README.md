# Episode 9: XDP and load balancing with Daniel Borkmann

[YouTube](https://youtu.be/OIyPm6K4ooY)
11 June 2021

with [Daniel Borkmann](https://github.com/borkmann), hosted by [Liz Rice](https://twitter.com/lizrice)

## Headlines

- eBPF Summit 2021 on August 18-19: [pre-register now!](https://ebpf.io/summit-2021)
- Quentin Monnet's summary of [Implementing eBPF on Windows](https://lwn.net/SubscriberLink/857215/1df543fa904b3f17/) 
  - Don't miss [next week's eCHO with Dave Thaler](https://youtu.be/LrrV-eo6fug)!
- [How Netflix uses eBPF flow logs at scale for network insight](https://netflixtechblog.com/how-netflix-uses-ebpf-flow-logs-at-scale-for-network-insight-e3ea997dca96)
- [hBPF - eBPF on Hardware](https://github.com/rprinz08/hBPF)

## XDP and load balancing with Daniel Borkmann

### XDP layer / general

- Background & history: [Daniel's KubeCon EU 2020 talk](https://static.sched.com/hosted_files/kccnceu20/8f/Aug19_eBPF_and_Kubernetes_Little_Helper_Minions_for_Scaling_Microservices_Daniel_Borkmann.pdf)
- [XDP intro in reference guide](https://docs.cilium.io/en/stable/bpf/#xdp)
- [ACM paper on XDP](https://dl.acm.org/doi/10.1145/3281411.3281443) ([slides](https://blog.tohojo.dk/slides/conext18-xdp.pdf))
- XDP as stack hardening use case
  - Slide 16-32 in [Daniel's presentation on BPF as a Fundamentally Better Dataplane](https://ebpf.io/summit-2020-slides/eBPF_Summit_2020-Keynote-Daniel_Borkmann-BPF_as_a_fundamentally_better_dataplane.pdf) ([video](https://www.youtube.com/watch?v=Qhm1Zn_BNi4))
  - Discussions around hardening from Google developers:
    - [More on network stack resilience](http://vger.kernel.org/netconf2019_files/netconf%202019%20-%20willem%20de%20bruijn%20-%20resilient%20rx%20+%20scaling%20udp.pdf)
    - [TCP SYN cookies from XDP](https://netdevconf.info/0x14/session.html?talk-issuing-SYN-cookies-in-XDP)
- Random Bits of History:
  - 2016: [XDP first merged into upstream kernel](https://lore.kernel.org/lkml/20160727.010753.2221383279830501569.davem@davemloft.net/)
  - 2017: [XDP MythBusters, netdevconf 2.1 keynote](https://netdevconf.info/2.1/slides/apr7/miller-XDP-MythBusters.pdf) ([video](https://www.youtube.com/watch?v=8Cxg7mpVIWw))
  - 2018: 
    - [1.5 years of XDP in production at Facebook](http://vger.kernel.org/lpc-networking2018.html#session-10) ([video](https://www.youtube.com/watch?v=E1QKn_AjuJk))
    - [AF_XDP, "the path to DPDK speeds"](http://vger.kernel.org/lpc-networking2018.html#session-11) ([video](https://www.youtube.com/watch?v=JmGfJok32Kw))
    - [bpfilter (aka iptables) at XDP](https://cilium.io/blog/2018/04/17/why-is-the-kernel-community-replacing-iptables)
      - Latter seeing interest in 2021 [again](https://lore.kernel.org/bpf/20210603101425.560384-1-me@ubique.spb.ru/)
  - Ongoing XDP developments:
    - [XDP bonding support](https://lore.kernel.org/bpf/20210609135537.1460244-1-joamaki@gmail.com/)
    - [XDP jumbo frame support](https://lore.kernel.org/bpf/cover.1617885385.git.lorenzo@kernel.org/)

### Cilium related

- XDP in releases:
  -  1.8: [kube-proxy at XDP layer](https://cilium.io/blog/2020/06/22/cilium-18#kube-proxy-replacement-at-the-xdp-layer)
  -  1.9: [Maglev consistent hashing for XDP (& tc)](https://cilium.io/blog/2020/11/10/cilium-19#maglev)
  -  1.10: [Standalone L4LB at XDP layer](https://cilium.io/blog/2021/05/20/cilium-110#xdp-based-standalone-load-balancer)
- Cilium XDP-related documentation:
  - [Getting Started Guide for kube-proxy replacement & XDP](https://docs.cilium.io/en/v1.10/gettingstarted/kubeproxy-free/#loadbalancer-nodeport-xdp-acceleration), including list of supported drivers
- Talk on Cilium's load balancing & introduction to K8s service handling:
  - [Linux Plumbers, 2020](https://linuxplumbersconf.org/event/7/contributions/674/attachments/568/1002/plumbers_2020_cilium_load_balancer.pdf) ([video](https://www.youtube.com/watch?v=UkvxPyIJAko&t=21s))
  - [KubeCon US, 2019](https://docs.google.com/presentation/d/1cZJ-pcwB9WG88wzhDm2jxQY4Sh8adYg0-N3qWQ8593I/edit#slide=id.g7055f48ba8_0_0) ([video](https://www.youtube.com/watch?v=bIRwSIwNHC0)
- Talk on Cilium and Maglev consistent hashing:
  - [FOSDEM, 2021](https://fosdem.org/2021/schedule/event/containers_ebpf_kernel/attachments/slides/4358/export/events/attachments/containers_ebpf_kernel/slides/4358/Advanced_BPF_Kernel_Features_for_the_Container_Age_FOSDEM.pdf) ([video](https://www.youtube.com/watch?v=PJY-rN1EsVw))

### XDP performance related

 - [Facebook's migration in 2017 from IPVS to XDP, performance graph in slide 4](https://netdevconf.info/2.1/slides/apr6/zhou-netdev-xdp-2017.pdf) ([video](https://www.youtube.com/watch?v=YEU2ClcGqts))
 - [ACM paper on XDP with performance comparison to DPDK](https://dl.acm.org/doi/10.1145/3281411.3281443)
 - [Cilium release blog with kube-proxy performance comparison](https://cilium.io/blog/2020/06/22/cilium-18#kube-proxy-replacement-at-the-xdp-layer)
 - [Verizon on performance analysis of XDP programs](https://www.usenix.org/system/files/lisa21_slides_jones.pdf) ([video](https://www.youtube.com/watch?v=TjAbnY5G4io))

### Other XDP projects outside of Cilium

- Platforms and runtimes related outside of Linux kernel:
  - [DPDK AF_XDP driver](https://doc.dpdk.org/guides/nics/af_xdp.html)
  - [XDP also on Windows](https://cloudblogs.microsoft.com/opensource/2021/05/10/making-ebpf-work-on-windows/)
- Other XDP L4LBs aside from Cilium used in production:
  - Facebook's [Katran project](https://github.com/facebookincubator/katran)
    - [Blog post](https://engineering.fb.com/2018/05/22/open-source/open-sourcing-katran-a-scalable-network-load-balancer/)
    - Also used by [Dropbox](https://dropbox.tech/infrastructure/boosting-dropbox-upload-speed) in production
  - Cloudflare's Unimog project
    - [Blog post](https://blog.cloudflare.com/unimog-cloudflares-edge-load-balancer/)
    - [Cloudflare architecture and how BPF eats the world](https://blog.cloudflare.com/cloudflare-architecture-and-how-bpf-eats-the-world/)
  - Github Load Balancer (GLB)
    - [Blog post](https://github.blog/2018-08-08-glb-director-open-source-load-balancer/)
    - [GLB also supporting XDP](https://github.com/github/glb-director)
    
### XDP in academia

- There are 85 citations on our ACM paper on XDP by now. More recent work that is XDP/BPF related:
  - Accelerating Memcached in XDP
    - [Paper](https://www.usenix.org/system/files/nsdi21-ghigoff.pdf)
    - [Blog post](https://pchaigno.github.io/ebpf/2021/04/12/bmc-accelerating-memcached-using-bpf-and-xdp.html)
  - hXDP - XDP BPF programs in FPGA-powered NICs
    - [Paper](https://arxiv.org/pdf/2010.14145.pdf)
    - [Blog post](https://pchaigno.github.io/ebpf/2020/11/04/hxdp-efficient-software-packet-processing-on-fpga-nics.html)
  - hBPF - eBPF CPU in hardware
    - [hBPF project](https://github.com/rprinz08/hBPF)
    - [Blog post](https://www.min.at/prinz/?x=entry:entry210403-164137)

### XDP tutorials/guides for writing own programs

 - [XDP Tutorial](https://github.com/xdp-project/xdp-tutorial/)

