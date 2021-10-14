# 部署一个简单的加减乘除计算器微服务

本示例使用 grpc 和 Gin 来构建一个简单的加减乘除计算服务。项目架构很简单，分为两部分：

* 一个是对外暴露的 Web 服务：api-server
* 一个是内部调用的微服务：calculator

### 前置依赖

protoc: 

```shell
brew install protoc
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

ingress-nginx 控制器：

```shell
kubectl apply -k ../ingress-nginx
```

### 构建镜像

```shell
make docker-build
```

### 部署

```shell
make kube-deploy
```