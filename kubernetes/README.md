# Kubernetes


# Minikube installation

```sh
# Installation instructions at https://github.com/kubernetes/minikube/releases:
# At the time of writing
curl -Lo minikube https://storage.googleapis.com/minikube/releases/v0.30.0/minikube-darwin-amd64 && chmod +x minikube && sudo cp minikube /usr/local/bin/ && rm minikube

# Starting minikube
minikube start
# Ensure 'docker' env variables (e.g. DOCKER_HOST) point to minikube's Docker. 
# Make sure to execute this command in every terminal that is being used
eval $(minikube docker-env)

# Enable kubectl autocompletion
source <(kubectl completion bash)

# Other commands:
minikube status
minikube stop
```

# Docker image building

These are instructions to build the docker images of the gRPC app. Make sure to run ```eval $(minikube docker-env)```
first if you want to do this with Kubernetes.

Note: Dockerfiles cannot add files using ADD or COPY from files outside of the directory the Dockerfile sits in.
That's why the Dockerfiles we will use for this will be sitting together with the source-code of the gRPC samples.

# Docker image details
## Backend

```sh
# Make sure you're in the backend dir before building
cd grpc/backend
docker build -t backend .

# Running the backend image with docker directly:
docker run -ti -p 50051:50051 backend

# FYI, that the images are based on alpine linux, which comes with the [ash](https://en.wikipedia.org/wiki/Almquist_shell) shell.
# Run it like so (instead of /bin/bash)
docker run -ti backend "/bin/sh"
```

## Web

Regular Dockerfile building:
```bash
# Make sure we're in the web directory, needed to build the web binary
cd grpc/web
docker build -t web .

# Running web image with docker directly
docker run -ti -p 1234:1234 web

# Note that if you're running using regular docker (non-minikube), you cna just use curl at this point
curl localhost:1234
```

## Fun with Kubernetes

### Single container deployments
```sh
# We need to set --image-pull-policy=Never to ensure that k8s will search our local images.
# More info: https://stackoverflow.com/questions/42564058/how-to-use-local-docker-images-with-minikube
# Format
# kubectl run <deployment-name> --image=<image> ...
kubectl run hello-web --image=web --port=1234 --image-pull-policy=Never
# Expose deployment
kubectl expose deployment hello-web --type=NodePort

# List info
kubectl get deployments
kubectl get pods # a deployment can contain multiple pods

# To delete (pods will get deleted automatically)
kubectl delete deployment hello-web
# Note that if you just delete a single pod, it will get recreated (that's a k8s feature!), you need to delete the deployment

# Since this is running locally, minikube needs to do some vodoo with networking. To actually query the URL, do this:
curl $(minikube service hello-web --url)
# Should return: "This is the index page!"
```

# Multi-container deployments

TODO