# IPVS

This README describes how to setup an IPVS VM that acts as loadbalancer to 2 nginx docker containers.
Based of https://medium.com/@benmeier_/a-quick-minimal-ipvs-load-balancer-demo-d5cc42d0deb4

## Target hosts

Setup 2 containers (A and B) running nginx, serving very simple index.html files.
```sh
# Run from this directory
docker run --rm -d -v "$(pwd)/servers/A:/usr/share/nginx/html" -p 8001:80  --name nginx-A nginx
docker run --rm -d -v "$(pwd)/servers/B:/usr/share/nginx/html" -p 8002:80  --name nginx-B nginx

# From host system:
curl localhost:8001
curl localhost:8002
```

## Loadbalancer
We need a VM to run the IPVS loadbalancer because IPVS requires the `ip_vs` kernel module to be enabled. This is not possible to with containers without modifying the host system (as containers run on the same kernel) - using a VM is the easy route.

```sh
vagrant up && vagrant ssh

sudo yum install -y ipvsadm
# Check if IPVS is enabled on the machine
sudo ipvsadmin -l
sudo modinfo ip_vs # required kernel module

# The nginx docker containers are accessible on the host system, which is the default gateway in Vagrant
# Use some shell magic to extract default gw (thanks https://stackoverflow.com/a/31549164/381010)
export HOST_IP=$(netstat -rn | grep "^0.0.0.0 " | tr -s ' ' | cut -d " " -f2)
curl $HOST_IP:8001
curl $HOST_IP:8002

# Set Guest IP, needed for next step
export GUEST_IP=$(ifconfig eth0 | awk '/inet /{print $2}')

# Setup IPVS loadbalancer on port 8000, loadbalacing to the 2 containers
# -A: add a new LB
# -t: service IP+port
# -s rr : scheduler = RoundRobin (rr)
sudo ipvsadm -A -t $GUEST_IP:8000 -s rr
sudo ipvsadm -a -t $GUEST_IP:8000 -r $HOST_IP:8001 -m
sudo ipvsadm -a -t $GUEST_IP:8000 -r $HOST_IP:8002 -m
sudo ipvsadm -l -n # -n = show numeric ports
# Note: this does NOT work with 0.0.0.0 as GUEST_IP, tried it...

# Do some requests, notice load balancing!
for i in {1..10}; do curl $GUEST_IP:8000; done

# Show stats
sudo ipvsadm -L -n --stats --rate

# Change load balancing mechanism to weighted round-robin, add more weight to webserver B
sudo ipvsadm -E -t $GUEST_IP:8000 -s wrr
sudo ipvsadm -e -t $GUEST_IP:8000 -r $HOST_IP:8002 -m -w 3

# Do some requests, notice the difference
for i in {1..10}; do curl $GUEST_IP:8000; done

# Cleanup
sudo ipvsadm -D -t $GUEST_IP:8000
```
