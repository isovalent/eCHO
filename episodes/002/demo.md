Kubernetes cluster running with some demo apps. 

You can see Cilium and Hubble components running in the kube-system namespace. 

`kubectl -n kube-system get pods`

Pick one of the Cilium pods and you can exec into it, and run both `cilium status` and `hubble status` in it

`kubectl exec -it -n kube-system cilium-<id> -- cilium status`

`kubectl exec -it -n kube-system cilium-<id> -- hubble status`

This Cilium pod is running on one node, and can see the traffic flowing on that node (only)

`kubectl exec -it -n kube-system cilium-<id> -- hubble observe`

The Hubble relay combines traffic flows from all the Cilium agents in the cluster.

`kubectl exec -it -n kube-system hubble-relay-<id> -- hubble observe`

Or you can port-forward the hubble relay to a local port. (Run this in the background, or in a separate terminal)

`cilium hubble port-forward` 

Using an environment variable means we won't have to specify `--server localhost:4245` on local hubble CLI commands. 

`export HUBBLE_SERVER=localhost:4245`

Now we can see flows from the whole cluster: 

`hubble observe` 

Note that Hubble is showing Kubernetes pod names. 
