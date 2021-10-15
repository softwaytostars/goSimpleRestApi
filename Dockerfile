FROM golang:alpine AS build

WORKDIR /app

COPY go.* ./
RUN go mod download

# copy source files and build the binary
COPY . .
RUN go build -o /docker-app-test

FROM alpine:latest

WORKDIR /

COPY --from=build /docker-app-test /docker-app-test
COPY --from=build /app/*.yml ./

ENTRYPOINT ["/docker-app-test"]

