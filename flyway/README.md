# Flyway

IMPORTANT NOTE (see https://stackoverflow.com/a/11618350/381010):

Flyway requires you to specify the schema against which your migration runs.

In MySQL, physically, a schema is synonymous with a database. You can substitute the keyword SCHEMA instead of DATABASE in MySQL SQL syntax, for example using CREATE SCHEMA instead of CREATE DATABASE.

Most other database products draw a distinction. For example, in the Oracle Database product, a schema represents only a part of a database: the tables and other objects owned by a single user.

## DB Server (MariaDB)
```sh
# Start DB server (it can take a min for DB server to come online)
docker run --rm --name some-mariadb -e MYSQL_ROOT_PASSWORD=my-secret-pw -d mariadb:10.4.10
# Create DB
docker exec -ti some-mariadb mysql -u root -pmy-secret-pw -e 'create schema mydb; show schemas;'
# In MYSQL, this is identical to:
docker exec -ti some-mariadb mysql -u root -pmy-secret-pw -e 'create databases mydb; show databases;'

# Show there's no tables before running migrations
docker exec -ti some-mariadb mysql -u root -pmy-secret-pw -e 'use mydb; show tables;'

# Cleanup (useful when playing around)
docker exec -ti some-mariadb mysql -u root -pmy-secret-pw -e 'drop database mydb; show databases;'

```

## Flyway
```sh
# Test connect to DB, print some info
docker run --rm --link some-mariadb -v $(pwd)/migrations:/flyway/sql flyway/flyway:6.2.1 -user=root -password=my-secret-pw -schemas=mydb -url=jdbc:mariadb://some-mariadb/mydb info

# Run migrations
docker run --rm --link some-mariadb -v $(pwd)/migrations:/flyway/sql flyway/flyway:6.2.1 -user=root -password=my-secret-pw -url=jdbc:mariadb://some-mariadb -schemas=mydb migrate

# Show data
docker exec -ti some-mariadb mysql -u root -pmy-secret-pw -e 'SELECT * FROM mydb.PERSON;'

# Pro version gives 'undo'
```