gen-proto:
	protoc --proto_path=proto proto/*.proto --go_out=pb

gen-rpc:
	protoc --proto_path=proto proto/*.proto --go-grpc_out=pb

clean:
	rm pb/github.com/fgmaia/task/*.go

run:
	go run main.go

test:
	go test -cover -race ./...

compose-up:
	docker-compose up -d

docker-exec:
	docker exec -it task /bin/bash

mockary:
	~/go/bin/mockery --all

create-volume:
	docker volume create --name=mysql_task_data

remove-volume:
	docker volume rm mysql_task_data

prepare_rabbitmq:
	mkdir -p etc/rabbitmq/conf	
	sudo cp rabbitmq.conf ./etc/rabbitmq/conf	
	docker-compose up -d
	docker-compose down

run-read-queue-worker:
	go run cmd/workers/taskqueueworker/main.go

run-grpc-server:
	go run cmd/grpcserver/main.go	
