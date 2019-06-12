const r = require("rethinkdb")
console.log("PORT:", process.env.RETHINKDB_PORT);

(async function run() {
    console.log("Connecting to DB...")
    let connection = await r.connect({ host: 'localhost', port: parseInt(process.env.RETHINKDB_PORT) })
    console.log("DONE\n");

    console.log("Updating all database records to add 'type' field...")
    try {
        // Add 'type' field to all records
        let result = await r.table('authors').update({ type: "fictional" }).run(connection)
        console.log(JSON.stringify(result, null, 2));
    } catch (e) { console.log("\t", e.msg) }
    console.log("DONE\n");


    console.log("Updating specific field for specific record...")
    try {
        let result = await r.table('authors').filter(r.row("name").eq("William Adama")).update({ rank: "Admiral" }).run(connection)
        console.log(JSON.stringify(result, null, 2));

    } catch (e) { console.log("\t", e.msg) }
    console.log("DONE\n");

    console.log("Closing connection...")
    await connection.close()
    console.log("DONE")

})();