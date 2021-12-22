#!/bin/bash

if ! gcloud container clusters list --project="${GKE_PROJECT_ID}" | grep -q ^"${GKE_CLUSTER_NAME}"; then \
  time gcloud container clusters create "${GKE_CLUSTER_NAME}" --project="${GKE_PROJECT_ID}" --machine-type="${GKE_CLUSTER_TYPE}" --num-nodes="${GKE_CLUSTER_NUM_NODES}" --zone="${GKE_CLUSTER_ZONE}" --node-version=1.20.10-gke.1600 -q; \
  echo "Writing config to ${KUBECONFIG}"; \
  gcloud container clusters get-credentials "${GKE_CLUSTER_NAME}" --project="${GKE_PROJECT_ID}" --zone="${GKE_CLUSTER_ZONE}" ; \
  gcloud compute firewall-rules create allow-proxy-nsm --action ALLOW --rules tcp:80 --project="${GKE_PROJECT_ID}"; \
  gcloud compute firewall-rules create allow-nsm --action ALLOW --rules tcp:5000-5100 --project="${GKE_PROJECT_ID}"; \
  gcloud compute firewall-rules create allow-vxlan --action ALLOW --rules udp:4789 --project="${GKE_PROJECT_ID}"; \
  gcloud compute firewall-rules create allow-wireguard --action ALLOW --rules udp:51820-52000 --project="${GKE_PROJECT_ID}"; \
  kubectl create clusterrolebinding cluster-admin-binding \
    --clusterrole cluster-admin \
    --user "$(gcloud config get-value account)"; \
fi
