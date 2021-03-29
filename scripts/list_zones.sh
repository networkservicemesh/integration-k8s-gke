#!/bin/bash

gcloud compute zones list --uri --project="$1" | grep -v asia | grep -v australia | cut -f 9 -d '/'
