```shell

docker pull mysql:5.7
# docker pull gcr.io/google-samples/xtrabackup:1.0
docker pull ist0ne/xtrabackup:1.0
docker tag ist0ne/xtrabackup:1.0 gcr.io/google-samples/xtrabackup:1.0
docker rmi ist0ne/xtrabackup:1.0
```

```shell
kubectl apply -f .
```