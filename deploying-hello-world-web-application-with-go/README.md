
> [原文](https://cloud.google.com/kubernetes-engine/docs/tutorials/hello-app) 为 Google Kubernetes Engine(GKE) 的部署应用的教程。本文稍作修改，使其可以部署在本地集群中。

在本示例中，我们将部署一个名为 hello-app 的示例 Web 应用[源代码](https://github.com/GoogleCloudPlatform/kubernetes-engine-samples/tree/master/hello-app)，这是一个用 Go 编写的 Web 服务器，响应请求并显示 Hello, World! 消息。

#### 拉取镜像

```bash
docker pull gcr.io/google-samples/hello-app:1.0
docker pull gcr.io/google-samples/hello-app:2.0
```

执行 `docker image ls` 命令，可见以下两个镜像: 

```
gcr.io/google-samples/hello-app   1.0   bc5c421ecd6c   3 years ago     9.86MB
gcr.io/google-samples/hello-app   2.0   c5607c30fb08   3 years ago     9.86MB
```


#### 创建服务

执行以下命令，创建 hello-web 服务

```bash
kubectl apply -k .
```

访问:

```bash
curl http://localhost
```

显示以下内容：

```
Hello, world!
Version: 1.0.0
Hostname: hello-web-6f844b8699-nmlls
```

#### 更新镜像版本
```bash
kubectl set image deployments/hello-web hello-web=gcr.io/google-samples/hello-app:2.0
```

访问:

```bash
curl http://localhost
```

显示以下内容：

```
Hello, world!
Version: 2.0.0
Hostname: hello-web-7ffff4ffd4-jz52k
```

#### 清理

```bash
kubectl delete -k .
```

