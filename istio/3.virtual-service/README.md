# 虚拟服务

### Demo 介绍

基于 `nginx` 和 `httpd` 镜像，我们部署两个简单的 Web 服务。访问它们时，分别返回 `Welcome to nginx!` 和 `It works!` 两个经典的欢迎页面。

### 部署

```shell
kubectl apply -f nginx-deployment.yaml
kubectl apply -f nginx-service.yaml
kubectl apply -f httpd-deployment.yaml
kubectl apply -f httpd-service.yaml
kubectl apply -f istio-gateway.yaml # istio 网关
kubectl apply -f istio-virtual-service.yaml # istio 虚拟服务
```

### 访问

```shell
curl http://localhost # 返回的 nginx 与 httpd 大致的比例是 4:1
```

### 清理

```shell
kubectl delete -k .
```