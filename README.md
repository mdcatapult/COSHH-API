## COSHH API

This GOLang based application provides the backend API service to complement the [COSHH webapp](https://gitlab.mdcatapult.io/informatics/coshh/coshh-ui). It uses a Postgres Database to store information about chemicals, their hazardous properties, expiry dates and links to safety docs.

### Running

`make run`

This starts the API and a local Postgres instance.

By default, the app starts with no data. To populate the app, follow the [ETL guide](scripts/etl/README.md). You can also start with an empty db if you prefer. The Postgres docker container uses an [sql script](scripts/init.sql) to create the schema.

### Debugging
You can run it all using the docker-compose file (use within the `make run` task above) or start up your locally changed versions of the components using `docker compose up -d`. This will build your local version of the app and run it via docker.  
You can also run and debug it in your IDE of choice by starting the server using `cmd/main.go`.  
To run your local versions you need to tell the backend what database connection params to use and where the files detailing the different lab and project names are.


```bash
export DBNAME=informatics \
export HOST=localhost \
export PASSWORD=postgres \
export PORT=5432 \
export USER=postgres \
export LABS_CSV="/Users/my.name/IdeaProjects/coshh-api/assets/labs.csv" \
export PROJECTS_CSV="/Users/my.name/IdeaProjects/coshh-api/assets/projects_041022.csv"
``` 

Start the database:
```bash
docker-compose up -d db
``` 

Start the API without docker:
```bash
cd cmd
go run main.go
```
OR

Start the API using docker:
```bash
docker compose up -d server
```

OR use your IDE and run the `cmd/main.go` file. Remember to set your env vars.

### Testing

`make test`

### Accessing the database locally from the command line

Ensure you set the schema, e.g.

```
psql -h localhost -U postgres -d informatics        \\ password is postgres
SET schema 'coshh';                                 
```

### Testing Authenticated Routes
Get the Auth0 client token from the Auth0 web portal. Use curl to auth against the example `protected` route.
```bash
curl --request GET \
  --url http:/localhost:8080/protected \
  --header 'authorization: Bearer INSERT AUTH0 TOKEN'
```  
Successful auth results in `"You have successfully authenticated"`. Failure to auth results in `{"message":"Requires authentication"}`.

### Gotchas

#### SQL

When writing any new sql queries always remember to commit the transaction!

#### CI

There was a glitch in the publish API stage in CI in October 2022 (which has since resolved itself) which meant that in order to deploy the API  the image 
had to be  built locally and pushed up to the registry manually.  In the event this should happen again use this command:

```docker build -t registry.mdcatapult.io/informatics/software-engineering/coshh/api:<tag name> . && docker push registry.mdcatapult.io/informatics/software-engineering/coshh/api:<tag name>```

N.B Mac M1 users may need to build the image for amd64 (as opposed to arm64) with `--platform linux/amd64`
