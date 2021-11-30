# Simple test application

This is a simple application that exposes a simple API for creating, deleting, retrieving document.

One could consider adding endpoints for bulk creation or deletion for a list of documents.
Then we should add a patch endpoint with a payload that would contain a list of objects telling if the state 
should be deleted or created, and the document itself.

The configuration file is config.yml.
By default, the application runs with an in memory data source.
However, one can run docker-compose in order to run a mongo db server. The mongo datastore will then be used by the app (thanks to an environment variable set in docker-compose.yml)


## Generate swagger documentation

Thanks to the swaggo module, one can annotate the endpoints and generate swagger documentation with:

`swag init  --output docs/resourcedocument --parseInternal --parseDependency`

## Run application in local
`go run main.go`

## Run application in docker
`docker build -t docker-app-test .`

run `docker container run --rm -p 8040:8040  docker-app-test` for using an in memory data source

or run `docker-compose up` for using mongodb as datastore
(you may need to create a docker network with `docker network create goapi_default --subnet 172.24.29.0/29` for example)

## Run tests
go test -v `go list ./... | grep -v goconvey` for running all tests except the ones in package goconvey.

Indeed, they are behavior tests that need to be run against the server api.

To run the behavior test, launch docker-compose , then

`go test -v ./goconvey`

## CI/CD

The pipeline is run under gitlabci and the .gitlabci.yml is its configuration file.
See gitlab/README.md

## Access to swagger api documentation

![Rest API documentation](./docs/images/swaggerapi.png)

The documentation is accessible at : http://localhost:8040/swagger/index.html

### Some call examples:

##### Add documents
`curl -X PUT --include http://localhost:8040/documents/toto  --header "Content-Type: application/json" --data '{"name":"monnom", "description":"mydesc"}'`
`curl -X PUT --include http://localhost:8040/documents/titi  --header "Content-Type: application/json" --data '{"name":"monnom2", "description":"mydesc2"}'`

##### Get all documents
`curl --include http://localhost:8040/documents`

##### Get a document given id
`curl --include http://localhost:8040/documents/toto`

##### Delete a document given id
`curl -X DELETE --include http://localhost:8040/documents/toto`

##### Post messages to kafka broker
`curl -X POST http://localhost:8040/emails`
