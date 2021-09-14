# 使用 Redis 部署 PHP 留言板应用程序

原文：https://v1-20.docs.kubernetes.io/docs/tutorials/stateless-application/guestbook/，基于原文做了一些相关调整。

本教程向您展示如何使用 Kubernetes 和 Docker 构建和部署 一个简单的_(非面向生产)的_多层 web 应用程序。

#### 部署

```shell
kubectl apply -k .
```

#### 清理

```shell
kubectl delete -k .
```

