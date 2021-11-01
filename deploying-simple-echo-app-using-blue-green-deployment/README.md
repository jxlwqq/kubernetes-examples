# 基于 Service 的 selector 实现蓝绿发布

> [原文](https://cloud.tencent.com/document/product/457/48877) 为腾讯云容器服务（Tencent Kubernetes Engine，TKE）部署应用的教程。本文对此基础上进行了较大的变更，使其可以部署在本地集群中。

以 Deployment 为例，集群中已部署两个不同版本的 Deployment，其 Pod 拥有共同的 label。但有一个 label 值不同，用于区分不同的版本。Service 使用 selector 选中了其中一个版本的 Deployment 的 Pod，此时通过修改 Service 的 selector 中决定服务版本的 label 的值来改变 Service 后端对应的 Deployment，即可实现让服务从一个版本直接切换到另一个版本。

#### 前提条件：部署 nginx ingress

```bash
kubectl apply -f ../ingress-nginx/deploy.yaml
```

#### 部署 v1 版本的 echo Deployment

执行以下命令：

```shell
kubectl apply -f echo-v1-deployment.yaml
```

`echo-v1-deployment.yaml`的解读：

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: echo-v1
spec:
  replicas: 1
  selector:
    matchLabels:
      app: echo
      version: v1
  template:
    metadata:
      labels:
        app: echo # 标签1
        version: v1 # 标签2
    spec:
      containers:
        - name: echo
          image: jxlwqq/http-echo
          args:
            - "--text=echo-v1" # 响应请求，返回"echo-v1"
          ports:
            - name: http
              protocol: TCP
              containerPort: 8080 # 容器端口号
```

#### 部署 v2 版本的 echo Deployment

执行以下命令：

```shell
kubectl apply -f echo-v2-deployment.yaml
```

`echo-v2-deployment.yaml`的解读：

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: echo-v2
spec:
  replicas: 1
  selector:
    matchLabels:
      app: echo
      version: v2
  template:
    metadata:
      labels:
        app: echo # 标签1
        version: v2 # 标签2
    spec:
      containers:
        - name: echo
          image: jxlwqq/http-echo
          args:
            - "--text=echo-v2" # 响应请求，返回"echo-v2"
          ports:
            - name: http
              protocol: TCP
              containerPort: 8080
```




#### 创建 Service

```shell
kubectl apply -f service.yaml
```

`service.yaml`的解读：

```yaml
apiVersion: v1
kind: Service
metadata:
  name: echo-svc
spec:
  selector: # 选择 v1 版本的 Pod
    app: echo 
    version: v1
  ports:
    - port: 80
      targetPort: 8080
```

#### 创建 Ingress

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: echo-ing
spec:
  rules:
    - http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: echo-svc
                port:
                  number: 8080
  ingressClassName: nginx
```

访问验证：

```shell
curl 127.0.0.1 # 返回 echo-v1
```

#### 修改 Service 的 selector，使其选中 v2 版本的服务

```shell
kubectl patch service echo-svc -p '{"spec":{"selector":{"app": "echo", "version": "v2"}}}'
```

访问验证：

```shell
curl 127.0.0.1 # 返回 echo-v2
```

#### 清理

```shell
kubectl delete -k .
```