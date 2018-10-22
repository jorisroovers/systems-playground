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

# Enable kubectl autocompletion (not sure if working on mac...)
source <(kubectl completion bash)

# Other commands:
minikube status
minikube stop
# Starting with a clean kubernetes cluster (takes a while):
minikube delete
rm -rf ~/.minikube
eval $(minikube docker-env)
```

# Docker image building

These are instructions to build the docker images of the gRPC app. Make sure to run ```eval $(minikube docker-env)```
first if you want to do this with Kubernetes.

What's happening here, is that minikube is running a docker instance in a VM. By running ```eval $(minikube docker-env)```,
we're pointing our docker command to the docker instance in the VM.

When we now build docker images, we will be doing so within the minikub environment, and kubernetes will be able to find
them under the name we give to the images.

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

# Fun with Kubernetes

Make sure to build the images with docker (pointing to the minikube docker instace) as described in the previous section.

## Intro to k8s

K8s has many different types of resources (service, deployment, pod, nodes). All can be created via yaml files, all can
be manipulated using the same set of commands. List all resource types with ```kubectl get --help```.

Basic resources:
- Pods: scalable containers
- Deployments: group of pods and their interconnectivity
- Services: Way to expose containers to the outside world
- Nodes: Infrastructure where k8s is hosting your containers on

Great for hands-on examples: http://kubernetesbyexample.com

## Single container deployments


```sh
# We need to set --image-pull-policy=Never to ensure that k8s will search our local images.
# More info: https://stackoverflow.com/questions/42564058/how-to-use-local-docker-images-with-minikube
# Format
# kubectl run <deployment-name> --image=<image> ...
kubectl run hello-web --image=web --port=1234 --image-pull-policy=Never
# Same, but using a file definition (only creates a pod, not a deployment):
kubectl create -f simple-pod.yaml

# List info
kubectl get deployments
kubectl get pods # a deployment can contain multiple pods
kubectl get services  # services expose ports from deployments/pods for external use

# Resource specific details
kubectl describe pod hello-web
kubectl describe pod web-from-file
kubectl describe deployment hello-web # only exists for hello-web, not for web-from-file

# Expose deployment, this will create a 'hello-web' k8s service.
kubectl expose deployment hello-web --type=NodePort
kubectl expose pod web-from-file --type=NodePort # expose a pod

# Details on created service
kubectl describe service hello-web
kubectl describe service web-from-file

# Since this is running locally, minikube needs to do some vodoo with networking. To actually query the URL, do this:
curl $(minikube service hello-web --url)
curl $(minikube service web-from-file --url)
# Should return: "This is the index page!"

# To delete (pods will get deleted automatically)
kubectl delete deployment hello-web
# Note that if you just delete a single pod, it will get recreated (that's a k8s feature!), you need to delete the deployment
kubectl delete pod web-from-file
# services need to be deleted separately!
kubectl delete service hello-web
kubectl delete service web-from-file
```

## Multi-container deployments

JR: Continue here :-) multi-pod-deployment.yaml doesn't work yet.

```shell
kubectl create -f multi-pod-deployment.yaml
```

## Misc Notes

```shell
# K8s config
kubectl config view

```