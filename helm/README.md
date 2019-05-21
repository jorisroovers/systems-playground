# Helm
Client-side: helm
Server-side: Tiller

Assumption: minikube is installed. See kubernetes directory for instructions if not.

## Installation

```sh
# Installing helm on mac
brew install kubernetes-helm

# Make sure minikube is running and the shell is initialized so kubectl points to minikube
minikube start
eval $(minikube docker-env)

# Init helm (will use kubectl under the hood and point to minikube)
helm init
kubectl get pod -n kube-system # Should show the Tiller pod
kubectl get deploy -n kube-system # Should show Tiller deploy
# Note, you can just delete the Tiller deploy and re-run 'helm init' to reinstall it (I tried it, it works!)

```

## Usage
```sh
helm repo update # update repos

# Install a helm Chart (find charts at https://hub.helm.sh)
helm install stable/mysql
# Once installed, a chart is referred to as a 'release'.
# Install specific version and name the release ourselves (vs. auto-generated)
helm install stable/consul --version 3.6.2 --name myconsul
# --replace to replace a release with same name
helm install stable/consul --version 3.6.2 --name myconsul --replace

# Inspect installed release
helm ls
helm status myconsul
kubectl get pod # Show pods

# Get details for a chart (doesn't need to be installed)
helm inspect stable/mysql

# Uninstall
helm delete myconsul
```