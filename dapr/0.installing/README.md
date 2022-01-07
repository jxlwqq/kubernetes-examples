# 在集群中安装 Dapr

Dapr 提供多种使用方式：

* 独立模式（适用于本地开发）
* Kubernetes 模式（适用于生产环境）
* 特定语言的 SDK

我们直接以 Kubernetes 模式开始。


### 在 Kubernetes 集群上设置 Dapr

首先，下载并安装最新版本的 Docker Desktop，并启动内置的 Kubernetes 集群。

然后执行以下两个命令，完成 Dapr 的基础设置：

```shell
# 安装 Dapr 客户端
brew install dapr/tap/dapr-cli

# 在 Kubernetes 集群中安装 Dapr 控制平面
dapr init -k
```

执行完成后，查看状态：

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


### 创建和配置状态存储

Dapr 支持多种不同的状态存储（如 Redis、CosmosDB、DynamoDB、Cassandra 等）来持久化和检索状态。本演示将使用 Redis。

首先，我们使用 helm 创建一个 高可用的 Redis 集群：

```shell
# 如果本地没有 helm，则安装
brew install helm
# 创建 Redis 集群
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update
helm install redis bitnami/redis
```

查看 Pod：

```shell
kubectl get pods
```

返回：
```shell
NAME               READY   STATUS    RESTARTS   AGE
redis-master-0     1/1     Running   0          83s
redis-replicas-0   1/1     Running   0          83s
redis-replicas-1   1/1     Running   0          46s
redis-replicas-2   1/1     Running   0          24s
```


应用 redis.yaml 文件并观察您的状态存储是否已成功配置：

```shell
kubectl apply -f redis.yaml
```

redis.yaml 文件内容如下：

```yaml
apiVersion: dapr.io/v1alpha1
kind: Component
metadata:
  name: statestore
spec:
  type: state.redis
  version: v1
  metadata:
    # These settings will work out of the box if you use `helm install
    # bitnami/redis`.  If you have your own setup, replace
    # `redis-master:6379` with your own Redis master address, and the
    # Redis password with your own Secret's name. For more information,
    # see https://docs.dapr.io/operations/components/component-secrets .
    - name: redisHost
      value: redis-master:6379
    - name: redisPassword
      secretKeyRef:
        name: redis
        key: redis-password
auth:
  secretStore: kubernetes
```

