# AWX
Based of https://github.com/geerlingguy/awx-container/

```sh
docker-compose up -d
```

- Wait a while for installation to complete (wait until DB migrations are complete).
- Then log in with username/password: `admin/password`

# Changes made to docker-compose.yml

1. Change image (see https://github.com/geerlingguy/awx-container/issues/46):
```sh
awx_web:
image: "geerlingguy/awx_web:latest"
#image: "ansible/awx_web:latest"
# ....
awx_task:
image: "geerlingguy/awx_task:latest"
# image: "ansible/awx_task:latest"
```

2. Change awx_web port to       - "8001:8052"

