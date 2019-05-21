# Traefik


## Installation
```bash
# Install traefik using helm
helm install stable/traefik --version 1.68.4 --name traefik --namespace kube-system --replace --set dashboard.enabled=true

# When doing this on minikube, you have to patch the installation because minikube doesn't support the LoadBalancer service type.
# So we need to delete the existing 'traefik' service and re-expose it as NodePort type
kubectl delete svc traefik -n kube-system
kubectl delete svc/traefik-dashboard -n kube-system # the separate dashboard service is not needed for minikube
kubectl expose deploy/traefik --type NodePort -n kube-system
kubectl get svc -n kube-system # show services

# To get host URLs, ask minikube for the URLs
# One of these (typically last one) will be the dashboard URL
minikube service -n kube-system traefik --url

# To clean up
helm delete traefik
```