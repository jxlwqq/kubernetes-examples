# 使用 Persistent Volumes 部署 WordPress 5 和 MySQL


#### 前提

ingress-nginx 控制器已经运行。

```bash
cd ingress-nginx/ # 切换到 ingress-nginx 目录
kubectl apply -f deploy.yaml
```

#### 准备镜像

```bash
docker pull mysql:5.7
docker pull wordpress:5-php8.0-apache
```

#### 安装

```shell
cd deploying-wordpress-and-mysql-with-persistent-volumes/
kubectl apply -k .
```

注意：

> The WORDPRESS_DB_NAME needs to already exist on the given MySQL server; it will not be created by the wordpress container. [source](https://hub.docker.com/_/wordpress)

5 版本的 WordPress 将不在自动创建数据库，所以这里我们需要手动登录 mysql 的容器中来 create 数据库。

```shell
kubectl exec mysql-deployment-5d585f4476-h7rgr -it -- mysql -uroot -p # mysql-deployment-5d585f4476-h7rgr pod名称替换成你自己的
```

密码: `!@#123`

```mysql
create database `wordpressdb`;
```

#### 查看各类资源 

```shell
kubectl get secret  # 查看 secret
kubectl get pod     # 查看 pod
kubectl get deploy  # 查看 deployment
kubectl get svc     # 查看 service
kubectl get ingress # 查看 ingress
kubectl get sc      # 查看 storageClass
kubectl get pv      # 查看 persistentVolume
kubectl get pvc     # 查看 persistentVolumeClaim
```

#### 设置 WordPress

浏览器打开 http://localhost:80 ，设置用户名和密码：

![wordpress:hello-world](hello-world.png)


#### 清场

```shell
kubectl delete -k .
```



