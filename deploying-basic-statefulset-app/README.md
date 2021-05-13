
```shell
docker pull registry.aliyuncs.com/google_containers/nginx-slim:0.8
docker tag registry.aliyuncs.com/google_containers/nginx-slim:0.8 k8s.gcr.io/nginx-slim:0.8
docker rmi registry.aliyuncs.com/google_containers/nginx-slim:0.8
```

使用集群默认自带的 StorageClass

```shell
kubectl get storageclass
```

```shell
NAME                 PROVISIONER          RECLAIMPOLICY   VOLUMEBINDINGMODE   ALLOWVOLUMEEXPANSION   AGE
hostpath (default)   docker.io/hostpath   Delete          Immediate           false                  14d
```

```shell
kubectl apply -f nginx-statefulset-and-service.yaml
```

```yaml
apiVersion: v1
kind: Service
metadata:
  name: nginx
  labels:
    app: nginx
spec:
  ports:
    - port: 80
      name: web
  clusterIP: None
  selector:
    app: nginx
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: web
spec:
  serviceName: nginx
  replicas: 2
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
        - name: nginx
          image: k8s.gcr.io/nginx-slim:0.8
          ports:
            - containerPort: 80
              name: web
          volumeMounts:
            - name: www
              mountPath: /usr/share/nginx/html
  volumeClaimTemplates: # 没有指定 StorageName，则使用集群默认的存储类
    - metadata:
        name: www
      spec:
        accessModes: [ "ReadWriteOnce" ] 
        resources:
          requests:
            storage: 1Gi
```