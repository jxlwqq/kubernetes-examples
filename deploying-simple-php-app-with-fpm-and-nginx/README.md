# 使用 php-fpm 和 nginx 部署一个简单的 phpinfo() 应用

##拉取镜像

源码在 php-info 目录中。我这里已经基于 Dockerfile 制作好了镜像，pull 后可以直接使用。

```yaml
docker pull jxlwqq/php-info
```

## 部署

```bash
kubectl apply -f configmap.yaml
kubectl apply -f php-fpm-nginx-deployment-and-service.yaml # 
kubectl apply -f horizontalpodautoscaler.yaml
kubectl apply -f ingress.yaml
```
也可以使用：
```bash
kubectl apply -f ./
```

```yaml
kind: ConfigMap
apiVersion: v1
metadata:
  name: nginx-config
data:
  nginx.conf: |
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
      containers:
        - name: php-fpm
          image: jxlwqq/php-info
          volumeMounts:
            - mountPath: /var/www/html
              name: nginx-www
          lifecycle:
            postStart:
              exec:
                command: ["/bin/sh", "-c", "cp -r /app/. /var/www/html"]
        - name: nginx
          image: nginx:1.7.9
          ports:
            - containerPort: 80
          volumeMounts:
            - mountPath: /var/www/html
              name: nginx-www
            - mountPath: /etc/nginx/nginx.conf
              subPath: nginx.conf
              name: nginx-config
      volumes:
        - name: nginx-www
          emptyDir: {}
        - name: nginx-config
          configMap:
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
      name: nginx
```

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
              serviceName: php-fpm-nginx
              servicePort: 80
```