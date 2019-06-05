# Traefik
Cloud-Native services router (reverse proxy/LB/SSL offload, etc), written in Go.

## Getting Started
```bash
# Install traefik using helm
# Options:
# serviceType=NodePort             : easy use with minikube by exposing the traefik ingress through a NodePort
# dashboard.enabled=true           : enable UI dashboard
# dashboard.serviceType=NodePort   : expose Dashboard through a NodePort as well
# metrics.prometheus.enabled=true  : enable prometheus metrics
# accessLogs.enabled=true          : enable access logs
helm install stable/traefik --version 1.68.4 --name traefik --namespace kube-system --replace --set serviceType=NodePort,dashboard.enabled=true,dashboard.serviceType=NodePort,metrics.prometheus.enabled=true,accessLogs.enabled=true

# Show all treafik related k8s resources
kubectl get all -n kube-system -l app=traefik

# To get host URLs, ask minikube for the URLs
minikube service -n kube-system traefik --url
minikube service -n kube-system traefik-dashboard --url # Dashboard is useful to follow what's going on in browser

# Create simple pod and service, query service
kubectl apply -f app.yaml
curl $(minikube service my-app --url)

# In separate terminal, monitor traefik pod
TRAEFIK_POD=$(kubectl get pod -n kube-system | awk '/traefik/{print $1}'); echo $TRAEFIK_POD
kubectl logs -n kube-system -f $TRAEFIK_POD

# Create traefik ingress
kubectl apply -f ingress.yaml

# Modify /etc/host to add hello-world.info to the minikube ip
# Note: this name mapping is required, you can't curl directly to the IP, Traefik will not serve the page in that case
grep $(minikube ip) /etc/hosts || echo "$(minikube ip) mytraefik.localhost" | sudo tee -a /etc/hosts

# Determine NodePort for Traefik's port 80, so we know which port to query on our host system
TRAEFIK_HTTP_NODE_PORT=$(kubectl get svc traefik -n kube-system -o=jsonpath='{.spec.ports[?(@.port==80)].nodePort}'); echo $TRAEFIK_HTTP_NODE_PORT

# Access app through Traefik ingress
curl "mytraefik.localhost:$TRAEFIK_HTTP_NODE_PORT"

# To clean up
kubectl delete all,ingress -l app=my-app # for some reason 'ingress' is not part of 'all'?
helm delete traefik
sudo sed -i.bak "/$(minikube ip)/d" /etc/hosts # cleanup /etc/hosts
```

## API
Traefik is Cloud Native, so it has APIs to do everything:

```sh
TRAEFIK_DASHBOARD_URL=$(minikube service -n kube-system traefik-dashboard --url); echo $TRAEFIK_DASHBOARD_URL
curl -s "$TRAEFIK_DASHBOARD_URL/health" | jq

# Some other API endpoints
curl -s "$TRAEFIK_DASHBOARD_URL/api" | jq
curl -s "$TRAEFIK_DASHBOARD_URL/api/providers" | jq
curl -s "$TRAEFIK_DASHBOARD_URL/api/providers/kubernetes" | jq

# If prometheus metrics are enabled, access them like so
curl $TRAEFIK_DASHBOARD_URL/metrics
```


### Extra Notes

By default, the serviceType of traefik and the traefik dashboard are set to ```LoadBalancer``` and ```ClusterIP``` respectively. This doesn't work on minikube. By using the helm options ```serviceType=NodePort``` and ```dashboard.serviceType=NodePort``, that get's fixed. To fix it manually:

```sh
# When doing this on minikube, you have to patch the installation because minikube doesn't support the LoadBalancer service type.
# So we need to delete the existing 'traefik' service and re-expose it as NodePort type
kubectl delete svc traefik -n kube-system
kubectl delete svc/traefik-dashboard -n kube-system # the separate dashboard service is not needed for minikube
kubectl expose deploy/traefik --type NodePort -n kube-system
kubectl get svc -n kube-system # show services
```