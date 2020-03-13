# 使用 MySQL 和 Redis 部署 Laravel 应用

#### 实验前提

* 需要你有 macOS 开发环境，本文以此为例，其他类型的开发环境请自行搭建。
* 需要你对 YAML 这一专门用来写配置文件的语言有所了解。
* 需要你对 Docker 有一些基本的了解。
* 需要你对 Kubernetes 中的 Node、Pod、ReplicaSet、Deployment、Service、Ingress、ConfigMap 等一些核心基础概念有一定的了解。

#### 安装 Docker for Mac

下载地址：https://hub.docker.com/editions/community/docker-ce-desktop-mac

启动并开启 Kubernetes 功能，功能开启过程中，Docker 将会自动拉取 Kubernetes 相关镜像，所以全程需要科学上网。

#### 本地端口准备

请确保本地 localhost 的 80 端口没有被占用，已在使用的请在实验期间暂时关闭占用 80 端口的服务。

#### 切换集群

如果你本地有多个 Kubernetes 的集群配置，请先切换至名为 docker-desktop 的集群：

````bash
kubectl config use-content docker-desktop
````

#### 创建 MySQL 服务

为了提高 Pod 的启动速度，我们首先准备好 MySQL 的镜像：

```bash
docker pull mysql:5.7
```
部署 MySQL 服务：

```bash
kubectl apply -f mysql-deployment-and-service.yaml
```

进入 Pod：

```bash
kubectl exec -it mysql-deployment-79cdbc594-rmhjk mysql # pod 名称改成你自己的
```

创建一个名为 laravel 的数据库：

```bash
create database laravel;
```
数据库 laravel 创建完成后，即可 exit。

#### 创建 Redis 服务

为了提高 Pod 的启动速度，我们首先准备好 Redis 的镜像：

```bash
docker pull redis
```

部署 Redis 服务：

```bash
kubectl apply -f redis-deployment-and-service.yaml
```

#### 创建 Laravel 服务

为了提高 Pod 的启动速度，我们首先准备好 Laravel-demo 的镜像：

镜像地址：https://hub.docker.com/repository/docker/jxlwqq/laravel-7-kubernetes-demo
源码地址：https://github.com/jxlwqq/laravel-7-kubernetes-demo

基于官方 Laravel v7 版本，做了以下修改：[compare](https://github.com/jxlwqq/laravel-7-kubernetes-demo/compare/e47e5cc7029408ed80e0cd0298d944f5b49b9cdd...master)

* 增加了 Docker 镜像构建相关的文件：Dockerfile 和 .dockerignore
* 增加了一个 config/apache2/sites-available/laravel.conf 
* 增加了一个 crontab 文件，执行 php artisan schedule:run
* 增加了一个 docker/entrypoint.sh 文件：Docker 容器启动时执行的命令

```bash
docker pull jxlwqq/laravel-7-kubernetes-demo
```

部署 Laravel 服务：

```bash
kubectl apply -f configmap.yaml # 将 env 环境变量配置在了 ConfigMap 对象里
kubectl apply -f laravel-deployment-and-service.yaml # 部署 Deployment 和 Service
 kubectl apply -f ingress.yaml # 部署 ingress 路由规则
```

#### 创建 Ingress-nginx 控制器
```bash
cd ingress-nginx
kubectl apply -f ingress-nginx-deployment-and-other-resources-mandatory.yaml
kubectl apply -f ingress-nginx-service.yaml
```

#### 访问

```bash
curl http://localhost
```