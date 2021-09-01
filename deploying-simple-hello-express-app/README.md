# 部署一个简单的 Express 应用

#### 源代码

一个最简单的 Express 应用程序看起来是这样的：

```js
const express = require('express')
const app = express()
const port = 3000

app.get('/', (req, res) => {
    res.send('Hello Express!')
})

app.listen(port, () => {
    console.log(`Example app listening at http://localhost:${port}`)
})
```

package.json 文件包含项目所需的依赖。

#### Docker 镜像

应用的 Dockerfile 如下所示：

```dockerfile
# 从官方仓库中获取最新版的 Node 基础镜像
FROM node:14
# 设置工作目录
WORKDIR /usr/src/app
# 复制项目依赖文件
ADD package*.json .
# 安装依赖
RUN npm install
# 复制项目文件
ADD . .
# 设置监听端口
EXPOSE 3000
# 配置启动命令
CMD [ "node", "index.js" ]
```

构建并提交镜像：

> jxlwqq 是我的 Docker Hub 账号，这里需要换成你自己的账号。

```shell
docker build -f Dockerfile -t jxlwqq/hello-express:1.0.0 . # 构建镜像
docker push jxlwqq/hello-express:1.0.0 # 提交镜像
```

#### 前提条件：部署 nginx ingress

```bash
cd ../ingress-nginx # 切换到 ingress-nginx 目录
kubectl apply -f deploy.yaml
```

#### 部署 hello express 应用

执行以下命令：

```shell
kubectl apply -f hello-express-deployment-and-service.yaml
kubectl apply -f ingress.yaml
```

`hello-express-deployment-and-service.yaml` 文件解读：

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: hello-express
spec:
  selector:
    matchLabels:
      app.kubernetes.io/name: hello-express # 选择匹配的 Pod 标签
      app.kubernetes.io/version: 1.0.0
  template:
    metadata:
      name: hello-express
      labels:
        app.kubernetes.io/name: hello-express # 选择匹配的 Pod 标签
        app.kubernetes.io/version: 1.0.0
    spec:
      containers:
        - name: hello-express
          image: jxlwqq/hello-express:1.0.0 # 镜像名称:镜像版本
          ports:
            - containerPort: 3000
---
apiVersion: v1
kind: Service
metadata:
  name: hello-express-svc
  labels:
    app.kubernetes.io/name: hello-express
spec:
  selector:
    app.kubernetes.io/name: hello-express # 选择匹配的 Pod 标签
  type: ClusterIP
  ports:
    - port: 3000
      targetPort: 3000
```

`ingress.yaml` 文件解读：

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: hello-express-ingress
spec:
  rules:
    - http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: hello-express-svc # service 名称
                port:
                  number: 3000 # 端口号
  ingressClassName: nginx
```

访问验证：

```shell
curl 127.0.0.1 # 返回 <p>Hello, Express!</p>
```

#### 清理
```shell
kubectl delete -k .
```

