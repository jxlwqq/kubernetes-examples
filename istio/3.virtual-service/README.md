# 虚拟服务

### Demo 介绍

基于 jxlwqq/http-echo 镜像，我们部署两个简单的 Web 服务。访问它们时，分别返回 `v1` 和 `v2`。

### 部署

```shell
kubectl apply -f echo-v1-deployment.yaml
kubectl apply -f echo-v1-service.yaml
kubectl apply -f echo-v2-deployment.yaml
kubectl apply -f echo-v2-service.yaml
kubectl apply -f istio-gateway.yaml # istio 网关
kubectl apply -f istio-virtual-service.yaml # istio 虚拟服务
```

### 访问

```shell
curl http://localhost # 返回的 v1 与 v2 大致的比例是 9:1
```

### 清理

```shell
kubectl delete -k .
```