# Episode 45: KubeCon EU 2022 Preview
{%youtube 0-mD02qYXHg %}



with [Duffie Cooley](https://twitter.com/mauilion) and [Nicholas Lane](https://twitter.com/apinick)

## Headlines
* [Blog alert! Next-Generation Mutual Authentication with Cilium Service Mesh](https://isovalent.com/blog/post/2022-05-03-servicemesh-security)

* [previous episode on service mesh!](https://youtu.be/XpccICEYqiA)

* **Book download!** [Security Observability with eBPF](https://isovalent.com/blog/post/2022-04-oreilly-security)


### eBPF Day!
[Cloud Native eBPF Day EU](https://events.linuxfoundation.org/cloud-native-ebpf-day-europe/)
[ebpfday schedule](https://cloudnativeebpfdayeu22.sched.com/)

Hosts: [Liz Rice](https://twitter.com/lizrice) & [Sarah Natovny](https://twitter.com/@sarahnovotny)


### other colo events! 
https://cloud-native.rejekts.io/


### Booths!

This year we will have two locations on the show floor!

The cilium oss booth and the Isovalent booth. Come by and say Hi :) 



### Cilium Talks!

[Blog Post](https://cilium.io/blog/2022/05/06/cilium-kubecon-eu-talks)

[IKEA Private Cloud, eBPF Based Networking, Load Balancing, and Observability with Cilium](https://sched.co/zrPW)

Monday May 16, 2022 10:45 - 11:15 CEST

The digital systems of IKEA are situated in public cloud and private data centers around the world. In this talk we’ll highlight some of the challenges – and opportunities - we faced in setting up a large scale, multi-cluster distributed Kubernetes environment across our data centers. We’ll share how we have used Cilium and its eBPF features to have a better scaling profile, to improve observability and even to replace some of our proprietary load balancers.

Connecting Kubernetes workloads across our BGP network
Protecting multi-tenant workloads with multi-cluster network policy
Cilium support for multi-homed pods * Mimicking availability zones with Cilium ClusterMesh
Use Cilium with XDP, ServiceType Loadbalancer and Ingress to replace our proprietary load balancer fronting workload.
You’ll leave this talk understanding how you can use Cilium and its eBPF capabilities to build and instrument your network and obtain great observability.

[Leveraging Cilium and SRv6 for Telco Networking](https://sched.co/zso2)

Monday May 16, 2022 10:50 - 11:20 CEST

In this session, Daniel Bernier from Bell Canada will be demonstrating how Cilium and its eBPF data plane was extended to support telco networking requirements in a cloud-native way. He will demonstrate how Cilium can provide network segmentation and Multi-VRF support with or without the use of multiple interfaces. With this new approach, he will also explain how to build simple multi-cluster VPN or simple integration to an MPLS provider network by leveraging natively IPv6 and SRv6.

[Connecting Klusters on the edge with deep dive into Cilium Cluster Mesh](https://sched.co/zsAE)

Tuesday May 17, 2022 11:35 - 11:45 CEST

Edge computing can require connecting hundreds of clusters across disparate locations and infrastructures. Without a networking solution to manage this scale and complexity, you will just have a bunch of computers talking to themselves rather than each other and your customers. Cilium is the next generation, eBPF powered open-source Cloud Native Networking solution, providing security, observability, scalability, and superior performance. Cilium has joined the CNCF as an incubating project. In this session, you’ll learn how you can leverage Cilium Cluster Mesh for providing connectivity for load-balancing, observability, and security between nodes across multiple clusters, enabling simple, high-performance cross-cluster connectivity at the edge. We’ll explore how Cluster Mesh allows endpoints in connected clusters to communicate while providing full security policy enforcement. The audience will walk away with an appreciation for how eBPF can help solve their networking challenges at the edge.

[Transparent Live Migration of Services Between Kubernetes Clusters](https://sched.co/ytpo)

Thursday May 19, 2022 15:25 - 16:00 CEST

Operating a distributed database on a single Kubernetes cluster is interesting, but how about transparently migrating it from one cluster to another–potentially between different cloud providers– without impacting user workloads? Kubernetes has become the de facto default deployment for ArangoDB, a distributed Graph database. Consider for example ArangoDB Oasis, a managed Cloud Database service with over 200 deployments (aka highly available database clusters) across three major cloud providers and many regions. But outages, (Kubernetes) upgrades, resource considerations, and cost optimizations require the underlying infrastructure to be very dynamic including migration between Kubernetes cluster, datacenter, or even cloud providers. This talk provides insights into how Kube-Arango, the OSS operator for ArangoDB, supports live migration of distributed stateful applications without impact on users. Challenges in such migration include for example networking, DNS, and persistent data.

[Cilium: Welcome, Vision and Updates](https://sched.co/ytq0)

Thursday May 19, 2022 15:25 - 16:00 CEST

If you’re interested in using Cilium, or contributing to the project, this session is for you. Our agenda for this session: 1. Introduction to Cilium A brief overview of the origin and vision for Cilium. 2. Working with Cilium An end user's perspective of using Cilium. 3. Cilium Service Mesh Cilium can be used as a highly efficient service mesh data plane. Let’s discuss the learnings from our beta, and the upcoming roadmap. We will leave time for Q&A, and an opportunity to meet Cilium maintainers and contributors.

[Choosing Cloud Native Technologies for the Journey to Multi-cloud](https://sched.co/ytpu)

Thursday May 19, 2022 15:25 - 16:00 CEST

Building, deploying and maintaining systems has become increasingly more complicated in recent years. Now, as engineers look toward migrating to multi-cloud architectures, systems and processes may need to be migrated to new technologies. But what choices are available, how do they fit together and how can the CNCF landscape help? This talk discusses the cloud native technologies that can be used to convert to a multi-cloud architecture and highlights some of the lessons learned from taking this journey on at Form3. The audience will learn: - How to decide if multi-cloud is essential for them - The fundamentals of deploying services across multiple clouds with Kubernetes - How to leverage Cilium to mesh together multiple clusters - The basics of event sourcing using NATS in the multi-cloud world - Resilient and performant data storage using CockroachDB This talk is useful for any new comers to the cloud native landscape, as well as those curious about going multi-cloud!

[Kubernetes Networking 101](https://sched.co/ytrV)

Friday May 20, 2022 11:00 - 12:30 CEST

Kubernetes Networking 101 will introduce attendees to the world of network communications in a hands on Cloud Native setting. This talk delivers a high level but completely practical end to end look at service communications within and without a Kubernetes cluster. Attendees will see how the many facets of Kubernetes networking come together to enable powerful communications solutions first hand. The tutorial begins with the simplest types of service communications, using Kubernetes services, DNS (CoreDNS) and CNI plugins (Cilium) to facilitate interprocess communications and load balancing. The tutorial builds additional scenarios on this base, including ingress (Emissary/Envoy), NodePort / HostPort features, load balancing (Metal-lb) and finally a short look at service mesh functionality (Linkerd). Upon completion of this tutorial, attendees will have a clear understanding of the Kubernetes communications possibilities and pointers to next steps in the learning journey.

[Logs Told Us It Was DNS, It Felt Like DNS, It Had To Be DNS, It Wasn’t DNS](https://sched.co/ytrw)

Friday May 20, 2022 11:55 - 12:30 CEST

It all started with a team reaching out because they had DNS issues during rolling updates. Business as usual when you host hundreds of applications on dozens of Kubernetes clusters… Four weeks later: We are reading kernel code to understand the corner cases of dropping Martian packets. Could this be the connection between gRPC client reconnect algorithms and the overflowing conntrack table we can feel but not see? In time, we solved the issue. And for once… it wasn't DNS! In this talk, we will focus on one of the most complex incidents we have faced in our Kubernetes environment. We will go through the debugging steps in detail, dive deep into the mysterious behaviors we discovered and explain how we finally addressed the incident by simply removing three lines of code.

[Better Bandwidth Management with eBPF](https://sched.co/ytsQ)

Friday May 20, 2022 14:00 - 14:35 CEST

Kubernetes provides many knobs for managing common system resources such as vCPUs and memory limits per Pod, but often forgotten is the effect of unbounded network communication in a cluster. A large churn of packets from several services can starve bandwidth for other services. Also, out of the box TCP congestion management is not optimal for Internet-facing services. In this talk we will explore how eBPF can be leveraged to dynamically insert logic for flexible, efficient and scalable rate limiting and bandwidth management on a per-Pod basis. This talk details: - The scalability limits of token bucket filters by the bandwidth plugin, and why EDT (Earliest Departure Time) combined with eBPF is a major step forward. - How TCP congestion control with BBR can now be leveraged for Pods thanks to eBPF for significantly improving application latency and throughput. - The benefits of enforcing bandwidth limits at the egress point and considerations when to use ingress enforcement.

[A Guided Tour of Cilium Service Mesh](https://sched.co/yttj)

Friday May 20, 2022 16:00 - 16:35 CEST

The Cilium project is adding Service Mesh features to its existing eBPF-enabled, identity-aware Kubernetes networking capabilities. This demo-driven talk explores how this works, and shows why it’s now possible to create a service mesh without sidecars. - Demonstrate why, before eBPF, the sidecar model was necessary for accessing an application pod’s network traffic - Explore how Cilium uses eBPF programs to connect Kubernetes endpoints - Show how this makes the sidecar model unnecessary for identity-aware connectivity - Demonstrate an example Cilium Service Mesh in use - Compare the resources used (in both userspace and the kernel) for both models Along the way, this talk will clarify some container and kernel concepts so that attendees can leave with a mental model of how eBPF-enabled service mesh really works.

[Composability is to Software as Compounding Interest is to Finance](https://sched.co/ytqy)

Friday May 20, 2022 14:55 - 15:30 CEST

The cloud native ecosystem is built of composable projects that can be stacked, recombined, reused, and built upon. This composability allows cloud native developers to iterate and ship functionality fast and creates compounding value to businesses from telcos to machine learning to gaming. This talk will trace the history of composability within the cloud native landscape from making Kubernetes pluggable and extensible through the CNI and CRI to standardizing observability with Prometheus and OTel to eBPF making security and networking composable with Cilium. Along the way we will discover how each interface and extension built the value of the project and the ecosystem as a whole creating a learning and business value flywheel. The audience will learn how the composability of cloud native has helped grow the public cloud, generated many successful startups, given meaningful careers to a wide variety of people, and why buying into composable ecosystems compounds business value.



### Virtual Office Hours!
Make sure you stop by the Cilium booth to get your Cilium swag and if you are visiting virtually, we will have virtual office hours on Wednesday 18th May 12:30-13:15 CEST. See you there!




### Kubecon EU!
[Schedule](https://events.linuxfoundation.org/kubecon-cloudnativecon-europe/program/schedule/)

#### Talks that Nicholas is looking forward to:
* [Bypassing Falco: How to Compromise a Cluster without Tripping the SOC - Shay Berkovich, BlackBerry](https://sched.co/ytl7)
* [Kubernetes Steering Committee AMA - Christoph Blecker, Red Hat; Bob Killen, Google; Tim Pepper & Davanum Srinivas, VMware; Paris Pittman, Apple; Stephen Augustus, Cisco](https://sched.co/ytnp)
* [The Hitchhiker's Guide to Pod Security - Lachlan Evenson, Microsoft](https://sched.co/ytnL)
* [Keynote: THE API IS PEOPLE! - Stephen Augustus, Head of Open Source, Cisco](https://sched.co/yuFh)
* [Allyship Workshop: A Human-Centered Approach to Allyship](https://sched.co/10ei9)


#### Talks that Duffie is looking forward to:
* [Cloud Native Security Day!](https://cloudnativesecurityconeu22.sched.com/)
* [Cilium on eks with Liz and Duffie](https://awscontainerdayseurope.splashthat.com/)
* [Make the Secure Kubernetes Supply Chain Work for You - Adolfo García Veytia @puerco ](https://sched.co/ytoz)
* [Too Much to Choose – Making Sense of a Smorgasbord of Security Standards - Anais Urlichs & Rory McCune](https://sched.co/ytsN)
* [Threat Modeling Kubernetes: A Lightspeed Introduction - Lewis Denham-Parry](https://sched.co/ytqj)
* [Keynote: THE API IS PEOPLE! - Stephen Augustus, Head of Open Source, Cisco](https://sched.co/yuFh)
* [The Hitchhiker's Guide to Pod Security - Lachlan Evenson, Microsoft](https://sched.co/ytnL)

