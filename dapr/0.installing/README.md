# 在集群中安装 Dapr

Dapr 提醒多种使用方式：

* 独立模式（适用于本地开发）
* Kubernetes 模式（适用于生产环境）
* 特定语言的 SDK

```shell
# install Dapr cli
brew install dapr/tap/dapr-cli

# Deploy the Dapr control plane to Kubernetes cluster
dapr init -k
```


查看状态：
```shell
dapr status -k
```

返回：

```shell
NAME                   NAMESPACE    HEALTHY  STATUS   REPLICAS  VERSION  AGE  CREATED              
dapr-operator          dapr-system  True     Running  1         1.5.1    12m  2022-01-05 22:55.44  
dapr-placement-server  dapr-system  True     Running  1         1.5.1    12m  2022-01-05 22:55.54  
dapr-sidecar-injector  dapr-system  True     Running  1         1.5.1    12m  2022-01-05 22:55.44  
dapr-dashboard         dapr-system  True     Running  1         0.9.0    12m  2022-01-05 22:55.44  
dapr-sentry            dapr-system  True     Running  1         1.5.1    12m  2022-01-05 22:55.44  
```