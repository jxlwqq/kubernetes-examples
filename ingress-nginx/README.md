#### 安装

为了让 Ingress 资源工作，集群必须有一个正在运行的 Ingress 控制器。Kubernetes 作为一个项目，目前支持和维护 AWS， GCE 和 nginx Ingress 控制器。这里我们推荐安装 nginx Ingress 控制器。

```bash
kubectl apply -f deploy.yaml
```

注：

deploy.yaml 文件内容来源自：https://github.com/kubernetes/ingress-nginx/blob/main/deploy/static/provider/cloud/deploy.yaml

详细操作说明见：https://github.com/kubernetes/ingress-nginx/blob/main/docs/deploy/index.md