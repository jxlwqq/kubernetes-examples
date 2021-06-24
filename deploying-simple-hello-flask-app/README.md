# 部署一个简单的 Flask 应用

#### 源代码

一个最迷你的 Flask 应用程序看起来是这样的：

```python
from flask import Flask

app = Flask(__name__)

@app.route("/")
def hello_world():
    return "<p>Hello, Flask!</p>"
    
if __name__ == "__main__":
    app.run(host='0.0.0.0')
```

requirements.txt 文件包含 app.py 所需的依赖，pip 将使用它来安装 Flask 包。

#### Docker 镜像

应用的 Dockerfile 如下所示：

```dockerfile
# 从官方仓库中获取最新版的 Python 基础镜像
FROM python
# 设置工作目录
WORKDIR /
# 复制项目文件
ADD . /
# 安装依赖
RUN pip install -r requirements.txt
# 设置监听端口
EXPOSE 5000
# 配置启动命令
CMD ["python", "app.py"]
```

构建并提交镜像：

> jxlwqq 是我的 Docker Hub 账号，这里需要换成你自己的账号。

```shell
docker build -f Dockerfile -t jxlwqq/hello-flask:latest . # 构建镜像
docker push jxlwqq/hello-flask:latest # 提交镜像
```

#### 前提条件：部署 nginx ingress

```bash
cd ../ingress-nginx # 切换到 ingress-nginx 目录
kubectl apply -f deploy.yaml
```

#### 部署 hello flask 应用

执行以下命令：

```shell
kubectl apply -f hello-flask-deployment-and-service.yaml
kubectl apply -f ingress.yaml
```

`hello-flask-deployment-and-service.yaml` 文件解读：

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: hello-flask
spec:
  selector:
    matchLabels:
      name: hello-flask # 选择匹配的 Pod 标签
  template:
    metadata:
      name: hello-flask
      labels:
        name: hello-flask # Pod 的标签
    spec:
      containers:
        - name: hello-flask
          image: jxlwqq/hello-flask:latest # 镜像名称:镜像版本
          ports:
            - containerPort: 5000
---
apiVersion: v1
kind: Service
metadata:
  name: hello-flask-svc
spec:
  selector:
    name: hello-flask # 选择匹配的 Pod 标签
  ports:
    - port: 80
      targetPort: 5000
```

`ingress.yaml` 文件解读：

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: hello-flask-ingress
spec:
  rules:
    - http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: hello-flask-svc # service 名称
                port:
                  number: 80 # 端口号
```

访问验证：

```shell
curl 127.0.0.1 # 返回 <p>Hello, Flask!</p>
```

#### 清理
```shell
kubectl delete -k .
```

