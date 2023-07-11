#!/bin/bash
if [ -z "$1" ]
  then
    echo "Please specify the OS distribution of the node"
fi
nodeos="$1"
kubectl apply -f ../install/falco_loader/$nodeos/falco_loader_pod.yaml
if [ $? -ne 0 ]; then
  exit 1
fi
echo "Waiting for falco loader pod to complete. This might take 6 mins..."
kubectl wait --for=condition=Completed pod/falco-loader --timeout=360s
if [ $? -ne 0 ]; then
  exit 1
fi

echo "Starting syscalls tracer..."
kubectl create ns kubesecgen
kubectl apply -f ../install/clusterrole.yaml
kubectl apply -f ../install/serviceaccount.yaml
kubectl apply -f ../install/clusterrolebinding.yaml
kubectl apply -f ../install/syscalls-tracer-pod.yaml
kubectl apply -f ../install/syscalls-tracer-service.yaml
kubectl get service/syscalls-tracer -n kubesecgen

