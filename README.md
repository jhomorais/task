# Task project demo

## Purpose
Task project demo

## Dependencies
- Docker
- Docker Compose

## Getting Started

First create the db docker volume:
```bash
docker volume create --name=mysql_task_data
```

Now execute

```bash
make prepare-rabbitmq
```

then

```bash
make compose-up
```

This command will start all containers with docker-compose.

### open 3 consoles and run
```bash
make docker-exec
```

Now we are ready to start the application.

### For each console opened

### Start grpc server in new terminal
```bash
make run-grpc-server
```

### Start queue worker to read messages from rabbitMQ queue
```bash
make run-read-queue-worker
```

### Run client to create task
```bash
go run cmd/grpclient/main.go TASK_SUMMARY
```

## Make commands

### Running tests locally
```bash
make test
```
### Create mocks from interface
```bash
make mock
```

### Gen proto files
```bash
make gen-proto
```

### Gen rpc files
```bash
make gen-rpc
```

## for more options open Makefile archive

## Example

![alt text](https://github.com/fgmaia/task/blob/master/how_to_test_console.png?raw=true)

## Some tests

![alt text](https://github.com/fgmaia/task/blob/master/how_to_test_console1.png?raw=true)