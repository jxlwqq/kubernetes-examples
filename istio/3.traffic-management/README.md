# 流量管理

### Demo 介绍

基于 hashicorp/http-echo 镜像，我们部署两个简单的 Web 服务。访问它们时，分别返回 `echo-v1` 和 `echo-v2`。

### 部署

```shell
kubectl apply -f echo-v1.yaml # 包含 deployment 和 service
kubectl apply -f echo-v2.yaml # 包含 deployment 和 service
kubectl apply -f gateway.yaml # istio 网关
kubectl apply -f virtual-service.yaml # istio 虚拟服务
```

### 访问

```shell
curl http://localhost # 返回的 echo-v1 与 echo-v2 大致的比例是 4:1
```

### 清理

```shell
kubectl delete -f echo-v1.yaml
kubectl delete -f echo-v2.yaml
kubectl delete -f gateway.yaml
kubectl delete -f virtual-service.yaml
```