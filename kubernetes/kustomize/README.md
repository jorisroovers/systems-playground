# Kustomize

Kustomize allows you to easily generate k8s resource configs. It's sort of a templating language that understands k8s concepts.

```sh
# From the kustomize directory in this repo:
kubectl kustomize simple/

# To apply this config to the k8s cluster
# While you can use kustomize to generate the NS for us, there's an annoyance with prefix/suffix generators, so let's create it manually
# https://github.com/kubernetes-sigs/kustomize/issues/235
kubectl create namespace my-namespace
kubectl kustomize simple/ | kubectl apply -f -

# Clean-up
kubectl delete --all pod,secret,configmap --namespace my-namespace
kubectl delete ns/my-namespace

# More complex example, using overlays (=~ template inheritance)
kubectl kustomize complex/overlays/dev
kubectl kustomize complex/overlays/prod
```