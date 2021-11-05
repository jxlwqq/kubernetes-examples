```shell
export FORTIO_POD=$(kubectl get pods -l app=fortio -o 'jsonpath={.items[0].metadata.name}')
kubectl exec "$FORTIO_POD" -c fortio -- /usr/bin/fortio curl -quiet http://nginx-svc:80
kubectl exec "$FORTIO_POD" -c fortio -- /usr/bin/fortio load -c 3 -qps 0 -n 30 -loglevel Warning http://nginx-svc:80
kubectl exec "$FORTIO_POD" -c istio-proxy -- pilot-agent request GET stats | grep nginx| grep pending
```