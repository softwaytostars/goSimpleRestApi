# Application for test

This is a simple application that exposes a simple API for creating, deleting, retrieving document.

One could consider adding endpoints for bulk creation or deletion for a list of documents.
Then we should add a patch endpoint with a payload that would contain a list of objects telling if the state 
should be deleted or created, and the document itself.

## Run tests
go test -v ./...

## Run application in local
`go run main.go`

## Run application in docker
`docker build -t docker-app-test .`

`docker container run --rm -p 8040:8040  docker-app-test`

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
