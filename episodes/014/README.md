# Episode 14: netkat with @Itsuugo

[YouTube](https://youtu.be/yabzjJMdI08)
16 July 2021

hosted by [Duffie Cooley](https://twitter.com/mauilion)

## Headlines

* [Detecting Kernel Hooking using eBPF](https://blog.tofile.dev/2021/07/07/ebpf-hooks.html)
* [Facebook Enforcing encryption at scale](https://engineering.fb.com/2021/07/12/security/enforcing-encryption/)
* [Halving the TezEdge node’s memory usage with an eBPF-based memory profiler](https://medium.com/tezedge/halving-the-tezedge-nodes-memory-usage-with-an-ebpf-based-memory-profiler-2bfd32f94f69)				More context: search for tezedge in https://ebpf.io/blog/ebpf-updates-2020-12					
* eBPF Summit on [Twitter](https://twitter.com/ebpfsummit) and [LinkedIn](https://linkedin.com/company/ebpf-summit)


## netkat with [Antonio Ojea](https://twitter.com/itsuugo) 
##### Twitter: @itsuugo Github: @aojea 

netkat is a netcat clone that uses raw sockets to avoid iptables and/or other OS filtering mechanisms.

### Motivation

Kubernetes environments have a LOT of iptables, hence the necessity to bypass iptables sometimes, for testing, debugging, troubleshooting, ...

![](https://i.imgur.com/ZwEMrbQ.jpg =400x)

### How it works

```
sudo ./bin/netkat
Usage: nk [options] [hostname] [port]

  -debug
        Debug
  -interface string
        Specify interface to use. Default interface with the route to the specified [hostname]
  -listen
        Bind and listen for incoming connections
  -source-port int
        Specify source port to use
  -udp
        Use UDP instead of default TCP
```

netkat works as the original netcat, but using raw sockets, so it is not affected by the netfilter hooks.

A RAW socket receives a copy of all packets, however, we don't need all of them, just the ones used by the netkat connection.

On the socket we filter the packets that we don't WANT
On ingress we filter the packets that we WANT, so the host doesn't close our connection

Netkat can obtain the tuple of the connection in advance, so it can filter using BFP and the socket SO_ATTACH_FILTER option, and the traffic control eBPF filtering capability.

The packets received on the socket, are raw packets (with ethernet headers), but since they are bypassing the kernel TCP/IP stack, something need to reassemble the packets and obtain a data stream. This is achieved using an userspace TCP/IP stack https://pkg.go.dev/gvisor.dev/gvisor/pkg/tcpip/stack

At high level, the code does:
- obtain the connection details: source and destination IP and port, and host interface
- parametrize the eBPF code with the connection details, bpf2go generate the eBPF code
- inject the BPF code to the raw socket and to the ingress interface
- 
 
As result, we have a new version of netcat that bypass iptables (needs CAP_NET_RAW)
It has a bonus, with the -d flag it has a sniffer XD
```
sudo ./netkat -d -l 127.0.0.1 80
2021/07/16 01:58:08 Creating raw socket
2021/07/16 01:58:08 Adding ebpf ingress filter on interface lo
2021/07/16 01:58:08 filter {LinkIndex: 1, Handle: 0:1, Parent: ffff:fff2, Priority: 0, Protocol: 3}
2021/07/16 01:58:08 Creating user TCP/IP stack
2021/07/16 01:58:08 Listening on 127.0.0.1:80
I0716 01:58:31.206864   22926 sniffer.go:418] recv tcp 127.0.0.1:44644 -> 127.0.0.1:80 len:0 id:a5ee flags:  S     seqnum: 1321663292 ack: 0 win: 65495 xsum:0xfe30 options: {MSS:65495 WS:7 TS:true TSVal:2715804214 TSEcr:0 SACKPermitted:true}
I0716 01:58:31.206937   22926 sniffer.go:418] recv tcp 127.0.0.1:44644 -> 127.0.0.1:80 len:0 id:a5ee flags:  S     seqnum: 1321663292 ack: 0 win: 65495 xsum:0xfe30 options: {MSS:65495 WS:7 TS:true TSVal:2715804214 TSEcr:0 SACKPermitted:true}
I0716 01:58:32.239938   22926 sniffer.go:418] recv tcp 127.0.0.1:44644 -> 127.0.0.1:80 len:0 id:a5ef flags:  S     seqnum: 1321663292 ack: 0 win: 65495 xsum:0xfe30 options: {MSS:65495 WS:7 TS:true TSVal:2715805247 TSEcr:0 SACKPermitted:true}
I0716 01:58:32.240044   22926 sniffer.go:418] recv tcp 127.0.0.1:44644 -> 127.0.0.1:80 len:0 id:a5ef flags:  S     seqnum: 1321663292 ack: 0 win: 65495 xsum:0xfe30 options: {MSS:65495 WS:7 TS:true TSVal:2715805247 TSEcr:0 SACKPermitted:true}
I0716 01:58:34.287962   22926 sniffer.go:418] recv tcp 127.0.0.1:44644 -> 127.0.0.1:80 len:0 id:a5f0 flags:  S     seqnum: 1321663292 ack: 0 win: 65495 xsum:0xfe30 options: {MSS:65495 WS:7 TS:true TSVal:2715807295 TSEcr:0 SACKPermitted:true}
I0716 01:58:34.288096   22926 sniffer.go:418] recv tcp 127.0.0.1:44644 -> 127.0.0.1:80 len:0 id:a5f0 flags:  S     seqnum: 1321663292 ack: 0 win: 65495 xsum:0xfe30 options: {MSS:65495 WS:7 TS:true TSVal:2715807295 TSEcr:0 SACKPermitted:true}
^C2021/07/16 01:58:35 Exiting: received signal
2021/07/16 01:58:35 Done
```

![](https://i.imgur.com/jLOApyQ.jpg)

### Tests!! 

I really like that Antonio wrote tests for this let's check them out. 



### Demo: Kubernetes LoadBalancer Services and ExternalTrafficPolicy

![](https://i.imgur.com/S9hMwEh.jpg)

1. Create a multinode cluster with [KIND](https://kind.sigs.k8s.io/)
```
cat <<EOF | kind create cluster --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
- role: worker
- role: worker
EOF
```

2. Run a fake loadbalancer controller
```sh
go get -v github.com/aojea/networking-controllers
cd $GOPATH/github.com/aojea/networking-controllers
cd cmd/loadbalancer
go build .
kind get kubeconfig > kind.conf
./loadbalancer --kubeconfig kind.conf --iprange "10.111.111.0/24"
I0716 02:58:43.535494   15454 controller.go:80] Starting controller fake-loadbalancer-controller
I0716 02:58:43.535567   15454 controller.go:84] Waiting for informer caches to sync
I0716 02:58:43.535589   15454 shared_informer.go:240] Waiting for caches to sync for fake-loadbalancer-controller
I0716 02:58:43.636251   15454 shared_informer.go:247] Caches are synced for fake-loadbalancer-controller 
I0716 02:58:43.636267   15454 controller.go:90] Starting workers
I0716 02:58:43.636285   15454 controller.go:149] Processing sync for service 
```

3. Create a deployment with one pod and expose it with a Loadbalancer
```sh
kubectl apply -f https://gist.githubusercontent.com/aojea/369ad6a5d4cbb6b0fbcdd3dd909d9887/raw/0ecd0564c1db4c7103de07accce99da1c7bf91c3/loadbalancer.yaml
```
4. Obtain the LoadBalancer assigned IP 
```sh
kubectl get service lb-service -o jsonpath='{.status.loadBalancer.ingress[0].ip}'
10.111.111.129
```
The loadbalancer has `externalTrafficPolicy: Local` so it should only be able to work
if the traffic is directed to the nodes that contains pods belonging to that Service.

```sh
kubectl get pods -l app=MyApp -o wide
NAME                               READY   STATUS    RESTARTS   AGE   IP           NODE           NOMINATED NODE   READINESS GATES
test-deployment-6f45696bd5-xlp96   1/1     Running   0          10m   10.244.2.2   kind-worker2   <none>           <none>
```


5. Emulate the external loadbalancer installing a route in the host to the loadbalancer IP through the node with the pod

Get the nodes IPs 
```sh
kubectl get nodes -o wide
NAME                 STATUS   ROLES                  AGE   VERSION   INTERNAL-IP   EXTERNAL-IP   OS-IMAGE       KERNEL-VERSION            CONTAINER-RUNTIME
kind-control-plane   Ready    control-plane,master   34m   v1.21.2   172.18.0.4    <none>        Ubuntu 21.04   4.18.0-301.1.el8.x86_64   containerd://1.5.2
kind-worker          Ready    <none>                 33m   v1.21.2   172.18.0.3    <none>        Ubuntu 21.04   4.18.0-301.1.el8.x86_64   containerd://1.5.2
kind-worker2         Ready    <none>                 33m   v1.21.2   172.18.0.2    <none>        Ubuntu 21.04   4.18.0-301.1.el8.x86_64   containerd://1.5.2
```

> example: pod is in kind-worker2 with ip 172.18.0.2 and loadbalancer IP 10.111.111.129
```sh
sudo ip route add 10.111.111.129 via 172.18.0.2
```
6. Check it preserves the source ip ...
```sh
echo -e "GET /clientip HTTP/1.1\nhost:myhost\n" | nc 10.111.111.129 80
HTTP/1.1 200 OK
Date: Thu, 15 Jul 2021 23:30:35 GMT
Content-Length: 16
Content-Type: text/plain; charset=utf-8

172.18.0.1:34028
```
7. and that it doesn't work if we target a node without pods backing that service
```
sudo ip route del 10.111.111.129 via 172.18.0.2
sudo ip route add 10.111.111.129 via 172.18.0.3

echo -e "GET /clientip HTTP/1.1\nhost:myhost\n" | nc -v -w 3 10.111.111.129 80
Ncat: Version 7.70 ( https://nmap.org/ncat )
Ncat: Connection timed out.
```
8. However, if we use `externalTrafficPolicy: Cluster` it works ... but the source IP is not preserved
```sh
kubectl patch service lb-service -p '{"spec":{"externalTrafficPolicy":"Cluster"}}'
service/lb-service patched
$ echo -e "GET /clientip HTTP/1.1\nhost:myhost\n" | nc -v -w 3 10.111.111.129 80
Ncat: Version 7.70 ( https://nmap.org/ncat )
Ncat: Connected to 10.111.111.129:80.
HTTP/1.1 200 OK
Date: Thu, 15 Jul 2021 23:34:22 GMT
Content-Length: 16
Content-Type: text/plain; charset=utf-8

172.18.0.3:14569
Ncat: 36 bytes sent, 133 bytes received in 0.01 seconds.
```
9. Set `externalTrafficPolicy: Local` again
```sh
kubectl patch service lb-service -p '{"spec":{"externalTrafficPolicy":"Local"}}'
service/lb-service patched
$ echo -e "GET /clientip HTTP/1.1\nhost:myhost\n" | nc -v -w 3 10.111.111.129 80
Ncat: Version 7.70 ( https://nmap.org/ncat )
Ncat: Connection timed out.
```
10. Test the service from within node without the pod
```sh
echo -e "GET /clientip HTTP/1.1\nhost:myhost\n" | nc -v -w 3 10.111.111.129 80
Connection to 10.111.111.129 80 port [tcp/http] succeeded!
HTTP/1.1 200 OK
Date: Thu, 15 Jul 2021 23:36:50 GMT
Content-Length: 16
Content-Type: text/plain; charset=utf-8
```
We think it shoudn't work because we are in a node without pods,
however, since the traffic comes from a node within the cluster
the traffic is considered internal and externalTrafficPolicy doesn't apply.
The reason is that kube-proxy installs some iptables rules
that capture the traffic from within the node.

This causes a problem to test LoadBalancers, because it requires traffic to be
sent from an external node, something that is not easy:
Kubernetes cluster use to be isolated, expose apiserver endpoint
e2e should not assume direct connectivity from the e2e binary

### Bypassing iptables with netkat

1. Run netkat in a pod (using hostnetwork) we are going to simulate an external connecting from inside a node
```sh
 kubectl apply -f https://gist.githubusercontent.com/aojea/952c82a58da625fbd9b8aca35f0e63f1/raw/0c5c625e037eb6431947c93152a3182c47004aa3/netkat.yaml
pod/netkat created
```
2. Check where is running
```sh
$ kubectl get pods -o wide
NAME                               READY   STATUS    RESTARTS   AGE   IP           NODE           NOMINATED NODE   READINESS GATES
netkat                             1/1     Running   0          9s    172.18.0.3   kind-worker    <none>           <none>
test-deployment-6f45696bd5-xlp96   1/1     Running   0          63m   10.244.2.2   kind-worker2   <none>           <none>
```
3. Login to the container

```sh
kubectl exec -it netkat ash
kubectl exec [POD] [COMMAND] is DEPRECATED and will be removed in a future version. Use kubectl kubectl exec [POD] -- [COMMAND] instead.
/ # ip route add 10.111.111.15 via 172.18.0.2
echo -e "GET /clientip HTTP/1.1\nhost:myhost\n" | ./netkat 10.111.111.15 80
2021/07/16 01:06:54 routes {Ifindex: 18 Dst: 10.111.111.15/32 Src: 172.18.0.3 Gw: 172.18.0.2 Flags: [] Table: 254}
2021/07/16 01:06:54 Creating raw socket
2021/07/16 01:06:54 Adding ebpf ingress filter on interface eth0
2021/07/16 01:06:54 filter {LinkIndex: 18, Handle: 0:1, Parent: ffff:fff2, Priority: 0, Protocol: 3}
2021/07/16 01:06:54 Creating user TCP/IP stack
2021/07/16 01:06:54 Dialing ...
2021/07/16 01:06:54 Connection established
2021/07/16 01:06:54 Connection error: <nil>
HTTP/1.1 200 OK
Date: Fri, 16 Jul 2021 01:06:54 GMT
Content-Length: 16
Content-Type: text/plain; charset=utf-8
172.18.0.3:26913
```

we can see it preserves the source ip

4. Check that if we target a different node it doesn't work as expected because the service has externalTrafficPolicy: local # and it considers the connection as external
```sh
/ # ip route del 10.111.111.15 via 172.18.0.3
/ # ip route add 10.111.111.15 via 172.18.0.4
/ # echo -e "GET /clientip HTTP/1.1\nhost:myhost\n" | ./netkat 10.111.111.15 80
2021/07/16 01:11:38 routes {Ifindex: 18 Dst: 10.111.111.15/32 Src: 172.18.0.2 Gw: 172.18.0.4 Flags: [] Table: 254}
2021/07/16 01:11:38 Creating raw socket
2021/07/16 01:11:38 Adding ebpf ingress filter on interface eth0
2021/07/16 01:11:38 filter {LinkIndex: 18, Handle: 0:1, Parent: ffff:fff2, Priority: 0, Protocol: 3}
2021/07/16 01:11:38 Creating user TCP/IP stack
2021/07/16 01:11:38 Dialing ...
2021/07/16 01:11:43 Dialing error: context deadline exceeded
2021/07/16 01:11:43 Can't connect to server: context deadline exceeded
```

### References:

1. https://en.wikipedia.org/wiki/Netcat
2. https://github.com/aojea/netkat
3. https://man7.org/linux/man-pages/man7/raw.7.html
4. https://man7.org/linux/man-pages/man7/packet.7.html
5. https://man7.org/linux/man-pages/man8/tc-bpf.8.html
6. https://pkg.go.dev/github.com/cilium/ebpf/cmd/bpf2go
7. https://docs.cilium.io/en/stable/bpf/
8. https://developers.redhat.com/blog/2021/04/01/get-started-with-xdp
9. https://blog.cloudflare.com/tag/ebpf/
10. https://www.asykim.com/blog/deep-dive-into-kubernetes-external-traffic-policies