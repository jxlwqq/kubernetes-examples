# 部署一个简单的 Flask 应用

#### 源代码

一个简单的 Gin 应用程序看起来是这样的：

```go
package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()

	r.GET("/hello", func(c *gin.Context) {
		c.String(200, "Hello, Gin!")
	})

	_ = r.Run(":8080")
}
```

go.mod 文件包含应用所需的依赖。

#### Docker 镜像

应用的 Dockerfile 如下所示：

```dockerfile
# 多阶段构建，减少镜像大小

# 从官方仓库中获取 1.16 的 Go 基础镜像
FROM golang:1.16-alpine AS builder
# 设置工作目录
WORKDIR /go/src/hello-gin
# 复制项目文件
ADD . /go/src/hello-gin
# 下载依赖
RUN go get -d -v ./...
# 构建名为"app"的二进制文件
RUN go build -o app .

# 获取轻型 Linux 发行版，大小仅有 5M 左右
FROM alpine:latest
# 将上一阶段构建好的二进制文件复制到本阶段中
COPY --from=builder /go/src/hello-gin/app .
# 设置监听端口
EXPOSE 8080
# 配置启动命令
CMD ["./app"]
```

构建并提交镜像：

> jxlwqq 是我的 Docker Hub 账号，这里需要换成你自己的账号。

```shell
docker build -f Dockerfile -t jxlwqq/hello-gin:latest . # 构建镜像
docker push jxlwqq/hello-gin:latest # 提交镜像
```

#### 前提条件：部署 nginx ingress

```bash
cd ../ingress-nginx # 切换到 ingress-nginx 目录
kubectl apply -f deploy.yaml
```

#### 部署 hello gin 应用

执行以下命令：

```shell
kubectl apply -f hello-gin-deployment-and-service.yaml
kubectl apply -f ingress.yaml
```

`hello-gin-deployment-and-service.yaml` 文件解读：

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: hello-gin
spec:
  selector:
    matchLabels:
      app: hello-gin # 标签选择器
  template:
    metadata:
      labels:
        app: hello-gin # Pod 标签
    spec:
      containers:
        - name: hello-gin
          image: jxlwqq/hello-gin:latest # 镜像名称:镜像版本
          resources:
            limits:
              memory: "128Mi"
              cpu: "500m"
          ports:
            - containerPort: 8080 # 端口号
---
apiVersion: v1
kind: Service
metadata:
  name: hello-gin-svc
spec:
  selector:
    app: hello-gin # 标签选择器
  ports:
    - port: 8080 # Service 端口号
      targetPort: 8080 # Pod 暴露的端口号
```

`ingress.yaml` 文件解读：

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: hello-gin-ingress
spec:
  rules:
    - http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: hello-gin-svc # service 名称
                port:
                  number: 8080 # 端口号
```

访问验证：

```shell
curl 127.0.0.1/hello # 返回 Hello, Gin!
```

#### 清理
```shell
kubectl delete -k .
```

