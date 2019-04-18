#!/usr/bin/env bash
docker-compose stop db

docker-clean

docker ps -a --format '{{.Names}} {{.Status}}' | grep 'Exited' | awk '{print $1}' | xargs docker rm

docker-compose up -d db

aws dynamodb create-table --table-name identity --attribute-definitions AttributeName=identifier,AttributeType=S --key-schema AttributeName=identifier,KeyType=HASH --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 --endpoint-url http://docker.devel:8000

#aws dynamodb put-item --table-name identity --endpoint-url http://docker.devel:8000 --item '{"company":{"BOOL":false},"email":{"S":"bob@bob.bob"},"identifier":{"S":"bob"},"phone":{"S":"111"},"registrations":{"L":[{"M":{"oversized":{"BOOL":false},"plate":{"S":"tester"},"vehicleType":{"S":"bike"}}}]}}'