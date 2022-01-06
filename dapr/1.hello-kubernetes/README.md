# Hello World

#### 部署

```shell
kubectl apply -f deploy.yaml
```

#### 查看日志

```shell
kubectl logs --selector=app=node -c node --tail=-1
```