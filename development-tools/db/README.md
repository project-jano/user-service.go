# Using MongoDB Test database

1. Start MongoDB

```shell
[sudo] docker-compose -f docker-compose.mongodb.yml up
```

2. Create database and assign user to it (using any MongoDB client, like MongoDB Compass)

````
use project-jano
db.createCollection("users")
db.createUser({ user: "jano", pwd: "jano", roles: [ "readWrite"] })
````

3. Use MongoDB Express if you would like to explore the collection

```shell
[sudo] docker-compose -f docker-compose.mongodb-browser.yml up
```

MongoDB Express starts in port 8081 but it's routed to port 9002.