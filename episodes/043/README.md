

# Episode 43: Cilium FQDN Policy Deepdive.


## [Cilium Twitter handle](https://twitter.com/ciliumproject)!
## [Echo News](https://isogo.to/echo-news-4)!




## FQDN Deep Dive.


### Network Policy:

This network policy sets the dns parser to watch for any dns resolution requests for our pods and enforces that they are able to only establish connectivity to echo.k8s.intra.


```
apiVersion: "cilium.io/v2"
kind: CiliumNetworkPolicy
metadata:
  name: "echo"
spec:
  endpointSelector:
    matchLabels:
      app: net
  egress:
  - toFQDNs:
    - matchPattern: "*.k8s.intra"  
  - toEndpoints:
    - matchLabels:
        "k8s:io.kubernetes.pod.namespace": kube-system
        "k8s:k8s-app": kube-dns
    toPorts:
    - ports:
      - port: "53"
        protocol: ANY
      rules:
        dns:
        - matchPattern: "*"
```

## Theory of operation


When you use Cilium Network Policy to apply the fqdn policy above. All of the pods in the default namespace with app=net labels are configured to only allow connections to the ip address that resolves to echo.k8s.intra


`curl http://echo.k8s.intra`
This command will use the resolver available in the pod to resolve the hostname to an ip. This action will cause the ip address to be stored in the fqdn cache on the node where the pod resides.

`curl http://echo.k8s.intra --resolve echo.k8s.intra:80:172.17.0.2`
This command will use curls built in local resolver mechanism. This means that curl will not make a dns call to resolve echo.k8s.intra it will  directly resolve it to whatever the ip address above is.

If we run this command before there is a mapping the flow will be denied by policy as there is no mapping for the fqdn and ip address.

If any pod on the node makes a dns query to resolve this fqdn -> ip mapping all pods on that node can make use of it. 



### Dns Caching
There are two stages of caching an fqdn -> ip address mapping. 

#### Stage 1. 
When the initial mapping is created the source of the record is "lookup". This record has a ttl of an hour. Any observed dns resolution will reset this clock to include another hour. A dns resolution here means that the client has queried the dns record and hubble can observe that. This 1 hour mapping is tunable via the `--tofqdns-min-ttl`

#### Stage 2.
This phase usually applies when after a dns record has been changed and is no longer resolving to the old ip address. This means that hubble will only observe new ip addresses moving forward.

When the mapping expires the source of the record changes from "lookup" to "connection" This means the record is still available able to be used and using it will keep the definition in the fqdn cache. The mapping is subject to removal once the ct record times out and is garbage collected.

### Configuration Defaults:

[link](https://docs.cilium.io/en/v1.11/cmdref/cilium-agent/)

```
      --tofqdns-dns-reject-response-code string              DNS response code for rejecting DNS requests, available options are '[nameError refused]' (default "refused")
      --tofqdns-enable-dns-compression                       Allow the DNS proxy to compress responses to endpoints that are larger than 512 Bytes or the EDNS0 option, if present (default true)
      --tofqdns-endpoint-max-ip-per-hostname int             Maximum number of IPs to maintain per FQDN name for each endpoint (default 50)
      --tofqdns-idle-connection-grace-period duration        Time during which idle but previously active connections with expired DNS lookups are still considered alive (default 0s)
      --tofqdns-max-deferred-connection-deletes int          Maximum number of IPs to retain for expired DNS lookups with still-active connections (default 10000)
      --tofqdns-min-ttl int                                  The minimum time, in seconds, to use DNS data for toFQDNs policies. (default 3600 )
      --tofqdns-pre-cache string                             DNS cache data at this path is preloaded on agent startup
      --tofqdns-proxy-port int                               Global port on which the in-agent DNS proxy should listen. Default 0 is a OS-assigned port.
      --tofqdns-proxy-response-max-delay duration            The maximum time the DNS proxy holds an allowed DNS response before sending it along. Responses are sent as soon as the datapath is updated with the new IP information. (default 100ms)
```

## Testing

### DNS:
I wanted to use something simple and quick to reconfigure for the dns records. 

I used the existing Coredns configuration to also serve a dns zone: k8s.intra

Here is an example of the configuration: 

```yaml=
kubectl get configmaps -n kube-system coredns -o yaml
apiVersion: v1
data:
  Corefile: |
    .:53 {
        errors
        health {
           lameduck 5s
        }
        ready
        kubernetes cluster.local in-addr.arpa ip6.arpa {
           pods insecure
           fallthrough in-addr.arpa ip6.arpa
           ttl 30
        }
        prometheus :9153
        forward . /etc/resolv.conf {
           max_concurrent 1000
        }
        cache 30
        loop
        reload
        loadbalance
        hosts /etc/coredns/customdomains.db k8s.intra {
          172.17.0.2 echo.k8s.intra
          fallthrough
        }
    }
kind: ConfigMap
metadata:
  creationTimestamp: "2022-04-18T20:37:20Z"
  name: coredns
  namespace: kube-system
  resourceVersion: "17629"
  uid: e1907f35-e8d3-44f2-975f-74f7dfe8814f

```

This configuration allows for any pod to resolve echo.k8s.intra.

When we reconfigure the configmap we also need to restart the coredns pods to ensure that the change in host -> ip mapping is seen.

The ttl for coredns is ~30s by default.

There are two docker containers both listening on port 80 and 443 that are part of the docker network on the hypervisor where our vms are running. These are serving a simple echo server at 172.17.0.2 and 172.17.0.3 


### Test Environment: 
#### Cluster Setup.

1 Control Plane node
2 worker nodes.

kubeadm based installation.

Cilium install: 

cilium install
cilium hubble enable


#### Namespace default:

I've deployed 2 pods per node so that I can validate the questions below.

[`pernode.sh`](https://gist.githubusercontent.com/mauilion/79ca44b583e52faa33113f2f52cf46d2/raw/b854e53af0163d6e41239d0b2d26378e3e647079/pernode.sh)
```
#!/bin/bash

# The purpose of this script is to deploy to each node in the cluster 2 pods. 
# Each pod will have an env var that shows it's zone.

function netpod () {
  kubectl run net${2}-${1} \
     --image overridden  --labels app=net,pod=net${2}-${1},node=${1}  --overrides \
    '{
      "spec":{
        "hostname": "net'${2}-${1}'",
	      "subdomain": "net",
        "nodeName": "'$1'",
        "containers":[{
          "name":"net",
          "image":"mauilion/debug"
        }]
      }
    }'
}

for worker in $(kubectl get nodes -o name | sed s/node.//)
  do
    for i in {1..2}
      do netpod $worker $i
    done
  done

kubectl create service clusterip net --tcp 8080 
```

Resulting in: 

```shell=
$▶ kubectl get pods  -o wide
NAME            READY   STATUS    RESTARTS   AGE   IP           NODE       NOMINATED NODE   READINESS GATES
net1-cp01       1/1     Running   0          87m   10.0.1.241   cp01       <none>           <none>
net1-worker01   1/1     Running   0          87m   10.0.2.242   worker01   <none>           <none>
net1-worker02   1/1     Running   0          87m   10.0.0.38    worker02   <none>           <none>
net2-cp01       1/1     Running   0          87m   10.0.1.140   cp01       <none>           <none>
net2-worker01   1/1     Running   0          87m   10.0.2.239   worker01   <none>           <none>
net2-worker02   1/1     Running   0          87m   10.0.0.152   worker02   <none>           <none>

```





### Phase 1

In the first phase of our evaluation we defined echo.k8s.intra to resolve to 172.17.0.2

FQDN Policy will watch for the resolution of echo.k8s.intra and record the result. This record can be viewed by execing into the cilium agent on the node where the pod is scheduled and looking at the output of cilium fqdn cache list

```shell=
root@cp01:/home/cilium# cilium fqdn cache list 
Endpoint   Source   FQDN              TTL    ExpirationTime             IPs          
66         lookup   echo.k8s.intra.   3600   2022-04-19T16:21:57.997Z   172.17.0.2
```

We can also observe this dns resolution via hubble cli the resolution of echo.k8s.intra to 172.17.0.2 

```bash=
Apr 19 15:23:46.971: default/net1-cp01:36505 <- kube-system/coredns-586b5b6865-8t6nq:53 to-endpoint FORWARDED (UDP)
Apr 19 15:23:46.972: default/net1-cp01:36657 -> kube-system/coredns-586b5b6865-5v4s7:53 L3-L4 REDIRECTED (UDP)
Apr 19 15:23:46.972: default/net1-cp01:36657 -> kube-system/coredns-586b5b6865-5v4s7:53 to-proxy FORWARDED (UDP)
Apr 19 15:23:46.972: default/net1-cp01:36657 -> kube-system/coredns-586b5b6865-5v4s7:53 dns-request FORWARDED (DNS Query echo.k8s.intra. AAAA)
Apr 19 15:23:46.972: default/net1-cp01:36657 -> kube-system/coredns-586b5b6865-5v4s7:53 dns-request FORWARDED (DNS Query echo.k8s.intra. A)
Apr 19 15:23:46.974: default/net1-cp01:36657 <- kube-system/coredns-586b5b6865-5v4s7:53 dns-response FORWARDED (DNS Answer "172.17.0.2" TTL: 30 (Proxy echo.k8s.intra. A))
Apr 19 15:23:46.974: default/net1-cp01:36657 <- kube-system/coredns-586b5b6865-5v4s7:53 dns-response FORWARDED (DNS Answer  TTL: 4294967295 (Proxy echo.k8s.intra. AAAA))
Apr 19 15:23:46.974: default/net1-cp01:36657 <- kube-system/coredns-586b5b6865-5v4s7:53 to-endpoint FORWARDED (UDP)
Apr 19 15:23:46.975: default/net1-cp01:44754 -> echo.k8s.intra:80 L3-Only FORWARDED (TCP Flags: SYN)
Apr 19 15:23:46.975: default/net1-cp01:44754 -> echo.k8s.intra:80 to-stack FORWARDED (TCP Flags: SYN)
Apr 19 15:23:46.976: default/net1-cp01:44754 <- echo.k8s.intra:80 to-endpoint FORWARDED (TCP Flags: SYN, ACK)
Apr 19 15:23:46.976: default/net1-cp01:44754 -> echo.k8s.intra:80 to-stack FORWARDED (TCP Flags: ACK)
Apr 19 15:23:46.976: default/net1-cp01:44754 -> echo.k8s.intra:80 to-stack FORWARDED (TCP Flags: ACK, PSH)
Apr 19 15:23:46.977: default/net1-cp01:44754 <- echo.k8s.intra:80 to-endpoint FORWARDED (TCP Flags: ACK, PSH)
Apr 19 15:23:46.978: default/net1-cp01:44754 -> echo.k8s.intra:80 to-stack FORWARDED (TCP Flags: ACK, FIN)
Apr 19 15:23:46.978: default/net1-cp01:44754 <- echo.k8s.intra:80 to-endpoint FORWARDED (TCP Flags: ACK, FIN)
Apr 19 15:23:46.978: default/net1-cp01:44754 -> echo.k8s.intra:80 to-stack FORWARDED (TCP Flags: ACK)

```

Note this ip -> hostname mapping in the fqdn cache will be present only on a node where a pod that matches this fqdn policy attempts to resolve echo.k8s.intra

Another pod on the same node can resolve make use of this cached result. The second pod will be able to run: 

`curl http://echo.k8s.intra --resolve echo.k8s.intra:80:172.17.0.2`
and cilium will use the cached dns result to allow this connection. 



### Phase 2.


In phase 2 of our evaluation we will change the A record for echo.k8s.infra to 172.17.0.3 and restart the coredns pods with
`kubectl rollout restart -n kube-system deploy/coredns`

At this point from our test pod we should see that the record has changed. 

We can then show that we continue to be able to connect to the new and old ip addresses leveraging the resolve feature of curl. 

A second pod on the same node will also be able to make use of the fqdn cache known on the node in question. 

This means that from net2-cp01 we can connect to the old or new ip address leveraging the local resolution feature of curl. 

A new connection from a new endpoint on worker02 after the change in A record will only have an entry for the new ip address. This cache is not persistent across nodes. Just across pods on the same node. 



# nginx-ingress CVEs
- [CVE-2021-25745: Ingress-nginx `path` can be pointed to service account token file](https://groups.google.com/g/kubernetes-security-announce/c/7vQrpDZeBlc/m/PiEHmrdXAgAJ)
- [CVE-2021-25746: Ingress-nginx directive injection via annotations](https://groups.google.com/g/kubernetes-security-announce/c/hv2-SfdqcfQ/m/__HJ3bdXAgAJ)
















