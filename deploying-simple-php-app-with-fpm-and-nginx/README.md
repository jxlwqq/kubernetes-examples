# 使用 php-fpm 和 nginx 部署一个简单的 phpinfo() 应用

## 实验前提

* 需要你有 macOS 开发环境，本文以此为例，其他类型的开发环境请自行搭建。
* 需要你对 YAML 这一专门用来写配置文件的语言有所了解。
* 需要你对 Docker 有一些基本的了解。
* 需要你对 Kubernetes 中的 Node、Pod、ReplicaSet、Deployment、Service、Ingress、ConfigMap 等一些核心基础概念有一定的了解。

## YAML 配置文件下载地址：

* YAML 文件：[jxlwqq/kubernetes-examples](https://github.com/jxlwqq/kubernetes-examples/tree/master/deploying-laravel-7-with-mysql-and-redis)。该项目还有其他一些 Kubernetes 的示例。欢迎 Star。

```bash
git clone https://github.com/jxlwqq/kubernetes-examples.git
cd deploying-simple-php-app-with-fpm-and-nginx
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
kubectl config use-content docker-desktop
````

## 拉取镜像

源码在 php-info 目录中。我这里已经基于 Dockerfile 制作好了镜像，pull 后可以直接使用。

```yaml
docker pull jxlwqq/php-info
```

源码逻辑很简单，打印 phpinfo 信息，Dockerfile 内容如下所示：

```Dockerfile
FROM php:7.4-fpm
WORKDIR /app
COPY index.php /app
```

```php
<?php
    phpinfo();
```

## 部署

```bash
kubectl apply -f configmap.yaml # 配置对象，本示例存放 nginx.config
kubectl apply -f php-fpm-nginx-deployment-and-service.yaml # php-fpm 和 nginx 双容器
kubectl apply -f ingress.yaml # ingress 路由规则
```

configmap.yaml 文件解读：

```yaml
kind: ConfigMap
apiVersion: v1
metadata:
  name: nginx-config
data:
  nginx.conf: | # 将 nginx config 配置写入 ConfigMap 中，经典的 php-fpm 代理设置，这里就不再多说了
    events {
    }
    http {
      server {
        listen 80 default_server;
        listen [::]:80 default_server;
        root /var/www/html;
        index index.php;
        server_name _;
        location / {
          try_files $uri $uri/ =404;
        }
        location ~ \.php$ {
          include fastcgi_params;
          fastcgi_param REQUEST_METHOD $request_method;
          fastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name;
          fastcgi_pass 127.0.0.1:9000;
        }
      }
    }
```

php-fpm-nginx-deployment-and-service.yaml 文件解读：

```yaml
kind: Deployment
apiVersion: apps/v1
metadata:
  name: php-fpm-nginx
spec:
  selector:
    matchLabels:
      app: php-fpm-nginx
  replicas: 1
  template:
    metadata:
      labels:
        app: php-fpm-nginx
    spec:
      containers: # 这里设置了多个容器
        - name: php-fpm
          image: jxlwqq/php-info
          ports:
            - containerPort: 9000 # php-fpm 端口
          volumeMounts: 
            - mountPath: /var/www/html # 挂载两个容器共享的 volume 
              name: nginx-www
          lifecycle: # 生命周期
            postStart: # 当容器处于 postStart 阶段时，执行一下命令
              exec:
                command: ["/bin/sh", "-c", "cp -r /app/. /var/www/html"] # 将 /app/index.php 复制到挂载的volume里
        - name: nginx
          image: nginx
          ports:
            - containerPort: 80 # nginx 端口
          volumeMounts:
            - mountPath: /var/www/html # 挂载两个容器共享的 volume 
              name: nginx-www
            - mountPath: /etc/nginx/nginx.conf #  挂载配置了 nginx.conf 的 volume
              subPath: nginx.conf
              name: nginx-config
      volumes:
        - name: nginx-www # 这个 volume 是 php-fpm 容器 和 nginx 容器所共享的，两个容器都 volumeMounts 了
          emptyDir: {}
        - name: nginx-config 
          configMap: # 这里为啥可以将 configMap 对象通过 volumeMounts 的方式注入到容器中呢，因为本质上 configMap 是一类特殊的 volume
            name: nginx-config
---
kind: Service
apiVersion: v1
metadata:
  name: php-fpm-nginx
spec:
  selector:
    app: php-fpm-nginx
  ports:
    - port: 80
      targetPort: 80
```

ingress.yaml 文件解读：

```yaml
kind: Ingress
apiVersion: networking.k8s.io/v1beta1
metadata:
  name: php-fpm-nginx
spec:
  rules:
    - http:
        paths:
          - backend:
              serviceName: php-fpm-nginx # 流量转发到名为 php-fpm-nginx 的 Server 是那个
              servicePort: 80
```

## 自动伸缩

```bash
kubectl apply -f horizontalpodautoscaler.yaml # hpa 水平自动伸缩对象
```

horizontalpodautoscaler.yaml 文件解读：

```yaml
apiVersion: autoscaling/v2beta2
kind: HorizontalPodAutoscaler
metadata:
  name: php-fpm-nginx
spec:
  scaleTargetRef: # 扩容的目标
    apiVersion: apps/v1
    kind: Deployment # 目标对象的类型
    name: php-fpm-nginx # 目标对象的名称
  minReplicas: 3 # 最小副本数
  maxReplicas: 10 # 最大副本书
  metrics: # 指标）
    - type: Resource # 类型：资源
      resource:
        name: memory # 内存
        target:
          type: Utilization # 利用率
          averageUtilization: 1 # 1% 这个值是为了实验，具体值请参考业务方实际情况而定
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

## 清场

删除本次示例所有的对象：

```bash
kubectl delete -f ./
```