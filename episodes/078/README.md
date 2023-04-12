# Episode 78: SCTP

[YouTube](https://youtu.be/2lD86qNHXXI)

## Introductions
* Dan

## Upcoming Events

* Security Con this week (congrats to Liz, and all the other presenting or performing keynotes)
* It's State of Open Source conference next week
* Civo conference also next week
* KCD Amsterdamn 23-24 Feb


## Recent News

- New lab available [Cilium Enterprise: Zero Trust Visibility](https://isovalent.com/labs/cilium-enterprise-zero-trust-visibility/) (all around network policies etc.)
- CNCF ToC election (congratulations to Duffie, Cathy and Nikhita !)
- Kubernetes 1.27 has entered *production readiness freeze* as part of the development cycle!

## SCTP

So what is SCTP.. ` ¯\_(ツ)_/¯`

Well, a lot of people have probably heard about IP, a large number of people are aware of TCP and a good amount of people have heard of UDP.

Effectively SCTP is the good bits of TCP and UDP.. 

Regarding TCP:

- Unicast protocols (providing reliable transport)
- Packets in sequence
- CRC ensuring error correction

![](https://www.computernetworkingnotes.org/images/cisco/ccna-study-guide/csg26-01-tcp-header.png)

Regarding UDP:

- Can support unreliable transport
- Out of order delivery
- Small header 12 bytes to UDPs 8 bytes (providing good throughput)

![](https://www.computernetworkingnotes.org/images/cisco/ccna-study-guide/csg26-02-udp-header.png)

### The Differences

- SCTP is message orientated vs a TCP data stream
- Can handle multiple streams vs a single stream per connection
- SCTP provides multi-homing, allowing an endpoint to use multiple addresses. Should a connection break SCTP will continue to send data from the alternative address, providing fault tolerance (and in some situations removing the need for things like BGP)
- SCTP uses an initial agreed cookie to provide protection against MiTM or masquerated/flooding.
- SCTP can detect dropped or duplicated packets

The combination of all of this provides a transport layer that gives functionality to applications that require "some" of the functionality of both TCP or UDP.

![image alt](https://image3.slideserve.com/5783144/sctp-packet-format-l.jpg)

### Adoption

Well SCTP has struggled for adoption, it requires support within the transport layer of the host. This has adoption in Linux, FreeBSD, Solaris etc.. However Windows doesn't appear to natively support SCTP :-( (neither does macOS)

Network devices such as switches / routers don't need know as they tend to focus just on the IP Headers to determine the destination and where to forward the packets.

### Use Cases

- Signaling transport in 4G and 5G mobile networks
- SIP (Session Initiation Protocol) signaling in VoIP (Voice over IP) systems 

### Cilium support 

SCTP support was introduced in Kubernetes 1.12 (Beta) and was eventually graduated to Stable in Kubernetes 1.20.

With Cilium 1.13, it is now possible for Pods to communicate with other Pods and Services using SCTP, and network policies can be applied to SCTP traffic to provide additional security controls.

### The Demo

```
curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash

helm repo add cilium https://helm.cilium.io/


helm upgrade --install cilium cilium/cilium \
--version 1.13.0-rc5     \
--namespace kube-system     \
--set sctp.enabled=true     \
--set hubble.enabled=true     \
--set hubble.metrics.enabled="{dns,drop,tcp,flow,icmp,http}"     \
--set hubble.relay.enabled=true     \
--set hubble.ui.enabled=true     \
--set hubble.ui.service.type=NodePort     \
--set hubble.relay.service.type=NodePort

cilium status --wait

vim sctp-deploy.yaml

kubectl apply -f ./sctp-deploy.yaml 

kubectl get pod

kubectl expose deployment sctp-deployment --port=9999 --type=NodePort

kubectl get svc/sctp-deployment

PORT=$(kubectl get svc/sctp-deployment -o jsonpath='{.spec.ports[0].nodePort}')
echo $PORT

NODE=192.168.0.192

ncat --sctp $NODE_IP $PORT
```

### Observing SCTP

```
hubble observe --server localhost:31234 --protocol sctp --follow

```

Also UI -> http://192.168.0.192:31235/?namespace=default

### Policies

```
cd policies

kubectl apply -f sctp-client-source-deny.yaml -f sctp-client-source.yaml

kubectl get pods --show-labels

SCTP_POD_IP=$(kubectl get pods -o wide -l=app=sctp -o=jsonpath='{.items[0].status.podIP}')
echo $SCTP_POD_IP

kubectl exec -it sctpclient -- nc $SCTP_POD_IP 9999 --sctp

kubectl exec -it sctpclient-deny -- nc $SCTP_POD_IP 9999 --sctp

cat sctp-net-pol.yaml

kubectl apply -f sctp-net-pol.yaml

kubectl exec -it sctpclient -- nc $SCTP_POD_IP 9999 --sctp

kubectl exec -it sctpclient-deny -- nc $SCTP_POD_IP 9999 --sctp
```

## Wrap up

Again, thanks everyone that has so far gotten involved in acquiring the new Isovalent badges. We're super excited to get Cilium 1.13 into your hands soon.

We've some super exciting guests and demos coming up in the next few episodes of eCHO. If anyone has any requests then we would love to hear from you!

Also, if anyone is at any of the events that we've mentioned then come find us, we love to talk. :D