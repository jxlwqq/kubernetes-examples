# 安装

基于 Knative v0.26.0 版本进行安装调试。

### 前置条件

* 安装 Docker Desktop，并启动内置的 Kubernetes 集群
* 安装 istioctl
* 安装 kn

```
brew install istioctl
brew install kn
```

### 安装 Istio
```shell
istioctl install -y
```

### 安装 Knative Operator
```shell
kubectl apply -f operator.yaml
```

### 安装 Knative Serving
```shell
kubectl apply -f serving.yaml
kubectl apply -f serving-default-domain.yaml
```

### 安装 Knative Eventing
```shell
kubectl apply -f eventing.yaml
```

### 安装第一个应用
```shell
kubectl apply -f hello.yaml
```

查看服务：

```
kn services list
```

返回：

```
AME    URL                                       LATEST        AGE     CONDITIONS   READY   REASON
hello  http://hello.default.127.0.0.1.sslip.io   hello-world   8m27s   3 OK / 3     True
```


访问：

```
curl http://hello.default.127.0.0.1.sslip.io
```

返回：

```
"Hello World!
```
