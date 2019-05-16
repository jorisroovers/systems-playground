# Kubernetes


# Minikube installation

```sh
# Installation instructions at https://github.com/kubernetes/minikube/releases:
# At the time of writing
curl -Lo minikube https://storage.googleapis.com/minikube/releases/v0.30.0/minikube-darwin-amd64 && chmod +x minikube && sudo cp minikube /usr/local/bin/ && rm minikube
brew install kubernetes-cli

# Starting minikube
minikube start
# Alternatively, enable swagger UI support for API browsing:
minikube start --extra-config=apiserver.enable-swagger-ui=true

# Ensure 'docker' env variables (e.g. DOCKER_HOST) point to minikube's Docker.
# Make sure to execute this command in every terminal that is being used
eval $(minikube docker-env)

# Enable kubectl autocompletion (not sure if working on mac...)
source <(kubectl completion bash)

# Other commands:
minikube status
minikube stop
minikube ssh  # SSH into the minkube VM

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

# Note that if you're running using regular docker (non-minikube), you can just use curl at this point
curl localhost:1234
```

# Fun with Kubernetes

Make sure to build the images with docker (pointing to the minikube docker instance) as described in the previous section.

## Intro to k8s

K8s has many different types of resources (service, deployment, pod, nodes). All can be created via yaml files, a bunch of them can also 
be created/manipulated using CLI commands (but not all resource types - some require yaml files). 
List all resource types with ```kubectl get --help```.

Basic resources:
- Pods: collection of scalable containers. A pod typically maps onto a single micro-service (which might be comprised of multiple containers).
- Deployments: supervisor for pods with fine-grained control over how and when a new pod version is rolled out as well as rolled back to a previous state.
- Services: Way to expose containers to the outside world
- Nodes: Infrastructure where k8s is hosting your containers on
- Secrets: Secure storage for sensitive information for easy consumption in your k8s cluster
- Volumes: shared storage

Great for hands-on examples: http://kubernetesbyexample.com

Under-the-hood processes:
- kubectl: runs on every Node to provide computing
- kube-proxy: runs on every Node to do mappings between VIPs and pods via iptables and IP namespaces.
- kube-scheduler
- kube-controller-manager
- kube-dns
- sidecar

## Notes
```sh
# Show available nodes:
kubectl get nodes
kubectl describe nodes

```

## Single container deployments

### Creation
```sh
# We need to set --image-pull-policy=Never to ensure that k8s will search our local images.
# More info: https://stackoverflow.com/questions/42564058/how-to-use-local-docker-images-with-minikube
# Format
# kubectl run <deployment-name> --image=<image> ...
kubectl run hello-web --image=web --port=1234 --image-pull-policy=Never
# Same, but also add labels. Labels make it easy to reference back to resources later on.
kubectl run hello-web --image=web --port=1234 --image-pull-policy=Never --labels "foo=bar,hurr=durr"
# Same, but using a file definition (only creates a pod, not a deployment):
kubectl create -f simple-pod.yaml

# Same thing (apply allows you to apply file against existing resources as well):
kubectl apply -f simple-pod.yaml

# More complex example file
kubectl apply -f complex-pod.yaml
```

### Inspection
```bash
kubectl get deployments
kubectl get pods # a deployment can contain multiple pods
kubectl get services  # services expose ports from deployments/pods for external use
kubectl get pods -l foo=bar # only pods with label 'foo=bar'
kubectl get pods -l 'foo in (bar, bazz)' # more complex selector

# Resource specific details
kubectl describe pod hello-web
kubectl describe pod web-from-file
kubectl describe deployment hello-web # only exists for hello-web, not for web-from-file

# Execute a command - very similar to docker. This works only for pods.
kubectl exec -it hello-web /bin/sh

# Get logs for a pod - very similar to docker. This works only for pods.
kubectl logs -f hello-web
```

### Services
```bash
# Expose deployment, this will create a 'hello-web' k8s service.
# --type=NodePort => expose deployment/deployment on its node (=host system)
kubectl expose deployment hello-web --type=NodePort
kubectl expose pod web-from-file --type=NodePort # expose a pod

# Details on created service
kubectl describe service hello-web
kubectl describe service web-from-file
# List services
kubectl get services # shorthand=svc

# Since this is running locally, minikube needs to do some voodoo with networking. To actually query the URL, do this:
curl $(minikube service hello-web --url)
curl $(minikube service web-from-file --url)
# Should return: "This is the index page!"
```

### Removal
```bash

# To delete (pods will get deleted automatically)
kubectl delete deployment hello-web
# Note that if you just delete a single pod, it will get recreated (that's a k8s feature!), you need to delete the deployment
# For the web-from-file deployment, we didn't create a deployment, so we can just remove the pod
kubectl delete pod web-from-file
# services need to be deleted separately!
kubectl delete service hello-web
kubectl delete service web-from-file

# You can delete multiple resources of the same type at once by just listing them
kubectl delete pod web-from-file web-from-file-complex

# Delete multiple resources of different types by specifying their resource type as well: 
# kubectel delete <resource-type>/<resource-name>
kubectl delete pod/web-from-file secret/apikey

```

## Deployments and rollouts
Deployments and roll-outs make it easy to do upgrades of applications.

```sh
# Deploy version 1 of the web app (my-app), expose the port and get index page
# Returned version is determined by the SERVICE_VERSION envar set in the yaml file.
kubectl apply -f deployments/deployment-v1.yaml
kubectl expose deployment my-app --type=NodePort
curl $(minikube service my-app --url)

# Roll out a newer version of the my-app
kubectl apply -f deployments/deployment-v2.yaml
curl $(minikube service my-app --url) # will show 'version 2'

# Show all version of my-app
kubectl rollout history deploy/my-app

# Roll back to a previous version
kubectl rollout undo deploy/my-app --to-revision=1
curl $(minikube service my-app --url) # shows 'version 1' again

# Roll out v3, which has 10 replicas configured. Use 'rollout status' to follow along.
kubectl apply -f deployments/deployment-v3.yaml
kubectl rollout status deploy/my-app
kubectl get pods
```

## Scaling (replica sets)
A k8s Deployment under-the-hood contains a Replicaset which is a wrapper around a Pod that allows it to horizontally scale the pod.
While you can manually create the ReplicaSet resource type, there's usually no reason to do so directly.
Instead, when you use a Deployment, k8s will automatically create the replicaset for you.

```sh
# Create deployment, show associated replicaset that k8s created
kubectl apply -f deployments/deployment-v1.yaml
kubectl get rs

# Scale deployment up, to 3 pods
kubectl scale --replicas 3 deploy/my-app
kubectl scale --replicas 2 deploy/my-app # scale down


```

TODO: related to replica sets



## Multi-container deployments

JR: Continue here :-) multi-container-pod.yaml doesn't work yet.

```bash
# Create multi-container deployment
kubectl create -f multi-container-deployment.yaml

# Get details on containers that are part of the pod
kubectl describe deployment multi-container-deployment

#
kubectl get replicaset
kubectl get rs

kubectl get pods

# Run a command in a specific container using the -c flag
kubectl exec <pod-name> -c backend ls
kubectl exec <pod-name> -c web ls


# Show logs of the backend container in the pod for this deployment
kubectl logs <pod-name> backend
```


## Volumes
Volumes can be shared across different containers.
Multiple types of containers exist: emptyDir, hostpath, awsElasticBlockStore, cephfs, nfs, etc. There's a long list.

```bash
# 'hostPath' volume type = shared directory between pod and host system.
# For the 'hostPath' volume type in pod-with-volumes.yaml, we need to make sure we have the directory created
# on the host system.
mkdir /tmp/k8s-shared
# We also need to mount this directory in minikube, because the containers are really running inside the minikube VM
# and not directly on our host machine.
# This command needs to run in a seperate terminal and stay running.
minikube mount /tmp/k8s-shared:/tmp/k8s-shared

# Create sample pod
kubectl create -f pod-with-volumes.yaml

# List files in 'emptyDir' type volume on container
kubectl exec pod-with-volumes -c container1 -- ls /tmp/foobar/
kubectl exec pod-with-volumes -c container2 -- ls /tmp/hurdur

# Write something to volume in first container
kubectl exec -t pod-with-volumes -c container1  -- sh -c "echo pizza > /tmp/foobar/test"
# Read it in the second container
kubectl exec pod-with-volumes -c container2 -- cat /tmp/hurdur/test

#  Create some test file in our host system
echo "testfile" > /tmp/k8s-shared/hostpath-test
# Read file on container
kubectl exec pod-with-volumes -c container1 -- cat /tmp/shared1/hostpath-test
```

## Persistent volumes
Shared across the cluster (=across nodes) and across node reboots. Like regular volumes, many types exist: nfs, GCEPersistentDisk, AWSElasticBlockStore, HostPath, CephFS, etc.

TODO

```bash

kubectl create -f persistent-volumes/persistent-volume.yaml
kubectl create -f persistent-volumes/persistent-volume-claim.yaml


kubectl get persistentvolumes
kubectl get pv

kubectl get persistentvolumeclaims
kubectl get pvc

```

## Namespaces
Way to group k8s resources together. Like a tenant/project on other IaaS platforms.

```bash
# Create and list namespaces
kubectl create namespace foobar
kubectl get namespace

# By default, kubectl commands run against the 'default' namespace, but you can always specify the namespace explicitly
kubectl get pods --namespace foobar
kubectl get pods -n foobar # shorthand

# Resources can exist with the same name in different namespaces
kubectl apply -f simple-pod.yaml # default namespace
kubectl apply --namespace foobar -f simple-pod.yaml
```

## Secrets

### Creation
```bash
# From CLI:
kubectl create secret generic mysecret --from-literal=key1=supersecret

# Note that secrets are key-value pairs, so you can specify multiple key-value pairs:
kubectl create secret generic mysecret2 --from-literal=key2=supersecret2 --from-literal=key2=supersecret2

# From File
echo -n "foobar" > /tmp/mysecret.txt
kubectl create secret generic mysecret-fromfile --from-file=key1=/tmp/mysecret.txt

```

### Inspection
```bash
# Basic info on secret (won't show the actual value)
kubectl get secret mysecret
kubectl get secret mysecret2
kubecetl describe mysecret

# To get value, output as yaml. Note that secret is still base64 encoded
kubectl get secret mysecret2 -o yaml
# Extract secret using awk + base64 --decode
kubectl get secret mysecret2 -o yaml | awk -F ":" '/key1/{print $2}' | base64 --decode
```

### Consuming in pod
```bash
# Create secret, start pod
kubectl create secret generic mysecret --from-literal=username=mysecretusername --from-literal=mysecret=foobar
kubectl apply -f pod-with-secret.yaml

# Show env variable, containing SECRET_USERNAME
kubectl exec pod-with-secret env
# Show volume mount containing secret
kubectl exec pod-with-secret cat /tmp/mysecrets/mysecret

# Update a previously set secret
# Bit of a hack to update resources: https://stackoverflow.com/questions/45879498/how-can-i-update-a-secret-on-kubernetes-when-it-is-generated-from-a-file
kubectl create secret generic mysecret --from-literal=username=joris --from-literal=mysecret=hurdur --dry-run -o yaml | kubectl apply -f -

# NOTE: See how the secret on the mounted volume is updated while the env var is not.
```


### Removal
```bash
kubectl delete secret apikey
```

## APIs
While it's possible to access the k8s API directly, it's usually more easy to use ```kubectl proxy``` to take care of the certificates, authentication and selecting the right endpoint.

Instructions on how to do it without ```kubectl proxy```: https://kubernetes.io/docs/tasks/administer-cluster/access-cluster-api/#without-kubectl-proxy

```bash
# Easiest is to have kubectl proxy API requests
kubectl proxy --port=8080 # Run this in separate terminal
# Query API:
curl http://localhost:8080/api/

# Alternatively, get API server details:
kubectl config view

# Or, just use kubectl do to queries and get json
kubectl get --raw=/api/v1
kubectl get --raw=/api/v1/pods

# Swagger UI (trailing slash matters!):
# Doesn't seem to work well? UI is unresponsive
http://localhost:8080/swagger-ui/

```


## Misc Notes

```shell
# K8s config
kubectl config view
```

