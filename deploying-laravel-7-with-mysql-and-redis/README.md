## 实验前提

* 需要你有 macOS 开发环境，本文以此为例，其他类型的开发环境请自行搭建。
* 需要你对 YAML 这一专门用来写配置文件的语言有所了解。
* 需要你对 Docker 有一些基本的了解。
* 需要你对 Kubernetes 中的 Node、Pod、ReplicaSet、Deployment、Service、Ingress、ConfigMap 等一些核心基础概念有一定的了解。

## YAML 配置文件下载地址：

* YAML 文件：[jxlwqq/kubernetes-examples](https://github.com/jxlwqq/kubernetes-examples/tree/master/deploying-laravel-7-with-mysql-and-redis)。该项目还有其他一些 Kubernetes 的示例。欢迎 Star。

```bash
git clone https://github.com/jxlwqq/kubernetes-examples.git
cd deploying-laravel-7-with-mysql-and-redis
```

## 安装 Docker for Mac

下载地址：https://hub.docker.com/editions/community/docker-ce-desktop-mac

启动并开启 Kubernetes 功能，功能开启过程中，Docker 将会自动拉取 Kubernetes 相关镜像，所以全程需要科学上网。

为啥不使用 minikube？minikube + virtualbox + kubectl 安装起来太繁琐了，而且即使科学上网了你也不一定能搞定。当然阿里云提供了一篇[安装教程](https://yq.aliyun.com/articles/221687)可以参考。

## 本地端口准备

请确保本地 localhost 的 80 端口没有被占用，已在使用的请在实验期间暂时关闭占用 80 端口的服务。

## 切换集群

如果你本地有多个 Kubernetes 的集群配置，请先切换至名为 docker-desktop 的集群：

````bash
kubectl config use-context docker-desktop
````

## 创建 MySQL 服务

为了提高 Pod 的启动速度，我们首先准备好 MySQL 的镜像：

```bash
docker pull mysql:5.7
```
部署 MySQL 服务：

```bash
kubectl apply -f mysql-deployment-and-service.yaml
```

注意：deployment 在生产场景中对 MySQL 这种有状态的服务并不适合。

yaml 文件解读：

```yaml
kind: Deployment # 对象类型
apiVersion: apps/v1 # api 版本
metadata: # 对象元数据
  name: mysql-deployment # 对象名称，加后缀 -deployment 主要是为了新手看的明白，也可以直接改为 name: mysql
  labels: # 对象标签
    app: mysql # 标签
spec: # 对象规约，注意这里的对象指的是 Deployment 对象
  selector: # 选择器，不太恰当的类比就相当于 CSS 选择器
    matchLabels: # 匹配标签
      app: mysql # Pod 标签
  strategy: # 部署策略
    type: Recreate # Recreate 重新创建  RollingUpdate 滚动升级
  template: # 模版
    metadata: # 元数据
      labels: # 标签
        app: mysql # Pod 的标签，这里的值与上面的 selector matchLabels 的标签对应
    spec: # 对象规约，注意这里的对象指的是 Pod 对象
      containers: # 容器
        - name: mysql # 容器名称
          image: mysql:5.7 # 容器镜像
          env: # 环境变量
            - name: MYSQL_ALLOW_EMPTY_PASSWORD
              value: 'true'
          ports: # 端口
            - containerPort: 3306 # 容器端口
          volumeMounts: 
            - mountPath: /var/lib/mysql
              name: mysql-storage
      volumes:
        - name: mysql-storage
          persistentVolumeClaim: # PersistentVolume（PV）是一块集群里由管理员手动提供，或 kubernetes 通过 StorageClass 动态创建的存储。 PersistentVolumeClaim（PVC）是一个满足对 PV 存储需要的请求。
            claimName: mysql-persistentvolumeclaim
---
kind: PersistentVolumeClaim # 对象类型
apiVersion: v1 # api 版本
metadata: # 元数据
  name: mysql-persistentvolumeclaim # 对象名称
  labels: # 标签
    app: mysql
spec: # 对象规约
  accessModes: # 访问方式
    - ReadWriteOnce
  resources: # 资源
    requests: # 请求配置
      storage: 1Gi # 大小，这里主要做了实验，我们就设置小点的
---
kind: Service # 对象类型
apiVersion: v1 # api 版本
metadata: # 元数据
  name: mysql-service # 对象名称
  labels: # 标签
    app: mysql
spec: # 对象规约
  selector: # 选择器
    app: mysql # Pod 的标签
  ports: # 端口
    - port: 3306 # 端口号
      targetPort: 3306 # 与 Pod  containerPort 端口号一致
 ```

进入 Pod：

```bash
kubectl get pods # 获取 pods 列表
kubectl exec -it mysql-deployment-79cdbc594-rmhjk mysql # pod 名称改成你自己
```

创建一个名为 laravel 的数据库：

```bash
create database laravel;
```
数据库 laravel 创建完成后，即可 exit。

## 创建 Redis 服务

为了提高 Pod 的启动速度，我们首先准备好 Redis 的镜像：

```bash
docker pull redis
```

部署 Redis 服务：

```bash
kubectl apply -f redis-deployment-and-service.yaml
```

yaml 文件解读：

```yaml
kind: Deployment # 对象类型
apiVersion: apps/v1 # api 版本
metadata: # 元数据
  name: redis-deployment # 对象名称
  labels: # 对象标签
    app: redis
spec: # 对象规约，这里的对象指的是 Deployment 对象
  selector: #  选择器
    matchLabels: # 匹配标签
      app: redis 
  template: # 模版
    metadata: # 元数据
      labels: # 标签
        app: redis
    spec: # 对象规约，这里的对象指的是 Pod 对象
      containers: # 容器 
        - name: redis # 容器名称
          image: redis # 容器镜像
          ports: # 端口
            - containerPort: 6379 # Redis 的服务端口
---
kind: Service # 对象类型
apiVersion: v1 # api 版本
metadata: # 元数据
  name: redis-service # 对象名称
spec: # 对象规约
  selector: # 选择器
    app: redis # 标签，选择 标签包含 app: redis 的 一组 Pod
  ports: # 端口
    - port: 6379 # 暴露的端口，如果 port 和 targetPort 一致，targetPort 可以不写
 ```

## 创建 Laravel 服务

为了提高 Pod 的启动速度，我们首先准备好 Laravel-demo 的镜像：

```bash
docker pull jxlwqq/laravel-7-kubernetes-demo
```

镜像地址：https://hub.docker.com/repository/docker/jxlwqq/laravel-7-kubernetes-demo

源码地址：https://github.com/jxlwqq/laravel-7-kubernetes-demo

源码基于官方 Laravel v7 版本，做了以下修改：[点击查看compare](https://github.com/jxlwqq/laravel-7-kubernetes-demo/compare/e47e5cc7029408ed80e0cd0298d944f5b49b9cdd...master)

* 增加了 Docker 镜像构建相关的文件：Dockerfile 和 .dockerignore
* 增加了一个 config/apache2/sites-available/laravel.conf 
* 增加了一个 crontab 文件，执行 php artisan schedule:run
* 增加了一个 docker/entrypoint.sh 文件：Docker 容器启动时执行的命令

我这里把核心的 Dockerfile 和 entrypoint.sh 的代码贴过来：

Dockerfile:

```bash
# First Stage
FROM node:alpine as frontend
COPY package.json package-lock.json /app/
RUN cd /app \
    && npm install
COPY webpack.mix.js /app/
COPY resources/js/ /app/resources/js/
COPY resources/sass/ /app/resources/sass/
RUN cd /app \
      && npm run production

# Second Stage
FROM composer as composer
COPY database/ /app/database/
COPY composer.json composer.lock /app/
RUN cd /app \
      && composer install \
           --optimize-autoloader \
           --ignore-platform-reqs \
           --prefer-dist \
           --no-interaction \
           --no-plugins \
           --no-scripts \
           --no-dev

# Third Stage
FROM php:7.4-apache-buster
RUN apt-get update \
    && apt-get install -y vim cron libmagickwand-dev imagemagick
RUN docker-php-ext-install intl pdo_mysql bcmath \
    && pecl install redis imagick \
    && docker-php-ext-enable opcache redis imagick \
    && cp /usr/local/etc/php/php.ini-production /usr/local/etc/php/php.ini
RUN apt-get clean \
    && apt-get autoclean \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

ARG LARAVEL_PATH=/var/www/laravel
WORKDIR ${LARAVEL_PATH}

COPY . ${LARAVEL_PATH}
COPY --from=composer /app/vendor/ ${LARAVEL_PATH}/vendor/
COPY --from=frontend /app/public/js/ ${LARAVEL_PATH}/public/js/
COPY --from=frontend /app/public/css/ ${LARAVEL_PATH}/public/css/
COPY --from=frontend /app/mix-manifest.json ${LARAVEL_PATH}/mix-manifest.json

RUN cd ${LARAVEL_PATH} \
      && php artisan package:discover \
      && chown www-data:www-data bootstrap/cache \
      && chown -R www-data:www-data storage/

RUN rm /etc/apache2/sites-enabled/*
COPY config/apache2 /etc/apache2/
RUN a2enmod rewrite headers \
    && a2ensite laravel

COPY docker/entrypoint.sh /usr/local/bin/entrypoint
RUN chmod +x /usr/local/bin/entrypoint
ENTRYPOINT ["/usr/local/bin/entrypoint"]
```
entrypoint.sh:

```bash
#!/usr/bin/env bash

set -e

cd /var/www/laravel
rm -f public/storage

echo 'migrate'
php artisan migrate --force

echo 'publish'
# php artisan vendor:publish --tag=laravel-pagination

echo 'cache'
php artisan config:cache
php artisan view:cache
php artisan route:cache
php artisan event:cache

echo "cron"
mkdir -p /var/spool/cron/crontabs/
cp crontab /var/spool/cron/crontabs/root
chmod 0644 /var/spool/cron/crontabs/root
crontab /var/spool/cron/crontabs/root
cron -f &

echo "queue"
php artisan queue:work --queue={default} --verbose --tries=3 --timeout=90 &

echo 'http'
exec apache2-foreground
```

部署 Laravel 服务：

```bash
kubectl apply -f configmap.yaml # 将 env 环境变量配置在了 ConfigMap 对象里
kubectl apply -f laravel-deployment-and-service.yaml # 部署 Deployment 和 Service
kubectl apply -f ingress.yaml # 部署 ingress 路由规则
```

configmap.yaml 文件解读：

```yaml
kind: ConfigMap # 对象类型：配置
apiVersion: v1 # api 版本
metadata: # 原数据
  name: laravel-env # 对象名称
data: # 变量数据，跟 laravel env 文件类似（不是太恰当的比方）
  APP_KEY: base64:zC8wVldUZfZJaGaZ7+CPh+5FzaXYmShm7G/Qh6GdRl8=
  APP_ENV: production
  DB_DATABASE: laravel
  DB_USERNAME: root
 ```

laravel-deployment-and-service.yaml 文件解读：

```yaml
kind: Deployment # 对象类型
apiVersion: apps/v1 # api 版本
metadata: # 元数据
  name: laravel-deployment # 对象名称
  labels: # 对象标签
    app: laravel
spec: # 对象规约
  selector: # 选择器
    matchLabels: # 标签匹配
      app: laravel
  replicas: 1 # 副本数量
  strategy: # 部署策略
    type: RollingUpdate # 滚动升级
    rollingUpdate: # 滚动升级的参数
      maxSurge: 1 # 在更新过程中最大同时创建 Pods 的数量
      maxUnavailable: 0 # 在更新过程中最大不可用的 Pods 的数量
  template: # 模版
    metadata: # 元数据
      labels: # 标签
        app: laravel 
    spec: # 对象规约
      containers: # 容器
        - name: laravel # # 容器名称
          image: jxlwqq/laravel-7-kubernetes-demo # 镜像
          ports: # 端口
            - containerPort: 80 # http web 服务默认端口 80
          env: # 环境变量，直接赋值
            - name: DB_HOST
              value : mysql-service # ！注意！这里的值，跟 MySQl Service 对象的 metadata name 一致
            - name: REDIS_HOST
              value: redis-service # ！注意！这里的值，跟 Redis Service 对象的 metadata name 一致
          envFrom: # 环境变量从哪里获得
            - configMapRef:
                name: laravel-env # 从名称为 laravel-env 的 configMap 对象中获得
          readinessProbe: # 就绪探针，或者叫做就绪探测器，集群通过它可以知道容器什么时候准备好了并可以开始接受请求流量
            httpGet: # 判断项目首页是不是返回 20X 状态
              port: 80
              path: /
              scheme: HTTP
          livenessProbe: # 存活探针，或者叫做存活探测器，集群通过它来知道什么时候需要重启容器
            httpGet: # 判断项目首页是不是返回 20X 状态
              port: 80
              path: /
              scheme: HTTP
---
kind: Service # 对象类型
apiVersion: v1 # api 版本
metadata: # 对象元数据
  name: laravel-service # 对象名称
  labels: # 标签
    app: laravel
spec: # 对象规约
  selector: # 选择器
    app: laravel # 选择那些包含 app: laravel 标签的一组 Pods
  ports: # 端口
    - port: 80 # 暴露的端口
      targetPort: 80 # 与 上面的 Pod containerPort 一致
 ```

ingress.yaml 文件解读：

```yaml
kind: Ingress # 对象类型
apiVersion: networking.k8s.io/v1beta1 # api 版本，这个资源对象还处于 beta 版本，但是已经很稳定了，k8s 社区还是比较谨慎的
metadata: # 元数据
  name: laravel-ingress # 对象名称
  labels: # 标签
    app: laravel
spec: # 对象规约
  rules:
    - http:
        paths:
          - backend:
              serviceName: laravel-service # 与 Laravel Service 对象的 metadata name 一致
              servicePort: 80 # Service 端口，与 Laravel Service 对象的 port 一致
 ```

## 创建 Ingress-nginx 控制器

有了 Ingress 对象还不够，还需要 Ingress-nginx 控制器。这里又有一个不太好的比方了，Ingress 对象类似 Nginx 的 nginx.conf 文件，单单有配置文件是万万不行的，我们需要 Nginx 服务（软件）本身。

为了让 Ingress 资源工作，集群必须有一个正在运行的 Ingress 控制器。 Kubernetes 官方目前支持和维护 GCE 和 nginx 控制器。

这里我们选择 Ingress-nginx 控制器：

```bash
cd ../ingress-nginx # 切换到 ingress-nginx 目录
kubectl apply -f ingress-nginx-deployment-and-other-resources-mandatory.yaml
kubectl apply -f ingress-nginx-service.yaml
```

注：
* ingress-nginx-deployment-and-other-resources-mandatory.yaml 文件内容来源自：https://github.com/kubernetes/ingress-nginx/blob/master/deploy/static/mandatory.yaml
* ingress-nginx-service.yaml 文件内容来源自：https://github.com/kubernetes/ingress-nginx/blob/master/deploy/static/provider/cloud-generic.yaml

详细操作说明见：https://github.com/kubernetes/ingress-nginx/blob/master/docs/deploy/index.md

## 访问

```bash
curl http://localhost
```
撒花，结束。

## 最后的话

现在的 Demo 将 web 服务、定时任务还有队列监听都放在了一个 Pod 中，无法对其进行扩容（因为定时任务和队列监听会重复）。如果需要对 Laravel 应用进行 HPA 扩容的话，还需要对 Laravel 项目的 docker/entrypint.sh 进行一些改造。将上述的 laravel-deploment 拆分成 3 个 Deploment，将容器分为三个角色，分为是 web、cron、queue。分别提供 web 服务、定时任务以及队列监听。最后对提供 web 服务的 Deploment 设置 HPA，根据 cpu 或者 内存占用率进行自动扩容。

另外，定时任务也可以使用 Kubernetes 的 CronJob 对象来实现。

这是后话，下期再细谈。