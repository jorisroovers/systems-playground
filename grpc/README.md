
# GRPC
Back-end and Front-end are in different programming languages on purpose (Python and Go) to showcase that gRPC makes it easy to communicate between programs in different languages using protobuff (similar to how JSON, msgpack, XML or Thrift do that).

Note that the ```Dockerfile```s present in the code are not used in this playground, but rather in that of kubernetes. The docker files need to be adjacent to the source-code, as Docker doesn't allow you to copy files into the docker images that are outside the current build context (also not by symlinking to them).

## Protobuf file
First, you need to write a protocol buffer file, see book_store.proto.

```bash
# Installation of protobuf on mac
brew install protobuf # protobuf also comes as part of anaconda
```

After this, we can generate the client and server code based on the protocol definition

## Go

Client generation:

```bash
# get the protoc-gen-go tool
go get -u github.com/golang/protobuf/protoc-gen-go
# Make sure the gopath is part of the shell path, so you can execute the protoc-gen-go tool
export PATH="$PATH:$(go env GOPATH)/bin"

# Generate Go stubs for our gRPC service, requires protoc-gen-go to be installed
# Execute this from the directory this README.md is in, this will generate web/book_store/book_store.pb.go:
protoc -I . book_store.proto --go_out=plugins=grpc:web/book_store

# Generic format:
# protoc -I <input_dir> <input_proto_file> --go_out=plugins=grpc:<output_dir>
```

Running the client:
```bash
cd web
go run web.go
```

## Python

Server generation:

```bash
pip install -r backend/requirements.txt

# Generate python stubs for our gRPC service
python -m grpc_tools.protoc -I. --python_out=backend --grpc_python_out=backend book_store.proto

# I believed that the same could be done by calling protoc (which would call the python generator behind the scene),
# but the following command only generated book_store_pb2.py and not book_store_pb2_grpc.py
protoc -I . book_store.proto --python_out=plugins=grpc:backend
```

Running the server:
```bash
export BACKEND_NAME="Joris"
python backend/backend.py
```

## Using it
From another terminal session:

```bash
curl localhost:1234/hello
```