#!/bin/bash
if [ -z "$1" ]
  then
    echo "Please specify the OS distribution of the node"
fi
nodeos="$1"
kubectl apply -f falco_loader/$nodeos/falco_loader_pod.yaml
if [ $? -ne 0 ]; then
  exit 1
fi
echo "Waiting for falco loader pod to complete. This might take 5 mins..."
kubectl wait --for=condition=complete pod/falco-loader --timeout=300s
if [ $? -ne 0 ]; then
  exit 1
fi

echo "Starting syscalls tracer..."
kubectl create ns kubesecgen
kubectl apply -f clusterrole.yaml
kubectl apply -f serviceaccount.yaml
kubectl apply -f clusterrolebinding.yaml
kubectl apply -f syscalls-tracer-pod.yaml
kubectl apply -f syscalls-tracer-service.yaml
kubectl get service/syscalls-tracer -n kubesecgen

