FROM alpine:3.10.3
WORKDIR /tmp

RUN apk add maven openjdk10-jdk

ADD zkexample /tmp/zkexample

WORKDIR /tmp/zkexample

# 'mvn dependency:resolve' is unnecessary, it would be executed as part of 'mvn package', but by doing this
# step seperately, we already download the packages as a separate step in docker so we don't have to do that again
# in case the next 'mvn package' step fails.
RUN mvn dependency:resolve
RUN mvn package
