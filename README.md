# kubernetes-examples

<img src='./ingress.png' width='400'>

#### 说明

* 所有示例均在 Docker for Mac 内置的 Kubernetes 的集群上部署，请下载或更新到 Docker for Mac 最新版本。
* 动手实验前，请先执行 `kubectl config current-content` 查看当前集群名称，确保切换到本地的 Kubernetes 集群。
* [安装 Ingress 控制器：ingress-nginx](./ingress-nginx)是所有示例的前提条件。 
* 请勿用于生产环境。

#### 经典示例

* [ingress 路由使用示例](deploying-simple-apple-and-banana-ingress)
* [使用 Redis 部署 PHP 留言板应用程序](deploying-php-guestbook-application-with-redis)
* [使用 mongoDB 部署 PHP 留言板应用程序](deploying-php-guestbook-application-with-mongodb)
* [使用 Persistent Volumes 部署 WordPress 和 MySQL](deploying-wordpress-and-mysql-with-persistent-volumes)
* [部署单节点 Elasticsearch 与 Kibana](deploying-simple-elasticsearch-and-kibana)
* [使用 mongoDB 部署 NodeJs 笔记应用程序](deploying-nodejs-note-application-with-mongodb)
* [部署 Laravel Demo](deploying-laravel-application)
* [使用 MySQL 和 Redis 部署 Laravel 7 Demo](deploying-laravel-7-with-mysql-and-redis)
* [使用 PHP-FPM 和 Nginx 部署一个 简单的 phpinfo 应用](deploying-simple-php-app-with-fpm-and-nginx)