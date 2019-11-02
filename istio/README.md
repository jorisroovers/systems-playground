
The commands below assume a default project and region is set using `gcloud init`.
Use `gcloud config configurations xxx` commands to modify the existing config.

# GKE Setup
As per https://istio.io/docs/setup/platform-setup/gke/

```sh
# Check that we're pointing to the right cloud and no cluster exist
gcloud config configurations list
gcloud container clusters list

# Create k8s cluster for istio on GKE
gcloud container clusters create istio-cluster --cluster-version latest --num-nodes 4

# Add GKE cluster to ~/.kube/config
gcloud container clusters get-credentials istio-cluster

# Grant admin permissions to the current user so it can create the necessary RBAC rules for Istio
kubectl create clusterrolebinding cluster-admin-binding --clusterrole=cluster-admin --user=$(gcloud config get-value core/account)
```

# Istio installation
As per https://istio.io/docs/setup/install/kubernetes/
```sh
# Download Istio
cd /tmp
curl -L https://git.io/getLatestIstio | ISTIO_VERSION=1.3.3 sh -
cd /tmp/istio-1.3.3
export PATH=$PWD/bin:$PATH

# Double check kubectl is pointing to the GKE k8s cluster
kubectl cluster-info
kubectl config get-contexts

# Apply CRD files to install istio
for i in install/kubernetes/helm/istio-init/files/crd*yaml; do kubectl apply -f $i; done

# Check which CRDs we just installed:
kubectl api-resources --verbs=list -o name | grep istio

# Install istio demo profile (TLS + non-TLS variant)
kubectl apply -f install/kubernetes/istio-demo.yaml
```

## Istio verification

### General verification
```sh
# Show services and pods
kubectl get svc -n istio-system
kubectl get pods -n istio-system
```

### Install BookInfo sample micro-service
From https://istio.io/docs/examples/bookinfo/

```sh
cd /tmp/istio-1.3.3
# Make sure that auto proxy injection is enabled for the namespace we want to use (i.e. default)
kubectl label namespace default istio-injection=enabled
kubectl apply -f samples/bookinfo/platform/kube/bookinfo.yaml

# Confirm everything is running:
kubectl get pod,svc

# Confirm services are working:
# Should return "<title>Simple Bookstore App</title>"
kubectl exec -it $(kubectl get pod -l app=ratings -o jsonpath='{.items[0].metadata.name}') -c ratings -- curl productpage:9080/productpage | grep -o "<title>.*</title>"

# Expose BookInfo micro-service so it's accessible from outside, using an Istio Gateway
kubectl apply -f samples/bookinfo/networking/bookinfo-gateway.yaml
kubectl get gateway,virtualservices

# Determine external Gateway URL
export INGRESS_HOST=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
export INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="http2")].port}')
export SECURE_INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="https")].port}')
export GATEWAY_URL=$INGRESS_HOST:$INGRESS_PORT

# Confirm bookstore service is accessible from outside
# Should return "<title>Simple Bookstore App</title>"
# You can also point your browser to this URL
curl -s http://${GATEWAY_URL}/productpage | grep -o "<title>.*</title>"


# Define Destination Rules so Istio is aware of what version routes of the app exists
kubectl apply -f samples/bookinfo/networking/destination-rule-all.yaml
kubectl get destinationrules -o yaml

```
# Playing around
```sh
# Accessing grafana
kubectl -n istio-system port-forward $(kubectl -n istio-system get pod -l app=grafana -o jsonpath='{.items[0].metadata.name}') 3000:3000 &
# Prometheys
kubectl -n istio-system port-forward $(kubectl -n istio-system get pod -l app=prometheus -o jsonpath='{.items[0].metadata.name}') 9090:9090 &
```

# Cleanup
From https://istio.io/docs/setup/install/kubernetes/#uninstall

```sh
cd /tmp/istio-1.3.3
# Delete demo profile
kubectl delete -f install/kubernetes/istio-demo.yaml

# Delete CRDs
for i in install/kubernetes/helm/istio-init/files/crd*yaml; do kubectl delete -f $i; done

# Delete GKE cluster
gcloud container clusters  delete istio-cluster
```