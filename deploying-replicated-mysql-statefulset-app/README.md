# 运行一个有状态的 MySQL 服务

原文：https://kubernetes.io/zh/docs/tasks/run-application/run-replicated-stateful-application/

#### 镜像准备

```shell
docker pull mysql:5.7
docker pull gcr.io/google-samples/xtrabackup:1.0
```

#### 部署

```shell
kubectl apply -f mysql.yaml
kubectl get pods -w
# 缩容
kubectl scale statefulset mysql --replicas=2
# 扩容
kubectl scale statefulset mysql --replicas=3
```

#### 查询操作

创建数据库：

```shell
# 直连写入
kubectl run mysql-client --image=mysql:5.7 -i --rm --restart=Never --\
  mysql -h mysql-0.mysql <<EOF
CREATE DATABASE test;
CREATE TABLE test.messages (message VARCHAR(250));
INSERT INTO test.messages VALUES ('hello');
EOF
```

查询：

```shell
# 通过 service 负载均衡查询
kubectl run mysql-client --image=mysql:5.7 -i -t --rm --restart=Never --\
  mysql -h mysql-read -e "SELECT * FROM test.messages"
  
# 直连查询
kubectl run mysql-client --image=mysql:5.7 -i -t --rm --restart=Never --\
  mysql -h mysql-0.mysql -e "SELECT * FROM test.messages"
 
# 直连查询 
kubectl run mysql-client --image=mysql:5.7 -i -t --rm --restart=Never --\
  mysql -h mysql-0.mysql -e "SELECT * FROM test.messages"
  
# 直连查询
kubectl run mysql-client --image=mysql:5.7 -i -t --rm --restart=Never --\
  mysql -h mysql-0.mysql -e "SELECT * FROM test.messages"
```

#### 清理

```shell
kubectl delete -f .
```

