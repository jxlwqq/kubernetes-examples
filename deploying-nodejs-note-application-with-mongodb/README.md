#### 示例来源

[Hands-on guide: developing and deploying Node.js apps in Kubernetes](https://learnk8s.io/nodejs-kubernetes-guide)

本仓库的示例均使用 Ingress 进行负载均衡，所以略有调整。也可以点击链接按照原文步骤进行安装和实验。

#### 镜像准备

```bash
docker pull learnk8s/knote-js:1.0.0 # 源码地址：https://github.com/learnk8s/knote-js/tree/master/01
docker pull mongo
```
#### 部署服务
```bash
kubectl apply -f frontend-deployment-and-service.yaml
kubectl apply -f mongo-deployment-and-service.yaml
kubectl apply -f ingress.yaml

#############
# 或者使用一个命令进行部署
kubectl apply -k ./
#############
```

#### 访问

打开浏览器访问：`http://localhost`


#### 清理
```shell
kubectl delete -k ./
```