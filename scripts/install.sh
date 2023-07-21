#!/bin/bash
if [ -z "$1" ]
  then
    echo "Please specify the OS distribution of the node"
fi
nodeos="$1"
case nodeos in
  centos-stream8)
    echo -n "OS Distribution Selected: centos-stream8"
    ;;

  *)
    printf "Invalid OS Distrbution. \nSupported distributions:\n    1.centos-stream8"
    exit 1
    ;;
esac
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
kubectl apply -f ../install
kubectl get service/syscalls-tracer -n kubesecgen

