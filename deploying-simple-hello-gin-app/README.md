# 部署一个简单的 Gin 应用

## 实验前提

* 需要你有 macOS 开发环境，本文以此为例，其他类型的开发环境请自行搭建。
* 需要你对 YAML 这一专门用来写配置文件的语言有所了解。
* 需要你对 Docker 有一些基本的了解。
* 需要你对 Kubernetes 中的 Node、Pod、ReplicaSet、Deployment、Service、Ingress、ConfigMap 等一些核心基础概念有一定的了解。

## YAML 配置文件下载地址：

* YAML 文件：[jxlwqq/kubernetes-examples](https://github.com/jxlwqq/kubernetes-examples/tree/master/deploying-simple-hello-gin-app)。该项目还有其他一些 Kubernetes 的示例。欢迎 Star。

```bash
git clone https://github.com/jxlwqq/kubernetes-examples.git
cd deploying-simple-hello-gin-app
```

## 安装 Docker for Mac

下载地址：https://hub.docker.com/editions/community/docker-ce-desktop-mac

启动并开启 Kubernetes 功能，功能开启过程中，Docker 将会自动拉取 Kubernetes 相关镜像，所以全程需要科学上网。

为啥不使用 minikube？minikube + virtualbox + kubectl 安装起来太繁琐了，而且即使科学上网了你也不一定能搞定。当然阿里云提供了一篇[安装教程](https://yq.aliyun.com/articles/221687)可以参考。

## 本地端口准备

请确保本地 localhost 的 80 端口没有被占用，已在使用的请在实验期间暂时关闭占用 80 端口的服务。

## 切换集群

如果你本地有多个 Kubernetes 的集群配置，请先切换至名为 docker-desktop 的集群：

````bash
kubectl config use-context docker-desktop
````

## 源代码

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

go.mod 文件包含应用所需的依赖：

```
module hello-gin

go 1.17

require github.com/gin-gonic/gin v1.7.4

require (
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/go-playground/validator/v10 v10.9.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/json-iterator/go v1.1.11 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/mattn/go-isatty v0.0.13 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/ugorji/go/codec v1.2.6 // indirect
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97 // indirect
	golang.org/x/sys v0.0.0-20210806184541-e5e7981a1069 // indirect
	golang.org/x/text v0.3.6 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
```

## Docker 镜像

应用的 Dockerfile 如下所示：

```dockerfile
# 多阶段构建：提升构建速度，减少镜像大小

# 从官方仓库中获取 1.17 的 Go 基础镜像
FROM golang:1.17-alpine AS builder

# 设置工作目录
WORKDIR /workspace

# 安装项目依赖
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

# 复制项目文件
COPY . .

# 构建名为"app"的二进制文件
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o app main.go

# 获取 Distroless 镜像，只有 650 kB 的大小，是常用的 alpine:latest 的 1/4
FROM gcr.io/distroless/static:nonroot
# 设置工作目录
WORKDIR /
# 将上一阶段构建好的二进制文件复制到本阶段中
COPY --from=builder /workspace/app .
# 设置监听端口
EXPOSE 8080
# 配置启动命令
ENTRYPOINT ["/app"]
```

构建并提交镜像：

> jxlwqq 是我的 Docker Hub 账号，这里需要换成你自己的账号。如果没有账号，需要先注册：https://hub.docker.com/signup

> 这一步如果想跳过的话，暂时可以直接拉取我制作好的镜像：`docker pull jxlwqq/hello-gin:latest`

```shell
docker build -f Dockerfile -t jxlwqq/hello-gin:latest . # 构建镜像
docker login # 登录
docker push jxlwqq/hello-gin:latest # 提交镜像
```

## 前提条件：部署 nginx ingress

为了让 Ingress 资源工作，集群必须有一个正在运行的 Ingress 控制器。 Kubernetes 官方目前支持和维护 GCE 和 nginx 控制器。

这里我们选择 Ingress-nginx 控制器：

```bash
kubectl apply -f ../ingress-nginx/deploy.yaml
```


> 注： deploy.yaml 文件内容来源自：https://github.com/kubernetes/ingress-nginx/blob/main/deploy/static/provider/cloud/deploy.yaml

> 详细操作说明见：https://github.com/kubernetes/ingress-nginx/blob/main/docs/deploy/index.md


## 部署 hello gin 应用

执行以下命令：

```shell
kubectl apply -f hello-gin-deployment-and-service.yaml
kubectl apply -f ingress.yaml
```

返回：
```shell
service/hello-gin-svc created
deployment.apps/hello-gin created
ingress.networking.k8s.io/hello-gin-ingress created
```

`hello-gin-deployment-and-service.yaml` 文件解读：

```yaml
apiVersion: apps/v1 # api 版本
kind: Deployment # 资源对象类型
metadata: # Deployment 元数据
  name: hello-gin # 对象名称
spec: # 对象规约
  selector: # 选择器，作用：选择带有下列标签的Pod
    matchLabels: # 标签匹配
      app: hello-gin # 标签KeyValue
  template: # Pod 模版
    metadata: # Pod元数据
      labels: # Pod 标签
        app: hello-gin # Pod 标签，与上述的 Deployment.selector中的标签对应
    spec: # Pod 对象规约
      containers: # 容器
        - name: hello-gin # 容器名称
          image: jxlwqq/hello-gin:latest # 镜像名称:镜像版本
          resources: # 资源限制
            limits: # 简单理解为max资源值
              memory: "128Mi"
              cpu: "500m"
            requests: # 简单理解为min资源值
              memory: "128Mi"
              cpu: "500m"
          ports: # 端口
            - containerPort: 8080 # 端口号
---
apiVersion: v1 # api 版本
kind: Service # 对象类型
metadata: # 元数据
  name: hello-gin-svc # 对象名称
spec: # 规约
  selector: # 选择器
    app: hello-gin # 标签选择器，与 Pod 的标签对应
  ports:
    - port: 8080 # Service 端口号
      targetPort: 8080 # Pod 暴露的端口号
```

`ingress.yaml` 文件解读：

```yaml
apiVersion: networking.k8s.io/v1 # api 版本
kind: Ingress # 对象类型
metadata: # 元数据
  name: hello-gin-ingress # 对应名称
spec: # 规约
  rules: # 规则
    - http:
        paths: # 路径
          - path: /
            pathType: Prefix
            backend: # 后端服务
              service: # 服务
                name: hello-gin-svc # service 名称
                port:
                  number: 8080 # 端口号
  ingressClassName: nginx
```

访问验证：

```shell
curl 127.0.0.1/hello # 返回 Hello, Gin!
```

## 清理
```shell
kubectl delete -k .
```

