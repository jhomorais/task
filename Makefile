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
