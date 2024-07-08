# Episode 140 - Cilium 1.16 and netkit devices

[![Episode 140 - Cilium 1.16 and netkit devices](https://img.youtube.com/vi/hldsOlLCO_Y/0.jpg)](https://www.youtube.com/watch?v=hldsOlLCO_Y "Episode 140 - Cilium 1.16 and netkit devices")

## Headlines

### CFPs open!

* [eBPF Summit](https://ebpf.io/summit-2024/)
  * 11 September, online
  * CFP deadline: 17 July
  * [Program committee application](https://forms.gle/vsDkbaELAiDRP5uU8)
* [Cilium & eBPF day](https://events.linuxfoundation.org/kubecon-cloudnativecon-north-america/co-located-events/cncf-hosted-co-located-events-overview/)
  * 12 November, Salt Lake City (KubeCon co-located event)
  * CFP deadline: 14 July

### Cilium & eBPF

* Digging into how [Tetragon agent reads Process Lifecycle data](https://yuki-nakamura.com/2024/05/23/tetragon-process-lifecycle-observation-tetragon-agent-part/) by Yuki Nakamura
* [Using TC BPF programs to redirect DNS traffic](https://keploy.io/blog/technology/using-tc-bpf-program-to-redirect-dns-traffic-in-docker-containers) by Keploy

## Netkit devices

* [KubeCon talk on netkit](https://sched.co/1R2s5)
* [Cilium netkit documentation](https://docs.cilium.io/en/latest/operations/performance/tuning/#netkit-device-mode)
* [Meta’s evaluation](https://lpc.events/event/17/contributions/1594/) on veth/ipvlan/netkit showing reduction of softirq load compared to veth
* [netkit netlink](https://github.com/vishvananda/netlink/pull/930) Golang support by Bytedance
* [ebpf-go library integration](https://github.com/cilium/ebpf/pull/1257) from Datadog
* Cilium code adding [tcx](https://github.com/cilium/cilium/pull/30103) and [netkit](https://github.com/cilium/cilium/pull/32429) merged for 1.16
* Kernel code for netkit [driver](https://git.kernel.org/pub/scm/linux/kernel/git/torvalds/linux.git/tree/drivers/net/netkit.c), [tcx](https://git.kernel.org/pub/scm/linux/kernel/git/torvalds/linux.git/tree/kernel/bpf/tcx.c), [mprog](https://git.kernel.org/pub/scm/linux/kernel/git/torvalds/linux.git/tree/kernel/bpf/mprog.c)
* [BIG TCP for Cilium](https://isovalent.com/blog/post/big-tcp-on-cilium/)
* [redirect_peer BPF helper](https://cilium.io/blog/2020/11/10/cilium-19/#veth) for Cilium
* [netkit LWN article](https://lwn.net/Articles/949960/)
