# Episode 115: KubeCon/CiliumCon NA Review

## Recent news

* Tetragon 1.0 release
    * [Thomas Graf's blog post](https://isovalent.com/blog/post/tetragon-release-10/)
    * [ITOps Times open source project of the week](https://www.itopstimes.com/os/itops-open-source-project-of-the-week-tetragon/)
    * [SDx Central](https://www.sdxcentral.com/articles/analysis/tetragon-adds-visibility-to-kubernetes-with-open-source-runtime-security-platform/2023/11/)
* [Cilium Graduation](https://www.cncf.io/announcements/2023/10/11/cloud-native-computing-foundation-announces-cilium-graduation/)
* [Children's eBPF book](https://isovalent.com/books/children-guide-to-ebpf/)!

## Recent Events (today's discussion)

* [Rejekts](https://youtube.com/playlist?list=PLnfCaIV4aZe-4zfJeSl1bN9xKBhlIEGSt&si=lucY1AIPLcwcn840) :us:
* [CiliumCon](https://youtube.com/playlist?list=PLj6h78yzYM2NeU8D0929lbiy6mmd5qeaX&si=u3KCZE-2eEA3nJYN) :us:
* KubeCon :us:
* [eBPF Documentary](https://youtu.be/Wb_vD3XZYOA?si=YPR8HPRKKhbliBnU) :film_projector:
  * [Behind the scenes material](https://isovalent.com/blog/post/ebpf-documentary-creation-story/?utm_source=hs_email&utm_medium=email&_hsenc=p2ANqtz-_6XcieQYba5M54Y530UDJlO2CmRyB5A078r_8Mr56n-GvH7POOdL9cWdSiA14ArB6awSWN)
* ContainerConf / Continuous LifeCycle in Manheim :de:

## Upcoming Events

* AWS Re:Invent
    * [PeerTalk platform](https://reinvent.awsevents.com/learn/peertalk/experts/) for booking in-person meetings (Liz is under AWS Community Leaders tab)
* CFP open for [KubeCon Paris](https://events.linuxfoundation.org/kubecon-cloudnativecon-europe/program/cfp/) and co-located [Cilium & eBPF Day](https://events.linuxfoundation.org/kubecon-cloudnativecon-europe/co-located-events/cilium-ebpf-day/)

---

### Rejekts

### CiliumCon highlights !

### Co-located day highlights !

- Feels like it's getting to be its own (un)-conference
    - Backstage
    - Telco
    - Istio Day
    - Edge
    - Wasm
    - Argo
    - App Developer
    - Observability Day
    - Multi-TenancyCon
    - AI + HPC

(all at the same time as the maintainer summits are taking place)!

### KubeCon

(We're already pretty exhausted from the Co-located events at this point...)

#### Day one

Keynote focussed largely on AI, with Priyanka running a demo of an LLM (Large Language Model) in Docker for Desktop. Some good panels on important themes such as sustainability, we paid tribute to some good friends and colleagues that are no longer with us. Then we wrapped things up with the graduation ceremony!



#### Day two

Keynote focussed on building good communities and incubation project updates!

[Service Mesh battle scars](https://kccncna2023.sched.com/event/1R2ts) panel took place.

[eBPF: Unlocking the Kernel !](https://t.co/4gdxEfby4G)
![](https://pbs.twimg.com/media/F-W_GcZWkAEJhzB?format=jpg&name=large)

The hive mind mingle took place at the end of KubeCon, including the infamous Cilium rap from Bart Farrell.

#### Day three

Keynote!

"Kubernetes should remain unfinished!", although it does beg the question as to who thinks Kubernetes is anywhere near finished.

Lots of work in the supply chain space, and we're seeing a big push on improvements to the developer experience through things like dapr.

Meta devices become NetKit devices -> [slides](https://static.sched.com/hosted_files/kccncna2023/b9/Borkmann-Turning-Up-Performance-To-11.pdf)

The gist of it (given my basic understanding), is a re-implementation of the device pair that would exist either side of a Pod. Today you would have a virtual ethernet pair (veth), which you can kind of think of as having two sides connected with a network cable:

- Left side, being a network switch connecting to other devices
- Right side, the host that is wanting to connect to other hosts.

It is this veth pair that emulates device connectivity, or in Kubernetes "pod" connectivity to other pods/services or the outside world.

The "problem" that is being addressed is that the `veth` code is in-efficent (not entirely true, but it can be improved upon). The virtual ethernet device simply implements things that we don't require, it also crosses multiple contexts (potentially causing CPU overhead (minimal)). The patchset from Daniel Borkman implements a minimal device that sits inside the pod, that resides in a single context and implements JUST the functionality required to enable communication.

Initial testing of the Netkit device put's it on par with baseline host throughput, which is the best case scenario.

Will land in the Linux Kernel 6.7, so might take a little while to get in peoples hands (but exciting non-the-less)

## Takeaways

- Has everyone recovered?
- If so are we ready to go again, Cfps are due soon `:lolsob:`
- For those writing CFPs [this](https://twitter.com/dims/status/1694705434311933959/photo/1) from dims, is an excellent set of criteria for helping write a good CFP.
