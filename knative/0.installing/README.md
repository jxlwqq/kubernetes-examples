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
kubectl apply -f http-echo.yaml
```

查看服务：

```
kn services list
```

返回：

```
NAME        URL                                           LATEST         AGE   CONDITIONS   READY   REASON
http-echo   http://http-echo.default.127.0.0.1.sslip.io   http-echo-v1   21s   3 OK / 3     True  
```


访问：

```
curl http://http-echo.default.127.0.0.1.sslip.io
```

返回：

```
v1
```
