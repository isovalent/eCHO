# Episode 51: Life of a Packet.

[YouTube](https://youtu.be/0BKU6avwS98)

# Headlines 

- [Liz on Kubecuddle](https://twitter.com/kubecuddlepod/status/1542544065966723074?s=20&t=8UxRx6Zmaya_Wo05vnlc4w)
- [ebpf Summit CFP is open!](https://twitter.com/eBPFsummit/status/1542625296708501509?s=20&t=orlHygWLoRvvL8LgcqNiRA)
- [Echo News episode 9!](https://isovalent-9197153.hs-sites.com/echo-news-episode-9?)

# Life of a packet.
Understanding how Cilium does what it does.

[Docs](https://docs.cilium.io/en/v1.11/concepts/ebpf/lifeofapacket/)


## Scenarios:
- [x] - pod 2 pod same node 
- [x] - pod 2 pod on adjacent node
- [x] - pod 2 proxy 2 pod on adjacent node 

## Tools:

cilium monitor
hubble observe
hubble relay


## Observation Points:

* to-endpoint
* to-overlay
* to-proxy