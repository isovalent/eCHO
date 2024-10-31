# Episode 148 - Exploring Cilium with Geneve DSR

[![Episode 148 - Exploring Cilium with Geneve DSR](https://img.youtube.com/vi/e4aLq5xCoJc/0.jpg)](https://www.youtube.com/watch?v=e4aLq5xCoJc&list=PLDg_GiBbAx-mY3VFLPbLHcxo6wUjejAOC&index=17 "Episode 148 - Exploring Cilium with Geneve DSR")

## Headlines

* [newsletter](https://cilium.io/newsletter)

## Agenda

* Bring up a kind cluster and deploy cilium with bgp based load balancing and geneve with DSR.

* Understanding the enforcement point for Network Policy.

container lab topology.

```yaml
name: lab
prefix: ""
topology:
  kinds:
    linux:
      cmd: bash
  nodes:
    router0:
      kind: linux
      image: frrouting/frr:v8.2.2
      labels:
        app: frr
        type: router
      exec:
      # NAT everything in here to go outside of the lab
      - iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE
      # Loopback IP (IP address of the router itself)
      - ip addr add 10.0.0.0/32 dev lo
      # Terminate rest of the 10.0.0.0/8 in here
      - ip route add blackhole 10.0.0.0/8 
      # Boiler plate to make FRR work
      - touch /etc/frr/vtysh.conf
      - touch /var/log/frr.log
      - chown frr:frr /var/log/frr.log
      - sed -i -e 's/bgpd=no/bgpd=yes/g' /etc/frr/daemons
      - /usr/lib/frr/frrinit.sh start
      # FRR configuration
      - >-
         vtysh -c 'conf t'
         -c 'log file /var/log/frr.log'
         -c 'frr defaults datacenter'
         -c 'router bgp 65000'
         -c '  bgp router-id 10.0.0.0'
         -c '  bgp bestpath as-path multipath-relax'
         -c '  no bgp ebgp-requires-policy'
         -c '  neighbor ROUTERS peer-group'
         -c '  neighbor ROUTERS remote-as external'
         -c '  neighbor ROUTERS default-originate'
         -c '  neighbor net0 interface peer-group ROUTERS'
         -c '  neighbor net1 interface peer-group ROUTERS'
         -c '  address-family ipv4 unicast'
         -c '    redistribute connected'
         -c '  exit-address-family'
         -c '!'
    tor0:
      kind: linux
      image: frrouting/frr:v8.2.2
      labels:
        app: frr
        type: tor
      exec:
      - ip link del eth0
      - ip addr add 10.0.0.1/32 dev lo
      - ip addr add 10.0.1.1/24 dev net1
      - ip addr add 10.0.2.1/24 dev net2
      - touch /etc/frr/vtysh.conf
      - touch /var/log/frr.log
      - chown frr:frr /var/log/frr.log
      - sed -i -e 's/bgpd=no/bgpd=yes/g' /etc/frr/daemons
      - /usr/lib/frr/frrinit.sh start
      - >-
         vtysh -c 'conf t'
         -c 'log file /var/log/frr.log'
         -c 'frr defaults datacenter'
         -c 'router bgp 65010'
         -c '  bgp router-id 10.0.0.1'
         -c '  bgp bestpath as-path multipath-relax'
         -c '  no bgp ebgp-requires-policy'
         -c '  neighbor ROUTERS peer-group'
         -c '  neighbor ROUTERS remote-as external'
         -c '  neighbor SERVERS peer-group'
         -c '  neighbor SERVERS remote-as internal'
         -c '  neighbor net0 interface peer-group ROUTERS'
         -c '  neighbor 10.0.1.2 peer-group SERVERS'
         -c '  neighbor 10.0.2.2 peer-group SERVERS'
         -c '  address-family ipv4 unicast'
         -c '    redistribute connected'
         -c '  exit-address-family'
         -c '!'
    tor1:
      kind: linux
      image: frrouting/frr:v8.2.2
      labels:
        app: frr
        type: tor
      exec:
      - ip link del eth0
      - ip addr add 10.0.0.2/32 dev lo
      - ip addr add 10.0.3.1/24 dev net1
      - ip addr add 10.0.4.1/24 dev net2
      - touch /etc/frr/vtysh.conf
      - touch /var/log/frr.log
      - chown frr:frr /var/log/frr.log
      - sed -i -e 's/bgpd=no/bgpd=yes/g' /etc/frr/daemons
      - /usr/lib/frr/frrinit.sh start
      - >-
         vtysh -c 'conf t'
         -c 'log file /var/log/frr.log'
         -c 'frr defaults datacenter'
         -c 'router bgp 65011'
         -c '  bgp router-id 10.0.0.2'
         -c '  bgp bestpath as-path multipath-relax'
         -c '  no bgp ebgp-requires-policy'
         -c '  neighbor ROUTERS peer-group'
         -c '  neighbor ROUTERS remote-as external'
         -c '  neighbor SERVERS peer-group'
         -c '  neighbor SERVERS remote-as internal'
         -c '  neighbor net0 interface peer-group ROUTERS'
         -c '  neighbor 10.0.3.2 peer-group SERVERS'
         -c '  neighbor 10.0.4.2 peer-group SERVERS'
         -c '  address-family ipv4 unicast'
         -c '    redistribute connected'
         -c '  exit-address-family'
         -c '!'
    lab:
      kind: k8s-kind
      startup-config: ./cluster.yaml
      extras:
        k8s_kind:
          deploy:
            wait: 0s

    lab-control-plane:
      kind: ext-container
      exec:
      # Cilium currently doesn't support BGP Unnumbered
      - ip addr add 10.0.1.2/24 dev net0
      # Cilium currently doesn't support importing routes
      - ip route replace default via 10.0.1.1
    lab-worker:
      kind: ext-container
      exec:
      - ip addr add 10.0.2.2/24 dev net0
      - ip route replace default via 10.0.2.1
    lab-worker2:
      kind: ext-container
      exec:
      - ip addr add 10.0.3.2/24 dev net0
      - ip route replace default via 10.0.3.1
    lab-worker3:
      kind: ext-container
      exec:
      - ip addr add 10.0.4.2/24 dev net0
      - ip route replace default via 10.0.4.1


  links:
  - endpoints: ["router0:net0", "tor0:net0"]
  - endpoints: ["router0:net1", "tor1:net0"]
  - endpoints: ["tor0:net1", "lab-control-plane:net0"]
  - endpoints: ["tor0:net2", "lab-worker:net0"]
  - endpoints: ["tor1:net1", "lab-worker2:net0"]
  - endpoints: ["tor1:net2", "lab-worker3:net0"]
```

Makefile
```sh

NAME=clab
VERSION=1.15.6
GW_API_VERSION=release-1.1
GATEWAY=$(shell docker exec router0 hostname -i)
HELM_REPO="cilium"

deploy:
  sudo containerlab -t topo.yaml deploy
	kind get kubeconfig --name $(NAME) > ~/.kube/config
	kubectl apply -k https://github.com/kubernetes-sigs/gateway-api/config/crd/experimental/?ref=$(GW_API_VERSION)

cilium-geneve:
	cilium install \
		--set bgpControlPlane.enabled=true \
		--set bpf.masquerade=true \
		--set egressGateway.enabled=true \
		--set envoy.enabled=true \
		--set gatewayAPI.enabled=true \
		--set hubble.enabled=true \
		--set hubble.relay.enabled=true \
		--set ingressController.enabled=true \
		--set ingressController.service.allocateLoadBalancerNodePorts=false \
		--set k8s.requireIPv4PodCIDR=true \
		--set loadBalancer.mode=dsr \
		--set loadBalancer.dsrDispatch=geneve \
		--set tunnelProtocol=geneve \
		--dry-run-helm-values > values.yaml
	helm install --kube-context kind-$(NAME) -n kube-system cilium $(HELM_REPO)/cilium --version $(VERSION) -f values.yaml

app-deploy:
	while ! kubectl get lbippools ; do sleep 1 ; done
	kubectl create -k https://github.com/mauilion/gw-api-demo

up: deploy cilium-geneve app-deploy apply-policy routes

reload: destroy deploy

destroy: 
	sudo containerlab -t topo.yaml destroy --cleanup
	sudo rm .topo.yaml.bak

apply-policy:
	kubectl apply -f cilium-bgp-peering-policies.yaml

routes:
	sudo ip route replace 20.0.10.0/24 via $(GATEWAY)
	sudo ip route replace 30.0.10.0/24 via $(GATEWAY)
	sudo ip route replace 40.0.10.0/24 via $(GATEWAY)


show-rib:
	@echo "======== router0 ========"
	docker exec -it router0 vtysh -c 'show bgp ipv4 wide'
	@echo "======== tor0    ========"
	docker exec -it tor0 vtysh -c 'show bgp ipv4 wide'
	@echo "======== tor1    ========"
	docker exec -it tor1 vtysh -c 'show bgp ipv4 wide'

show-fib:
	@echo "======== router0 ========"
	docker exec -it router0 ip r
	@echo "======== tor0    ========"
	docker exec -it tor0 ip r
	@echo "======== tor1    ========"
	docker exec -it tor1 ip r

show-neighbors:
	@echo "======== router0 ========"
	docker exec -it router0 vtysh -c 'show bgp ipv4 summary wide'
	@echo "======== tor0    ========"
	docker exec -it tor0 vtysh -c 'show bgp ipv4 summary wide'
	@echo "======== tor1    ========"
	docker exec -it tor1 vtysh -c 'show bgp ipv4 summary wide'

show-bgp:
	@echo "======== router0 ========"
	docker exec -it router0 vtysh -c 'show bgp ipv4'
	@echo "======== tor0    ========"
	docker exec -it tor0 vtysh -c 'show bgp ipv4'
	@echo "======== tor1    ========"
	docker exec -it tor1 vtysh -c 'show bgp ipv4'

```
