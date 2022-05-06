# CoreDNS

Loosely following: https://dev.to/robbmanes/running-coredns-as-a-dns-server-in-a-container-1d0

# Server
```sh
# NOTE: ensure you're executing these from this directory!
docker run -d --name coredns --restart=always --volume=$(pwd):/root/ -p 6553:53/udp coredns/coredns -conf /root/Corefile
docker logs coredns

# Cleanup:
docker kill coredns
docker container prune
```

# Client
```sh
# use -p to define non-standard port
dig -p 6553 @0.0.0.0 dns.example.com

# get SOA record
dig -p 6553 @0.0.0.0 example.com

# get TXT record
dig -p 6553 @0.0.0.0 example.com TXT
```