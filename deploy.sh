#!/usr/bin/env bash

# deploy.sh <jobNumber>
set -o errexit
set -o nounset
set -o pipefail

JOBNUMBER=$1

# create service
kubectl create -f <(cat << EOF
apiVersion: v1
kind: Service
metadata:
  name: gangservice
  labels:
    app: gang
spec:
  selector:
    app: gang
  type: ClusterIP
  ports:
  - protocol: TCP
    port: 8863
    targetPort: 8863
EOF)

# create job counter web server
kubectl create -f <(cat << EOF
apiVersion: v1
kind: Pod
metadata:
  name: gangweb
  labels:
    app: gang
spec:
  containers:
    - name: gangweb
      image: gangweb:latest
      imagePullPolicy: Never
      command: ["go"]
      args: ["run", "webserver"]
      ports:
        - containerPort: 8863
EOF)

# wait for web server to be running
until grep 'Running' <(kubectl get pod gangweb -o=jsonpath='{.status.phase}'); do
  sleep 1
done

# create gang jobs
for i in $(seq "$JOBNUMBER"); do
  kubectl create -f <(cat << EOF
apiVersion: batch/v1
kind: Job
metadata:
  name: gang-job-$i
spec:
  completions: 20
  parallelism: 20
  template:
    spec:
      containers:
      - name: gang
        image: gang:latest
        imagePullPolicy: Never
        env:
        - name: jobName
          value: gang-job-$i
        - name: serviceName
          value: gangservice
        - name: memberAmount
          value: "20"
        command: ["go"]
        args: ["run", "gang"]
      restartPolicy: Never
EOF)
done