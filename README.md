# integration-k8s-gke

## Setup GKE cluster

```bash
gcloud container clusters create "${GKE_CLUSTER_NAME}" --project="${GKE_PROJECT_ID}" --machine-type="${GKE_CLUSTER_TYPE}" --num-nodes="${GKE_CLUSTER_NUM_NODES}" --zone="${GKE_CLUSTER_ZONE}" -q
gcloud container clusters get-credentials "${GKE_CLUSTER_NAME}" --project="${GKE_PROJECT_ID}" --zone="${GKE_CLUSTER_ZONE}"
kubectl create clusterrolebinding cluster-admin-binding --clusterrole cluster-admin --user "$(gcloud config get-value account)"
```

To have **AF_XDP** support, you need to use [gVNIC](https://cloud.google.com/compute/docs/networking/using-gvnic). Add `--enable-gvnic` to the cluster creation command.

## Destroy GKE cluster

```bash
gcloud container clusters delete "${GKE_CLUSTER_NAME}" --project="${GKE_PROJECT_ID}" --zone="${GKE_CLUSTER_ZONE}" -q
```