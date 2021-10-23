# 在 Docker Desktop 集群中安装 Metrics Server

### 部署

部署完组件后，需要 patch 下一个参数，才能正常运行：

```shell
# 部署组件
kubectl apply -f metrics-server.yaml
# 增加一个 kubelet-insecure-tls 参数
kubectl patch deployments.apps \
  metrics-server \
  --namespace kube-system \
  --type='json' \
  -p='[{"op": "replace", "path": "/spec/template/spec/containers/0/args", "value": [
  "--cert-dir=/tmp",
  "--secure-port=443",
  "--kubelet-preferred-address-types=InternalIP,ExternalIP,Hostname",
  "--kubelet-use-node-status-port",
  "--metric-resolution=15s",
  "--kubelet-insecure-tls"
]}]'
```

### 试验

```shell
# 部署一个 nginx deployment
kubectl apply -f nginx.yaml
# 自动扩缩
kubectl autoscale deployment nginx --min=1 --max=10 --cpu-percent=10
# 暴露服务
kubectl expose deployment nginx --port=80 --type=LoadBalancer
# 查看
kubectl top nodes
kubectl top pods
```


```shell
# 获取 LoadBalancer 的 IP 地址
kubectl get services nginx -o jsonpath='{.status.loadBalancer.ingress[0].ip}'
# 增大负载
while true; do wget -q -O- http://localhost; done
# 打开另外一个窗口观察
kubectl get pods --watch
```


### 清理

```shell
kubectl delete -f metrics-server.yaml
kubectl delete -f nginx.yaml
kubectl delete svc nginx
kubectl delete hpa nginx
```