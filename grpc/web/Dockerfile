FROM alpine:3.8

RUN apk add --no-cache git make musl-dev go curl

# Configure Go
ENV GOROOT /usr/lib/go
ENV GOPATH /go
ENV PATH /go/bin:$PATH

# In go < 1.11, all packages need to be in $GOPATH/src
# Since we import "web/book_store" in web.go, Go will be looking for that package in $GOPATH/src/web/book_store.
# For this to work, we need to copy our source-files to $GOPATH/src/web
WORKDIR /go/src/web
COPY . /go/src/web

# Install all dependencies (= all import statements, including "web/book_store")
# Run 'go get -v' for detailed output
RUN go get

EXPOSE 1234

CMD ["go", "run", "web.go"]