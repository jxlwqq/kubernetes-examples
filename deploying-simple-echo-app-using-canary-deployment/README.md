# 使用 Nginx Ingress 实现金丝雀发布

> [原文](https://cloud.tencent.com/document/product/457/48907) 为腾讯云容器服务（Tencent Kubernetes Engine，TKE）部署应用的教程。本文对此基础上进行了较大的变更，使其可以部署在本地集群中。

目前 Nginx Ingress 支持基于 Header、Cookie 和服务权重3种流量切分的策略。 通过给 Ingress 资源指定 Nginx Ingress 所支持的 annotation 可实现金丝雀发布。需给服务创建 2 个 Ingress，其中 1 个常规 Ingress，另 1 个为带 `nginx.ingress.kubernetes.io/canary: "true"` 固定的 annotation 的 Ingress，称为 Canary Ingress。Canary Ingress 一般代表新版本的服务，结合另外针对流量切分策略的 annotation 一起配置即可实现多种场景的金丝雀发布。

本文以 canary-by-header 策略作为示例。实现以下场景：

> 来自北京和上海的用户体验新版本的服务，而来自其他地区的用户继续使用老版本的服务。

其他流量切分策略及其匹配的使用场景，请参考 [腾讯云的原文](https://cloud.tencent.com/document/product/457/48907) 或者 [Nginx Ingress 官方文档](https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/annotations/#canary)。

#### 前提条件：部署 nginx ingress

```bash
cd ../ingress-nginx # 切换到 ingress-nginx 目录
kubectl apply -f deploy.yaml
```

#### 提前拉取 Docker 镜像

```bash
docker pull jxlwqq/http-echo
```

[jxlwqq/http-echo](https://github.com/jxlwqq/http-echo) 是一个非常小巧的，运行在内存中的 web 服务器，它呈现一个 HTML 页面，其中包含提供给它的参数的内容。这对于测试或演示“hello world”程序特别有用。

#### 部署 v1 版本的 echo 服务

执行以下命令：

```shell
kubectl apply -f echo-v1-deployment-and-service.yaml
```

`echo-v1-deployment-and-service.yaml`的解读：

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: echo-v1
spec:
  selector:
    matchLabels:
      app: echo
      version: v1
  template:
    metadata:
      labels:
        app: echo
        version: v1
    spec:
      containers:
        - name: echo
          image: jxlwqq/http-echo
          args:
            - "--text=echo-v1" # 响应请求，返回"echo-v1"
          ports:
            - containerPort: 8080 # 容器端口号

---
apiVersion: v1
kind: Service
metadata:
  name: echo-v1-svc
spec:
  selector:
    app: echo
    version: v1
  ports:
    - port: 80
      targetPort: 8080
```

#### 部署 v2 版本的 echo 服务

执行以下命令：

```shell
kubectl apply -f echo-v2-deployment-and-service.yaml
```

`echo-v2-deployment-and-service.yaml`的解读：

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: echo-v2
spec:
  selector:
    matchLabels:
      app: echo
      version: v2
  template:
    metadata:
      labels:
        app: echo
        version: v2
    spec:
      containers:
        - name: echo
          image: jxlwqq/http-echo
          args:
            - "--text=echo-v2" # 响应请求，返回"echo-v2"
          ports:
            - containerPort: 8080 # 容器端口号

---
apiVersion: v1
kind: Service
metadata:
  name: echo-v2-svc
spec:
  selector:
    app: echo
    version: v2
  ports:
    - port: 80
      targetPort: 8080
```

#### 创建 Ingress，对外暴露服务，指向 v1 版本的 echo 服务

执行以下命令：

```shell
kubectl apply -f ingress-primary.yaml
```

`ingress-primary.yaml`文件的解读：

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: echo-ing
  annotations:
    kubernetes.io/ingress.class: nginx
spec:
  rules:
    - host: canary.example.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: echo-v1-svc
                port:
                  number: 80
  ingressClassName: nginx
```

执行以下命令，进行访问验证：

```shell
curl -H "Host: canary.example.com" 127.0.0.1   # 返回：echo-v1
```

#### 基于 Header 的流量切分，创建 Canary Ingress，指定 v2 版本的 echo 服务

执行以下命令：

```shell
kubectl apply -f ingress-canary-by-header.yaml
```

`ingress-canary-by-header.yaml`的解读：

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: echo-canary-ing
  annotations:
    nginx.ingress.kubernetes.io/canary: "true" # 支持金丝雀
    nginx.ingress.kubernetes.io/canary-by-header: "Region" # 基于请求头中的"Region"字段切分流量
    nginx.ingress.kubernetes.io/canary-by-header-pattern: "shanghai|beijing" #当请求头中的"Region"的值匹配"shanghai"或者"beijing"的时候切分流量
spec:
  rules:
    - host: canary.example.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: echo-v2-svc
                port:
                  number: 80
  ingressClassName: nginx
```

执行以下命令，进行访问验证：

```shell
curl -H "Host: canary.example.com" 127.0.0.1 # 返回：echo-v1
curl -H "Host: canary.example.com" -H "Region: shanghai" 127.0.0.1   # 返回：echo-v2
curl -H "Host: canary.example.com" -H "Region: beijing" 127.0.0.1    # 返回：echo-v2
curl -H "Host: canary.example.com" -H "Region: shenzhen" 127.0.0.1   # 返回：echo-v1
curl -H "Host: canary.example.com" -H "Region: guangzhou" 127.0.0.1  # 返回：echo-v1
```


#### 清理

```shell
kubectl delete -k .
```
