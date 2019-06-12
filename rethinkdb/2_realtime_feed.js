const r = require("rethinkdb")
console.log("PORT:", process.env.RETHINKDB_PORT);

(async function run() {
    console.log("Connecting to DB...")
    let connection = await r.connect({ host: 'localhost', port: parseInt(process.env.RETHINKDB_PORT) })
    console.log("DONE\n");

    console.log("Logging all changes below...")
    r.table('authors').changes().run(connection, function (err, cursor) {
        if (err) throw err;
        cursor.each(function (err, row) {
            if (err) throw err;
            console.log(JSON.stringify(row, null, 2), "\n");
        });
    });

})();