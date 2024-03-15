# Episode 118: Exploring EKS prefix delegation

## Upcoming Events

* [Cilium Advent!](https://labs-map.isovalent.com/holidays/)


## Recent news
* [What's new in 1.14 with Thomas Graf](https://isovalent.com/events/2023-11-30-cilium-114-webinar/)
* [Hand on with 1.14 workshop](https://isovalent.com/events/2023-12-14-cilium-114-workshop/)
* [cilium.io/newsletter](https://cilium.io/newsletter)


---

EKS with Prefix Delegation.

### Why?
[aws docs](https://aws.github.io/aws-eks-best-practices/networking/prefix-mode/index_linux/)

### How?
[Amit Blogpost!](https://medium.com/@amitmavgupta/cilium-support-for-eni-prefix-delegation-in-an-eks-cluster-feddf894160b)

### Let's do it!

eks-config.yaml
``` yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: dcooley-delegation
  region: us-west-1

addons:
  - name: vpc-cni
    version: latest
    configurationValues: |-
      env:
        ENABLE_PREFIX_DELEGATION: "true"
    resolveConflicts: overwrite
managedNodeGroups:
- name: ng-1
  desiredCapacity: 2
  privateNetworking: true
  # taint nodes so that application pods are
  # not scheduled/executed until Cilium is deployed.
  # Alternatively, see the note below.
  taints:
   - key: "node.cilium.io/agent-not-ready"
     value: "true"
     effect: "NoExecute"

```
