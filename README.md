C:\Users\chico\go\bin\mockery.exe --all
### Running tests locally
```bash
go test -v ./...
```

### Create mocks from interface
```bash
~/go/bin/mockery --all
```

## Running locally
We have some workers that are located in the path `/workers/{choose your worker}`.

To run it you should go to the desired path and run `go run main.go`.

All the workers are deppended on some AWS features and to run them locally you need to have docker, aws-cli, aws sam and node installed.

We use a node dependency called LocalStack to ease the local development by emulating the aws features locally.

Use aws SAM to run local lambda


### Initialize localstack through docker-compose
```bash
docker-compose up -d
```
---
If everything is setup correctly, you should be able to receive the status of the SQS running locally here: http://localhost:4566/.

Than, you can use the aws-cli locally like this:

```bash
# List existing queues
aws --endpoint-url=http://localhost:4566 sqs list-queues

# Receive message
aws --endpoint-url=http://localhost:4566 sqs receive-message --queue-url http://localhost:4566/000000000000/pix_response
aws --endpoint-url=http://localhost:4566 sqs receive-message --queue-url http://localhost:4566/000000000000/pix_qr_code_request

# Create queue
aws --endpoint-url=http://localhost:4566 sqs create-queue --queue-name pix_response
aws --endpoint-url=http://localhost:4566 sqs create-queue --queue-name pix_qr_code_request
aws --endpoint-url=http://localhost:4566 sqs create-queue --queue-name pix_pagarme_postback
aws --endpoint-url=http://localhost:4566 sqs create-queue --queue-name pix_expired_qr_code

# Send message
aws --endpoint-url=http://localhost:4566 sqs send-message --queue-url http://localhost:4566/000000000000/pix_qr_code_request --message-body file://mocks/json/qrcode_request.json
aws --endpoint-url=http://localhost:4566 sqs send-message --queue-url http://localhost:4566/000000000000/pix_pagarme_postback --message-body file://mocks/json/pagarme_postback.json
aws --endpoint-url=http://localhost:4566 sqs send-message --queue-url http://localhost:4566/000000000000/pix_qr_code_request --message-body file://mocks/json/qrcode_request_with_customer_id.json

aws --endpoint-url=http://localhost:4566 sqs send-message --queue-url http://localhost:4566/000000000000/pix_expired_qr_code --message-body file://mocks/json/qrcode_expired_message.json
```

## Localstack - In case of No connection:

Remove the folder .localstack on root of project

## DynamoDB LocalStack Commands:

# List Tables
aws dynamodb list-tables --endpoint-url http://localhost:4566

# Describe Table
aws dynamodb describe-table --table-name check_expired_qr_code
aws dynamodb describe-table --table-name check_expired_qr_code | grep ID

# Create Table
aws dynamodb create-table \
    --table-name check_expired_qr_code \
    --attribute-definitions \
        AttributeName=id,AttributeType=S \
        AttributeName=pix_expiration_at,AttributeType=N \
    --key-schema \
        AttributeName=id,KeyType=HASH \
        AttributeName=pix_expiration_at,KeyType=RANGE \
--provisioned-throughput \
        ReadCapacityUnits=10,WriteCapacityUnits=5 \
--endpoint-url http://localhost:4566
