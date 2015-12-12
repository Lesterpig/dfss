### DFSS - MGDB lib ###

This library is used in order to manage a connection to mongoDB
It uses the mgo driver, but aims at simplifying the queries to the database

## Mongo Manager ##

The struct handling the connection is MongoManager. It requires the environment variable MONGOHQ_URL in order to initialize a connection : it is the uri containing the informations to connect to mongo.
For example, in a test environment, we may have :
MONGOHQ_URL=localhost (require a mongo instance running on default port 27017)
In a prod environment however, it will more likely be :
MONGOHQ_URL=adm1n:AStr0ngPassw0rd@10.0.4.4:27017

## Mongo Collection ##

Please refer to the example to see the API in practice.
