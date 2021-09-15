# KBP Journal App

《Kubernetes Best Practices》第一章「搭建一个基本服务」的示例应用程序。随书仓库 https://github.com/brendandburns/kbp-sample 代码不全。

本示例中的应用程序包含：
* 一个用于存储数据的 Redis 后端服务 `redis/`
* 一个简单的日志系统 `frontend/`
* 一个 Nginx 静态文件服务器 `fileserver/`

#### 部署有状态 Redis 后端服务

使用 StatefulSet 资源来部署 Redis 集群，使用卷说明来编写可复制的模版，为多副本中的每个 Pod 分配自己独有的 PV。集群中的 Leader 和 Follower 使用存储在 ConfigMap 中启动脚本区分角色。

```shell
kubectl apply -f redis/
```

#### 部署日志服务

日志系统的前端采用 TS 实现的一个 Node.js 应用程序，该应用程序使用暴露在 8080 端口上的 HTTP服务来处理请求，并采用 Redis 作为后端来 CURD 当前日志条目。

```shell
kubectl apply -f frontend/
```

#### 部署静态文件服务器

使用 Deployment 来声明多副本的 Nginx 服务器。

```shell
kubectl apply -f fileserver/
```

#### 部署 Ingress

```shell
kubectl apply -f ingress.yaml
kubectl apply -f ../ingress-nginx/deploy.yaml
```

#### 清理

```shell
kubectl delete -k .
```

