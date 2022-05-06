# Bind

TODO: continue
See here: https://linuxtut.com/en/d1fc46fd636f79025ab0/

```sh
# NOTE: ensure you're executing these from this directory!
docker run -d --name=bind9 -p 6553:53/udp --volume=$(pwd):/etc/bind/ --volume  /tmp:/var/cache/bind internetsystemsconsortium/bind9:9.18
docker logs bind9

# Cleanup:
docker kill bind9
docker container prune
```

# Client
```sh
# Recursive lookup just works by default
dig -p 6553 google.com @0.0.0.0
```