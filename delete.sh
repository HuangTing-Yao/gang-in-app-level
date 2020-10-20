#!/usr/bin/env bash

# delete.sh
set -o errexit
set -o nounset
set -o pipefail

kubectl delete job --selector=app=gang
kubectl delete pod gangweb
kubectl delete svc gangservice  