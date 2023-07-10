kubectl create ns kubesecgen
kubectl apply -f clusterrole.yaml
kubectl apply -f serviceaccount.yaml
kubectl apply -f clusterrolebinding.yaml
kubectl apply -f syscalls-tracer-pod.yaml
kubectl apply -f syscalls-tracer-service.yaml
kubectl get service/syscalls-tracer -n kubesecgen

