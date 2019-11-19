# Zookeeper
```sh
# Start zookeeper server instance
docker run --name some-zookeeper --restart always -d zookeeper:3.5.6

# Start zkCli.sh client in docker
docker run -it --rm --link some-zookeeper:zookeeper zookeeper:3.5.6 zkCli.sh -server zookeeper
```

## Zookeeper shell (zkCli.sh)
See https://github.com/apache/zookeeper/blob/a6c36b69cc72d7d67e392dab5360007d6f737bef/zookeeper-docs/src/main/resources/markdown/zookeeperCLI.md

```sh
# Basics
help # use 'help' to list all commands
quit # self-explanatory, note that data persists across client sessions since it's stored on the server :-)
ls / # List all top-level nodes

# Create 'zk_test' zookeeper node, get and set some data
create /zk_test my_data
get /zk_test
set /zk_test foobar
get /zk_test
stat /zk_test # stats

# You can create child nodes and list them
create /zk_test/child1 child1data
create /zk_test/child2 child2data
create /zk_test/child3 child2data
get /zk_test/child1
ls /zk_test
ls -R /zk_test # recursive
ls -s /zk_test # include stats

# Cleanup
delete /zk_test/child2 # specific node
deleteall /zk_tests # delete 'directory'


# Create ephemeral node: it will disappear once the client disconnects (or times out)
# This is useful to implement registor/service discovery/quorums in zookeeper: every client creates an ephemeral
# node for itself with the zookeeper service. As long as the node is there we know the node is alive.
create /zk_test/ephemeral1 ephemeral1data
create /zk_test/ephemeral2 ephemeral2data
# Try quitting and rejoining, ephemeral data will be gone

# Watchers: watch for change to a path
get -w /zk_test
# In separate terminal (notice how first terminal prints NodeDataChanged message)
set /zk_test newvalue123
removewatches /zk_test # remove watchers
```

# Zookeeper Java client

IMPORTANT: Make sure the server is running :-)

```sh
# Make sure you're in the current dir, then build and run docker image
docker build -t zk-example .
docker run -ti --link some-zookeeper:zookeeper zk-example sh

# Inside docker container
# Set classpath (there's better ways to do this using maven, but that's not our focus here :-) )
CLASSPATH=$(find /root/.m2/repository/ -iname *.jar | tr '\n' ":")
# Run simple zookeeper example
java -cp "$CLASSPATH:/tmp/zkexample/target/zkexample-1.0-SNAPSHOT.jar" com.mycompany.zkexample.ZkConnect
```
