# Episode 53: Life of a Packet: Services.

{%youtube SGfMEpjq07Q%}

# Headlines 

- [ebpf Summit CFP and registration is open!](https://ebpf.io/summit-2022/)

- [whens it close?](https://twitter.com/breakawaybilly/status/1547878638279876608?s=20&t=Q9Nr9xR2lUis_TqhRsDZNQ)

- [Cilium Standalone Layer 4 Load Balancer XDP](https://cilium.io/blog/2022/04/12/cilium-standalone-L4LB-XDP/)

- [ebpf news episode 10](https://isogo.to/echo-news-10)



# Life of a packet.
Understanding how Cilium does what it does.

[Docs](https://docs.cilium.io/en/v1.11/concepts/ebpf/lifeofapacket/)



## Scenarios: Episode 51: 
- [x] - pod 2 pod same node 
- [x] - pod 2 pod on adjacent node
- [x] - pod 2 proxy 2 pod on adjacent node 


## Scenarios: Episode 53:
- [ ] [kube-proxy](https://www.youtube.com/watch?v=wEyLsEaomfA)
- [ ]
- [ ]
- [ ]
- [ ]


## Tools:

cilium monitor
hubble observe
hubble relay

pernode.sh
```shell=
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
## Observation Points:

* to-endpoint
* to-overlay
* to-proxy