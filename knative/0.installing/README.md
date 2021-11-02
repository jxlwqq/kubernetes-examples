# 安装

基于 Knative v0.26.0 版本进行安装调试。

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