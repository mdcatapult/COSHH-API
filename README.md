### Running

`make run`

You may need to `npm i` if running for the first time, and `npm install -g @angular/cli` if you don't have Angular.

By default, the app starts with no data. To populate the app, follow the [ETL guide](etl/README.md).

### Debugging
You can run it all using the docker-compose file (use within the `make run` task above) or start up your locally changed versions of the components. 
To run your local versions you  need to tell the backend what database connection params to use and where the file detailing the different lab names is.

```bash
export DBNAME=informatics \
export HOST=localhost \
export PASSWORD=postgres \
export PORT=5432 \
export USER=postgres \
export LABS_CSV="/Users/my.name/IdeaProjects/coshh/api/labs.csv" \
export PROJECTS_CSV="/Users/my.name/IdeaProjects/coshh/api/projects_041022.csv"
``` 

Start the database:
```bash
docker-compose up db
``` 

Start the API:
```bash
cd api
go run main.go
```

Start the UI:
```bash
cd ui
ng s
```

Go to `http://localhost:4000` in your browser.

### Testing

`make test`

### Accessing the database locally from the command line

Ensure you set the schema, e.g.

```
psql -h localhost -U postgres -d informatics        \\ password is postgres
SET schema 'coshh';                                 
```

### Angular Environment Variables

Note that Angular cannot read env vars at run time.  They are injected in via the `.environment.ts` and `environment.prod.ts`
files at build time.  *If you want to change the API URL you will need to build a new image after updating the environment file.* 

### Gotchas

When writing any new sql queries always remember to commit the transaction!

### CI and wopr deployment

The app is deployed on https://coshh.wopr.inf.mdc/

There was a glitch in the publish API stage in CI in October 2022 (which has since resolved itself) which meant that in order to deploy the API  the image 
had to be  built locally and pushed up to the registry manually.  In the event this should happen again use this command:

```docker build -t registry.mdcatapult.io/informatics/software-engineering/coshh/api:<tag name> . && docker push registry.mdcatapult.io/informatics/software-engineering/coshh/api:<tag name>```

N.B Mac M1 users will need to build the image for amd64 (as opposed to arm64) with `--platform linux/amd64`  

Once you have checked the image has successfully been pushed to the registry (see [here](https://gitlab.mdcatapult.io/informatics/software-engineering/coshh/container_registry/249)) the `coshh-api` workload can be updated to use the new image.
