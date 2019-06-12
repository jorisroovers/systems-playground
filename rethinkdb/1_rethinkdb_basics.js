const r = require("rethinkdb")
console.log("PORT:", process.env.RETHINKDB_PORT);

(async function run() {
    console.log("Connecting to DB...")
    let connection = await r.connect({ host: 'localhost', port: parseInt(process.env.RETHINKDB_PORT) })
    console.log("DONE\n");

    console.log("Cleaning up from previous runs");
    console.log("\tDropping table 'authors'")
    try {
        await r.db('test').tableDrop("authors").run(connection)
    } catch (e) { console.log("\t", e.msg) }
    console.log("DONE\n")

    // Create 'authors' table in DB 'test'.
    console.log("Creating 'authors' table...")
    try {
        let result = await r.db('test').tableCreate('authors').run(connection);
        console.log(JSON.stringify(result, null, 2));
    } catch (e) { console.log("\t", e.msg) }
    console.log("DONE\n")


    // insert some data
    console.log("Inserting test data in 'authors' table...")
    try {
        let result = await r.table('authors').insert([
            {
                name: "William Adama", tv_show: "Battlestar Galactica",
                posts: [
                    { title: "Decommissioning speech", content: "The Cylon War is long over..." },
                    { title: "We are at war", content: "Moments ago, this ship received word..." },
                    { title: "The new Earth", content: "The discoveries of the past few days..." }
                ]
            },
            {
                name: "Laura Roslin", tv_show: "Battlestar Galactica",
                posts: [
                    { title: "The oath of office", content: "I, Laura Roslin, ..." },
                    { title: "They look like us", content: "The Cylons have the ability..." }
                ]
            },
            {
                name: "Jean-Luc Picard", tv_show: "Star Trek TNG",
                posts: [
                    { title: "Civil rights", content: "There are some words I've known since..." }
                ]
            }
        ]).run(connection)
        console.log(JSON.stringify(result, null, 2));
    } catch (e) { console.log("\t", e.msg) }
    console.log("DONE\n")


    // Find data
    try {
        console.log("Finding all data...");
        let cursor = await r.table('authors').run(connection)
        let result = await cursor.toArray();
        console.log(JSON.stringify(result, null, 2));
        console.log("DONE\n")

        console.log("Finding specific data...");
        cursor = await r.table('authors').filter(r.row('name').eq("William Adama")).run(connection)
        result = await cursor.toArray();
        console.log(JSON.stringify(result, null, 2));

        console.log("Getting specific row by id '%s'...", result[0].id);
        result = await r.table('authors').get(result[0].id).run(connection)
        console.log(JSON.stringify(result, null, 2));


    } catch (e) { console.log("\t", e.msg) }


    console.log("Closing connection...")
    await connection.close()
    console.log("DONE")

})();