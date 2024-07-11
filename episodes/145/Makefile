#PROXY := $(shell kubectl get svc --namespace kong kong-cp-kong-proxy -o jsonpath='{.status.loadBalancer.ingress[0].ip}'):$(shell kubectl get svc --namespace kong kong-cp-kong-proxy -o jsonpath='{.spec.ports[0].port}')
#ADMIN := $(shell kubectl get svc --namespace kong kong-cp-kong-admin -o jsonpath='{.status.loadBalancer.ingress[0].ip}'):$(shell kubectl get svc --namespace kong kong-cp-kong-admin -o jsonpath='{.spec.ports[0].port}')

#:kubectl get svc --namespace kong kong-cp-kong-proxy -o jsonpath='{.spec.ports[0].port}')
#export PROXY_IP=${HOST}:${PORT}

alias:
	alias proxy_curl="docker exec -it kind-control-plane curl $(shell kubectl get svc --namespace kong kong-cp-kong-proxy -o jsonpath='{.status.loadBalancer.ingress[0].ip}'):$(shell kubectl get svc --namespace kong kong-cp-kong-proxy -o jsonpath='{.spec.ports[0].port}')"
	alias admin_curl="docker exec -it kind-control-plane curl -k https://$(shell kubectl get svc --namespace kong kong-cp-kong-admin -o jsonpath='{.status.loadBalancer.ingress[0].ip}'):$(shell kubectl get svc --namespace kong kong-cp-kong-admin -o jsonpath='{.spec.ports[0].port}')"

admin:
	socat TCP4-LISTEN:8002,fork TCP4:172.18.0.3:$(shell kubectl get svc --namespace kong kong-cp-kong-manager -o jsonpath='{.spec.ports[0].nodePort}')| socat TCP4-LISTEN:8001,fork TCP4:172.18.0.3:$(shell kubectl get svc --namespace kong kong-cp-kong-admin -o jsonpath='{.spec.ports[0].nodePort}')

kind:
	echo "Standing up cluster"
	kind create cluster --image kindest/node:v1.30.0 --config ./kind.yaml

crds:
	kubectl apply -f https://raw.githubusercontent.com/kubernetes-sigs/gateway-api/v1.0.0/config/crd/standard/gateway.networking.k8s.io_gatewayclasses.yaml
	kubectl apply -f https://raw.githubusercontent.com/kubernetes-sigs/gateway-api/v1.0.0/config/crd/standard/gateway.networking.k8s.io_gateways.yaml
	kubectl apply -f https://raw.githubusercontent.com/kubernetes-sigs/gateway-api/v1.0.0/config/crd/standard/gateway.networking.k8s.io_httproutes.yaml
	kubectl apply -f https://raw.githubusercontent.com/kubernetes-sigs/gateway-api/v1.0.0/config/crd/standard/gateway.networking.k8s.io_referencegrants.yaml
	kubectl apply -f https://raw.githubusercontent.com/kubernetes-sigs/gateway-api/v1.0.0/config/crd/experimental/gateway.networking.k8s.io_grpcroutes.yaml
	kubectl apply -f https://raw.githubusercontent.com/kubernetes-sigs/gateway-api/v1.0.0/config/crd/experimental/gateway.networking.k8s.io_tlsroutes.yaml

cilium:
	helm install cilium cilium/cilium --version 1.15.6 \
		  --namespace kube-system \
		  --set ipam.mode=kubernetes \
		  --set kubeProxyReplacement=true \
		  --set k8sServiceHost=kind-control-plane \
		  --set k8sServicePort=6443 \
		  --set hostFirewall.enabled=true \
		  --set l2announcements.enabled=true \
		  --set externalIPs.enabled=true \
		  --set gatewayAPI.enabled=true \
		  --set envoy.enabled=true \
		  --set debug.enabled=true \
		  --set debug.verbose=flow
	cilium config view | grep "enable-gateway-api"
	@echo "Wait 5 to 10 seconds before applying layer2"

layer2:
	kubectl apply -f ./layer2.yaml

kong:
	helm repo add kong https://charts.konghq.com
	helm repo update
	kubectl create namespace kong
	kubectl create secret generic kong-enterprise-license --from-literal=license="'{}'" -n kong
	openssl req -new -x509 -nodes -newkey ec -pkeyopt ec_paramgen_curve:secp384r1 -keyout ./tls.key -out ./tls.crt -days 1095 -subj "/CN=kong_clustering"
	kubectl create secret tls kong-cluster-cert --cert=./tls.crt --key=./tls.key -n kong
	helm install kong-cp kong/kong -n kong --values ./kong-values-cp.yaml
	helm install kong-dp kong/kong -n kong --values ./kong-values-dp.yaml

gateway:
	kubectl apply -f https://raw.githubusercontent.com/cilium/cilium/1.15.6/examples/kubernetes/gateway/basic-http.yaml
demo:
	kubectl apply -f https://raw.githubusercontent.com/istio/istio/release-1.11/samples/bookinfo/platform/kube/bookinfo.yaml

clean:
	kind delete cluster
	rm tls.crt
	rm tls.key
