With [Duffie Cooley](https://twitter.com/mauilion)

## Headlines

* [eBPF Summit 2023](https://ebpf.io/summit-2023/)
* [Echo Newsletter](https://cilium.io/newsletter)


# Migrating an aks cluster from kubenet to cilium cni

## Stand up a cluster with kubenet.
### Inspect kubenet cluster.
### deploy sample application.
### Generate Cilium config.
### deploy cilium.
### Validate sample application.
### Migrate workloads.
### Validate.
### Complete the migration.


`setup.sh`
```
# Global Variables
RESOURCE_GROUP_NAME=KubenetTest
LOCATION=westus

VNET_NAME=KubenetVNet
ADDRESS_PREFIXES="192.168.0.0/16"
SUBNET_NAME=KubenetSubnet
SUBNET_PREFIX="192.168.1.0/24"
CLUSTER_NAME=AKSKubenet
IDENTITY=KubeNetIdentity

#Create the resourcegroup
az group create --name $RESOURCE_GROUP_NAME --location $LOCATION

#Create an identity
az identity create --name $IDENTITY --resource-group $RESOURCE_GROUP_NAME

#define the PRINCIPAL_ID
PRINCIPAL_ID=$(az identity show --name $IDENTITY --resource-group $RESOURCE_GROUP_NAME --query principalId -o tsv)

#Create the network
az network vnet create \
    --resource-group $RESOURCE_GROUP_NAME \
    --name $VNET_NAME \
    --address-prefixes $ADDRESS_PREFIXES \
    --subnet-name $SUBNET_NAME \
    --subnet-prefix $SUBNET_PREFIX
#Grab the SUBNET_ID
SUBNET_ID=$(az network vnet subnet show --resource-group $RESOURCE_GROUP_NAME --vnet-name $VNET_NAME --name $SUBNET_NAME --query id -o tsv)

#Create the AKS Cluster.
az aks create \
    --resource-group $RESOURCE_GROUP_NAME \
    --name  $CLUSTER_NAME \
    --network-plugin kubenet \
    --node-count 3 \
    --service-cidr 10.0.0.0/16 \
    --dns-service-ip 10.0.0.10 \
    --pod-cidr 10.244.0.0/16 \
    --vnet-subnet-id $SUBNET_ID
#Update kubeconfig to point to it.
az aks get-credentials --overwrite-existing --name $CLUSTER_NAME --resource-group $RESOURCE_GROUP_NAME

#Generate the cilium config.
cilium install --dry-run-helm-values --set nodinit.enabled=true,azure.resourceGroup=$RESOURCE_GROUP_NAME,ipam.mode=Kubernetes,hubble.enabled=true,hubble.relay.enabled=true > values.yaml

#install cilium.
helm install cilium -n kube-system --version v1.13.4 cilium/cilium --values values.yaml


#cleanup az group delete $RESOURCE_GROUP_NAME



```
