# integration-k8s-gke

Integration K8s GKE runs NSM system tests on GKE.

[cloudtest](https://github.com/networkservicemesh/cloudtest) is used to run the tests from [deployments-k8s](https://github.com/networkservicemesh/deployments-k8s/) in GKE.

You can see exactly what cloudtest does to setup a cluster in GKE [here](cloudtest/gke.yaml).

Effectively it just sets the indicated environment variables
```bash
GKE_PROJECT_ID
CLUSTER_RULES_PREFIX=gke
GKE_CLUSTER_NAME
KUBECONFIG
GKE_CLUSTER_ZONE
GKE_CLUSTER_TYPE
GKE_CLUSTER_NUM_NODES
GCLOUD_SERVICE_KEY
GCLOUD_PROJECT_ID
GITHUB_RUN_NUMBER
```

and then runs the [gke-start.sh](scripts/gke-start.sh)

## Setup GKE cluster

```bash
gcloud container clusters create "${GKE_CLUSTER_NAME}" --project="${GKE_PROJECT_ID}" --machine-type="${GKE_CLUSTER_TYPE}" --num-nodes="${GKE_CLUSTER_NUM_NODES}" --zone="${GKE_CLUSTER_ZONE}" -q
gcloud container clusters get-credentials "${GKE_CLUSTER_NAME}" --project="${GKE_PROJECT_ID}" --zone="${GKE_CLUSTER_ZONE}"
kubectl create clusterrolebinding cluster-admin-binding --clusterrole cluster-admin --user "$(gcloud config get-value account)"
 
```

## Destroy GKE cluster

```bash
gcloud container clusters delete "${GKE_CLUSTER_NAME}" --project="${GKE_PROJECT_ID}" --zone="${GKE_CLUSTER_ZONE}" -q
```