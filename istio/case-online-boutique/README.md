```shell
docker pull gcr.io/google-samples/microservices-demo/emailservice:v0.3.0
docker pull gcr.io/google-samples/microservices-demo/checkoutservice:v0.3.0
docker pull gcr.io/google-samples/microservices-demo/recommendationservice:v0.3.0
docker pull gcr.io/google-samples/microservices-demo/frontend:v0.3.0
docker pull gcr.io/google-samples/microservices-demo/paymentservice:v0.3.0
docker pull gcr.io/google-samples/microservices-demo/productcatalogservice:v0.3.0
docker pull gcr.io/google-samples/microservices-demo/cartservice:v0.3.0
docker pull gcr.io/google-samples/microservices-demo/loadgenerator:v0.3.0
docker pull gcr.io/google-samples/microservices-demo/currencyservice:v0.3.0
docker pull gcr.io/google-samples/microservices-demo/shippingservice:v0.3.0
docker pull gcr.io/google-samples/microservices-demo/adservice:v0.3.0
```


```shell
git clone git@github.com:GoogleCloudPlatform/microservices-demo.git
kubectl apply -f release/kubernetes-manifests.yaml
kubectl apply -f release/istio-manifests.yaml
```