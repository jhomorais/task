#TODO FINALIZAR

# Task project demo test

## Purpose
Task project demo test

## Dependencies
- Docker
- Docker Compose

## Getting Started

First create the db docker volume:

`docker volume create --name=mysql_task_data`

If want to clean volume run and create it again

`docker volume rm mysql_task_data`

Now execute

`make prepare-rabbitmq`

and

`docker-compose up -d`

This command will start all containers with docker-compose.

Now we are ready to start the application.

`docker exec -it task /bin/bash`

To start all workers:

`go run main.go`

And start the server with:

`go run cmd/server/main.go`

## Make commands

### Running tests locally
```bash
make test
```
### Create mocks from interface
```bash
make mock
```