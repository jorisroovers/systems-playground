# NETCONF-YANG

```sh
virtualenv .venv && source .venv/bin/activate
# ncclient = populat NETCONF Client
pip install ncclient==0.6.6

## SERVER
# NetOpeer = linux-based NETCONF Server
# Sysrepo = linux-based YANG datastore
docker pull sysrepo/sysrepo-netopeer2

# NOTE: v0.7.7 seems to have SSH issues, v0.7.6 and 'latest' work fine.
docker run -it --name sysrepo -p 830:830 --rm sysrepo/sysrepo-netopeer2:latest

# SSH into container, you'll notice netconf 'hello' response
# password = "netconf"
ssh netconf@localhost -p 830 -s netconf

# Alternatively, run netconf.py which uses ncclient
python netconf.py
```