# RethinkDB
- Database that pushes JSON to your apps in realtime. No need for polling.
- Query language is called ReQL.
- When not using the real-time change feed, it looks and feels a bit like MongoDB.
- Written in C++. Client libraries in NodeJS, Ruby, Python, Java, others are community supported.
- Real-time feed scales well: You should be able to open thousands of concurrent active feeds on a single RethinkDB node, and scale to tens or hundreds of thousands of feeds across a RethinkDB cluster.
- Company behind RethinkDB shut down, RethinkDB was open-sourced and now managed by the Linux Foundation, but development seems to have stagnated.

```sh
# Start rethinkDB container
docker run -d -P --name rethink1 rethinkdb
# Go to webinterface to play around, based on exposed port
open "http://$(docker port rethink1 8080)"

# Install ReThinkDB library
npm install rethinkdb@2.3.3

# Determine port, run basics example
export RETHINKDB_PORT=$(docker port rethink1 28015 | awk -F ":" '{print $2}'); echo $RETHINKDB_PORT
node ./1_rethinkdb_basics.js

# In new terminal, run realtime feed script ()
export RETHINKDB_PORT=$(docker port rethink1 28015 | awk -F ":" '{print $2}'); echo $RETHINKDB_PORT
node ./2_realtime_feed.js

# In other terminal:
node ./3_updates.js
# Note that running './3_updates.js' twice doesn't cause the realtime feed to update twice, those changes are idempotent.

```

## Notes

- Examples/tutorials on the RethinkDB website use traditional callbacks which makes it hard to put examples consecutively in a single file. The examples were slightly modified to use async/await to work around this.
- Not all operations in ReThinkDB are idempotent. I.e. ```tableDrop``` with a non-existing table or ```tableCreate``` with a pre-existing table will raise exceptions. Those operations also don't seem to have optional parameters to make them idempotent.
