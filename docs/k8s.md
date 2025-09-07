# Kubernetes

Simple Kubernetes manifests can be found in the `k8s` folder. The configuration there, as well as the instructions in this document, are designed for a local bare-metal cluster (particularly using minikube) for testing how the distributed system behaves across multiple replicas.

## Configuration
There are 3 mainfests in the `k8s` folder:
- `app.yaml`: The `Deployment` and `Service` for the `session-queue` app
- `redis.yaml`: Simple Redis `Deployment` and `Service` for the `session-queue` replicas.
- `config.yaml`: Environment variable configuration for `session-queue` with some sensible testing defaults.

## Run Locally (minikube)
Make sure minikube is running on your machine with:
```bash
minikube start
```
The app manifest points at a local image rather than pulling from a remote registry for local testing ease. Build the image in minikube's environment with the following:
```bash
eval $(minikube docker-env)
docker build -t session-queue:latest .
```
Now with the image built and available, you should be able to apply all the manifests, with:
```bash
kubectl apply -f k8s
```
The URL of `session-queue`'s `NodePort` `Service` can now be retrieved with minikube:
```bash
minikube service session-queue --url
```
You should now be able to interact with the system, using a HTTP client of your choice, with this URL. For example, with HTTPie:
```bash
http POST {url}/join
http GET {url}/status "Authorization: Bearer {token}"
```
You can also access logs of the replicas with:
```bash
kubectl get pods # This will list random-generated pod names
kubectl logs {pod} --follow
```

## Production
As mentioned, the out-the-box mainfests are made for local development testing. There are various changes that could/should be made for a production-ready k8s deployment, including:
- Pushing the image to a remote registry and not using `imagePullPolicy: Never`
- Using a `Service` of type `LoadBalancer` to expose the deployment (particularly if on a cloud environment)
- Storing the JWT secret outside of config/source control, and using something like [ESO](https://external-secrets.io/latest/) to inject it