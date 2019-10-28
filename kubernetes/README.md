# Kubernetes


# Minikube installation

```sh
# Installation instructions at https://github.com/kubernetes/minikube/releases:
# At the time of writing
brew cask install minikube
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
minikube dashboard # Show Web dashboard

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
List all resource types with ```kubectl api-resources```.

Basic resources:
- Pods: collection of scalable containers. A pod typically maps onto a single micro-service (which might be comprised of multiple containers).
- Deployments: supervisor for pods with fine-grained control over how and when a new pod version is rolled out as well as rolled back to a previous state.
- Services: Way to expose containers to the outside world
- Nodes: Infrastructure where k8s is hosting your containers on
- Secrets: Secure storage for sensitive information for easy consumption in your k8s cluster
- Volumes: shared storage

Great for hands-on examples: http://kubernetesbyexample.com

Under-the-hood processes:
- kubelet: runs on every Node to provide computing
- kube-proxy: runs on every Node to do mappings between VIPs and pods via iptables and IP namespaces.
- kube-scheduler
- kube-controller-manager
- kube-dns
- sidecar


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

# You can create multiple resources using a multi-resource yaml file that contains separator '---' lines
kubectl apply -f multi-resource.yaml
```

### Inspection
```bash
kubectl get deployments
kubectl get pods # a deployment can contain multiple pods
kubectl get services  # services expose ports from deployments/pods for external use

# Pod finding
kubectl get pods -l foo=bar # only pods with label 'foo=bar'
kubectl get all -l foo=bar # all resources with label 'foo=bar'
kubectl get pods -l 'foo in (bar, bazz)' # more complex selector
kubectl get pods --field-selector status.phase=Running
kubectl get pods -o json # output in json. Other options: yaml, wide, name, custom-columns, go-template and more
kubectl get pod -o go-template --template "{{range .items}}{{.metadata.name}} {{end}}"
kubectl get pod --watch #  More like 'tail -f' than like 'watch', it appends changes to the terminal.

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

# Note that when using deployments, the container hostname is automatically mapped to the podname
# This is NOT the case when manually creating a single POD
PODNAME=$(kubectl get pod -l label=my-app -o name | sed  "s/pod\///")
kubectl exec $PODNAME hostname

# Roll out a newer version of the my-app
kubectl apply -f deployments/deployment-v2.yaml
curl $(minikube service my-app --url) # will show 'version 2'

# Show all version of my-app
kubectl rollout history deploy/my-app

# Roll back to a previous version
kubectl rollout undo deploy/my-app --to-revision=1
curl $(minikube service my-app --url) # shows 'version 1' again

# Roll out v3, which has 10 replicas configured. Use 'rollout status' to follow along as k8s brings up new pods and tears down the old ones.
#  The --record flag will just record the command we execute in the rollout history
kubectl apply --record -f deployments/deployment-v3.yaml
kubectl rollout status deploy/my-app
kubectl get pods

# Roll out v4, which has the deployment strategy explicitly set to 'Recreate'. Strategy is 'RollingUpdate' by default.
kubectl apply -f deployments/deployment-v4.yaml
# Show name and version of pods
watch "kubectl get pod -o custom-columns=NAME:.metadata.name,VERSION:.spec.containers[0].env[0]"
# Compare with rollout of V3 with V4. V4 will destroy all pods before recreating the new pods. V3 will do a rolling upgrade.
```

## Scaling (replica sets)
A k8s Deployment under-the-hood contains a Replicaset which is a wrapper around a Pod that allows it to horizontally scale the pod.
While you can manually create the ReplicaSet resource type, there's usually no reason to do so directly.
Instead, when you use a Deployment, k8s will automatically create the replicaset for you.

In the past, k8s used the ReplicationController for the same purpose, but the use of Deployment (+ReplicaSet) is now advised over the ReplicationController. Deployments add additional features to ReplicationControllers such as rollback.

Historical context:
Note that the ReplicationController manipulates Pods directly, it does NOT use the ReplicaSet resource type. ReplicaSet's were created to be used by Deployments, and to encapsulate and seperate the replication behavior from the other Deployment capabilities such as versioning, rollback, etc. In the past, with the ReplicationController, the capabilities of replication and deployment (versioning, rollout, rollback, etc) where much more intertwined within the single ReplicationController resource type.

```sh
# Create deployment, show associated replicaset that k8s created
kubectl apply -f deployments/deployment-v1.yaml
kubectl get rs

# Scale deployment up, to 3 pods
kubectl scale --replicas 3 deploy/my-app
kubectl scale --replicas 2 deploy/my-app # scale down

# Note that if you do a GET request against the hosts that k8s will be doing round-robin LB between the pods
curl $(minikube service my-app --url)
# Example Output, note that the 'Host' part changes when you do consecutive requests:
# This is the index page! (Version: 1.0, Host: my-app-66f5b89f64-jnksd)
```

## Services
A service is an abstraction for pods, providing a stable, so called virtual IP (VIP) address.
While pods may come and go and with it their IP addresses, a service allows clients to reliably
connect to the containers running in the pod using the VIP.

Services can be created for pods, deployments, replicasets or replicaset controllers.

Under the hood, the kube-proxy takes care of the actual mapping of the VIPs. THere's multiple drivers under-the-hood,
like IPTables, IPVS, user-space, etc.

```sh
# List services
kubectl get services # shorthand = svc
# Create simple pod to play with
kubectl apply -f simple-pod.yaml

# Create service via yaml
kube apply -f services/my-service.yaml

# 'expose' is a convenience command to quickly create services for already created resources
kubectl expose pod web-from-file --type=NodePort

# Note that the resources in the selector don't have to be available for a service to be created.
# This is called a 'headless' service
kubectl apply -f services/another-service.yaml
curl $(minikube service another-service --url) # Connection refused
kubectl apply -f complex-pod.yaml
curl $(minikube service another-service --url) # Works!
```

### Services and DNS
We creating a service and DNS is enabled in the cluster, all containers will have DNS access to all other hosts in the cluster.

```sh
# Start backend pod, 2 replicaes
kubectl run backend --image=backend --port=50051 --replicas 2 --image-pull-policy=Never --labels "app=grpc,tier=backend"
# Start web pod, 2 replicas, set env var SERVICE_BACKEND to point to backend
kubectl run web --image=web --port=1234 --replicas 2 --image-pull-policy=Never --labels "app=grpc,tier=web" --env="SERVICE_BACKEND=backend:50051"
# Expose the services, this will automatically make them available via DNS within the cluster
kubectl expose deploy/web deploy/backend --type NodePort --labels "app=grpc"

# See how an nslookup in a 'backend' pod resolves to the backend IP
WEB_POD=$(kubectl get pod -l "app=grpc,tier=web" | awk 'NR==2{print $1 }'); echo $WEB_POD
kubectl exec $WEB_POD -ti -- nslookup backend

# Try curling from backend to web using DNS
BACKEND_POD=$(kubectl get pod -l "app=grpc,tier=backend" | awk 'NR==2{print $1 }'); echo $BACKEND_POD
kubectl exec $BACKEND_POD -ti -- curl web:1234

# Note that you cannot access services outside the current namespace
kubectl get svc -n kube-system # this will list kube-dns
kubectl exec $WEB_POD -ti -- nslookup kube-dns # Does not work

# Hit the /hello url on web, which will do a gRPC call to the backend
# Notice the loadbalancing to the backend on subsequent hits
curl $(minikube service web --url)/hello
```

### Service types
ClusterIP, NodePort, Loadbalancer, ExternalName

TODO


## Volumes
Volumes can be shared across different containers.
Multiple types of containers exist: emptyDir, hostpath, awsElasticBlockStore, cephfs, nfs, etc. There's a long list.
Note: Docker also provides volumes, but not as flexible in terms of underlying drivers

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

## Side-cars
Sidecars are containers **within** a pod that provide auxillary services to the main pod container.
Typical examples: caching, proxying, log aggregation, monitoring agents, etc.

```sh
# First create configmap with Nginx config
kubectl create configmap nginx --from-file=sidecar/nginx.conf

# Create pod with sidecar and expose it
kubectl apply -f sidecar/my-pod.yaml
kubectl expose pod pod-with-sidecar --port 4567 --type=NodePort

# Curl the service. Note that the nginxi "Server" header in the HTTP response
curl -v $(minikube service pod-with-sidecar --url)

# Alternatively, open result of 'minikube service pod-with-sidecar --url' in browser, and watch nginx logs
kubectl logs -f pod-with-sidecar -c web-proxy
kubectl logs -f pod-with-sidecar -c web
```

## Init containers
Init containers allow you to perform some startup actions in a container before a different container starts.

```bash
# Create container
kubectl apply -f init-container.yaml

# Note that the pod stays in 'PodInitializing' for 120 sec as the init container is running
kubectl get pod

# Show logs of main container -> will show it waiting for the init container to finish.
PODNAME=$(kubectl get pod -l app=ic -o name | sed  "s/pod\///")
kubectl logs $PODNAME
```

## Health Checks

```bash
# Create healthcheck
kubectl apply -f healthchecks/simple-pod.yaml


```

## Downward API
Way to expose specific k8s info to pods, without the pods needing to do explicit REST API calls to the k8s API (you probably also don't want to allow pods to access the k8s API directly.)

The available fields are not very well documented. Best I could find is this link:
https://kubernetes.io/docs/tasks/inject-data-application/downward-api-volume-expose-pod-information/#capabilities-of-the-downward-api

```sh
kubectl apply -f downward-api.yaml
# Env vars
kubectl exec downward-api-demo env

# Volumes
kubectl exec downward-api-demo cat /etc/podinfo/annotations
kubectl exec downward-api-demo ls /etc/podinfo

# Benefit of volume: config can be updated
kubectl exec downward-api-demo cat /etc/podinfo/labels
kubectl label pod downward-api-demo foo=bar
kubectl exec downward-api-demo cat /etc/podinfo/labels
```

## Persistent volumes
Shared across the cluster (=across nodes) and across node reboot - i.e. a PV is not backed by locally-attached storage
on a worker node but by networked storage system such as nfs, GCEPersistentDisk, AWSElasticBlockStore, HostPath, CephFS, etc.

Primarily interesting for stateful applications (i.e. with databases).

### Persistent Volume Claims
Used by **end-users** to request a persistent volume from k8s.

```bash
# Create Persistent Volume Claim
kubectl create -f persistent-volumes/persistent-volume-claim-v1.yaml
kubectl get persistentvolumeclaims # shorthand=pvc
# Note how k8s automatically allocated a persistent volume for you
kubectl get persistentvolume # shorthand=pvc

# In minikube's case, the PV type is 'hostPath' by default (which means it's just a directory on the minikube VM)
PV_NAME=$(kubectl get pv | awk '/myclaim/{print $1}'); echo $PV_NAME
PV_MINIKUBE_PATH=$(kubectl describe pv $PV_NAME | awk '/Path:/{print $2}'); echo $PV_MINIKUBE_PATH
kubectl describe pv $PV_NAME
minikube ssh ls $PV_MINIKUBE_PATH # Show contents on PV in minikube

# Create a deployment with a pod that mounts the PVC
kubectl apply -f persistent-volumes/deployment-v1.yaml

# Create a file on the persistent storage
PODNAME=$(kubectl get pod -l label=my-app -o name | sed  "s/pod\///"); echo $PODNAME
kubectl exec $PODNAME touch /tmp/persistent/myfile
kubectl exec $PODNAME ls /tmp/persistent

# Delete deployment and recreate it
kubectl delete deploy/my-app
kubectl apply -f persistent-volumes/deployment-v1.yaml

# Show content of PVC on new pod. Note 'myfile' is still there.
PODNAME=$(kubectl get pod -l label=my-app -o name | sed  "s/pod\///"); echo $PODNAME
kubectl exec $PODNAME ls /tmp/persistent
```

### Persistent Volume
Used by **administrators** to pre-create persistent volumes that PersistentVolumeClaims can be mapped to.

```bash
# Create persistent volume
kubectl create -f persistent-volumes/persistent-volume.yaml
kubectl get pv

# When creating a new PVC with the same storage class as the previously created PV (i.e. my-storage-class),
# the PVC will be mapped to the PV
kubectl create -f persistent-volumes/persistent-volume-claim-v2.yaml
kubectl get pv,pvc
```

## StatefulSet
StatefulSet: resource type used to manage stateful applications. Used instead of a Deployment, which only manages stateless apps.
A Stateful set will provide persistent volumes per pod. StatefulSets provide guarantees about the ordering and uniqueness of Pods. This is useful for when you have applications that require stable, unique network identifiers, persistent storage or deterministic deployment/scaling patterns.

```sh
 # Create statefulset + headless service (=no clusterIP)
kubectl apply -f statefulset.yaml

kubectl get statefulsets # shorthand=svc
kubectl get sts,svc,pod,pvc # show all created resources from statefulset.yaml

# You'll notice that the PVC and PODs following deterministic naming conventions,
# as opposed to Deployments which will generate unique ids for pods, etc.

# When you delete a pod, a pod with the same name will be recreated
kubectl delete pod mystatefulset-1; kubectl get pod

# Compared to Deployment, where upon POD deletion, a new POD is created with a different identifier
kubectl apply -f deployments/deployment-v3.yaml
kubectl get pod
POD=$(kubectl get pod -l "label=my-app" | awk 'NR==2{print $1 }'); echo $POD
kubectl delete pod $POD
kubectl get pod
kubectl delete deploy my-app # clean up


# Data is also persisted between statefulset recreations (to avoid data loss):
# Create file
kubectl exec mystatefulset-1 ls /mydata
kubectl exec mystatefulset-1 touch /mydata/foobar
# Delete Statefulset, note how PersistentVolumeClaims are not deleted automatically
kubectl delete sts mystatefulset
kubectl get pvc

# When recreating StatefulSet, same volumes will be mounted, data is still there
kubectl apply -f statefulset.yaml
kubectl exec mystatefulset-1 ls /mydata
```

## DaemonSet

A DaemonSet ensures that all (or some) Nodes run a copy of a Pod. As nodes are added to the cluster, Pods are added to them. As nodes are removed from the cluster, those Pods are garbage collected. Deleting a DaemonSet will clean up the Pods it created.
Useful for running e.g. storage, log collection or monitoring daemons (cepth, gluster, logstash, fluent, Prometheues NodeExporter, AppD agent, etc) on every node.

## Namespaces
Way to group k8s resources together. Like a tenant/project on other IaaS platforms.

```bash
# Create and list namespaces
kubectl create namespace foobar
kubectl create namespace foobar --dry-run -o yaml # Don't create the NS, just show yaml
kubectl get namespace

# By default, kubectl commands run against the 'default' namespace, but you can always specify the namespace explicitly
kubectl get pods --namespace foobar
kubectl get pods -n foobar # shorthand

# Resources can exist with the same name in different namespaces
kubectl apply -f simple-pod.yaml # default namespace
kubectl apply --namespace foobar -f simple-pod.yaml
```

## Configmap
Easy way to store configuration (i.e. config management).

```sh
# Literal configmap configuration
kubectl create configmap my-config --from-literal=foo.bar=hurdur
kubectl get configmap my-config
kubectl get configmap my-config -o yaml

# From file
kubectl create configmap game-config --from-file=configmap/game.properties

# Consume in pod
kubectl apply -f configmap/configmap-example.yaml
kubectl exec configmap-example env
kubectl exec configmap-example cat /etc/config/game.properties
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

## Nodes
Nodes are the k8s compute hosts on which pods get scheduled.

```sh
# Show available nodes:
kubectl get nodes
kubectl describe nodes
```

## Jobs
Jobs are wrapper constructs for pods running tasks that are expected to terminate. Contrast this with regular pods which are typically managed by replicasets or replicaset controllers and not expected to terminate.
The Job resource will do the pod management for you.

```sh
# Create + Run job
kubectl apply -f jobs/simple-job.yaml
# Jobs cannot be re-run (confirmed), you have to delete and re-create them

# Parallel jobs (run job 200 times with 3 jobs in parallel)
kubectl apply -f jobs/parallel-job.yaml
watch kubectl get job # watch completions go up
kubectl get job --watch # alternative

# Cronjob: run every 1 min
kubectl apply -f jobs/cron-job.yaml
kubectl get cronjob # cronjob definitions
watch kubectl get job # cronjob instantiations
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

## Ingress
Resource type that manages external access to the services in a cluster, typically HTTP.
Ingress can provide load balancing, SSL termination and name-based virtual hosting.

https://kubernetes.io/docs/tasks/access-application-cluster/ingress-minikube/
```sh
# Minikube comes with an NGINX Ingress controller
minikube addons enable ingress

# Needs few mins the first time (to download image), but nginx will show up after a while
kubectl get pods -n kube-system # nginx-ingress-controller-xxx

# Create sample app and service (=2 resources in example-app.yaml)
kubectl apply -f ingress/example-app.yaml
kubectl get svc,pod -l app=my-app

# Sample-app can now be accessed through the service
curl $(minikube service --url my-app)

# In new terminal: open ingress controller logs to see what happens in next steps:
NGINX_POD=$(kubectl get pod -n kube-system | awk '/nginx/{print $1}'); echo $NGINX_POD
kubectl logs -f $NGINX_POD -n kube-system

# Create ingress
kubectl apply -f ingress/example-ingress.yaml

# Modify /etc/host to add hello-world.info to the minikube ip
# Note: this name mapping is required, you can't curl directly to the IP, Nginx will not serve the page in that case
grep $(minikube ip) /etc/hosts || echo "$(minikube ip) hello-world.info" | sudo tee -a /etc/hosts

# This works
curl "hello-world.info/foobar"
# These won't
curl "$(minikube ip)/foobar" # hostname required
curl "hello-world.info/hurdur" # only foobar is mapped

# Clean up
kubectl delete ingress/example-ingress svc/my-app pod/my-app
sudo sed -i.bak "/$(minikube ip)/d" /etc/hosts # cleanup /etc/hosts
```



## Under the hood

```shell
# K8s config
kubectl config view

# See what system pods/deployments k8s is running to run k8s itself
kubectl get pod -n kube-system
kubectl get deployments -n kube-system

# Show config of CoreDNS pod
kubectl describe configmap -n kube-system coredns
```