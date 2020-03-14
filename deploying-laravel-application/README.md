# 轻松部署 Laravel 应用

原文链接：[Kubernetes: deploy Laravel the easy way](https://learnk8s.io/blog/kubernetes-deploy-laravel-the-easy-way)
译文链接：[使用 Kubernetes 来部署你的 Laravel 程序](https://learnku.com/server/t/22017)

注意：原文使用了 minikube 来部署，本文使用的是 Docker for Mac 自带的 Kubernetes 集群。

#### 拉取镜像 

基于原文，制作了一个 Docker 镜像 [laravel-kubernetes-demo](https://hub.docker.com/repository/docker/jxlwqq/laravel-kubernetes-demo)
，方便大家快速拉取镜像：

```bash
docker pull jxlwqq/laravel-kubernetes-demo
```

#### 部署

```bash
kubectl apply -k ./
```

上述命令会应用本目录下的 kustomization.yaml 文件。它的信息如下：
```yaml
resources: # 需要apply的资源文件
  - laravel-deployment-and-service.yaml
  - ingress.yaml

configMapGenerator: # 生成一个 ConfigMap 对象
  - name: laravel-env
    literals:
      - APP_KEY=base64:zC8wVldUZfZJaGaZ7+CPh+5FzaXYmShm7G/Qh6GdRl8=
```

#### 部署 Ingress-nginx 控制器

```bash
cd ../ingress-nginx
kubectl apply -k ./ 
```

#### 访问
```bash
curl http://localhost
```

#### 清理
```bash
kubectl delete -k ./
```
