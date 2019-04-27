.PHONY: all clean

export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64

BUILD=$(shell git rev-parse --short=7 HEAD)
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.Build=$(BUILD)"

all:
	go build -o $(SERVICENAME) $(LDFLAGS)

clean:
	-rm $(SERVICENAME) \
        && docker rmi "carprks/$(SERVICENAME):$(VERSION)" \
        && docker rmi "carprks/$(SERVICENAME):latest"

docker:
	docker build -t "carprks/$(SERVICENAME):$(VERSION)" \
		--build-arg build=$(BUILD) \
		--build-arg version=$(VERSION) \
		--build-arg serviceName=$(SERVICENAME) \
		--build-arg AWS_DB_REGION=$(AWS_DB_REGION) \
		--build-arg AWS_DB_ENDPOINT=$(AWS_DB_ENDPOINT) \
		--build-arg AWS_DB_TABLE=$(AWS_DB_TABLE) \
		--build-arg AWS_ACCESS_KEY_ID=$(AWS_ACCESS_KEY_ID) \
		--build-arg AWS_SECRET_ACCESS_KEY=$(AWS_SECRET_ACCESS_KEY) \
		--build-arg DATABASE_DYNAMO=$(DATABASE_DYNAMO) \
		-f Dockerfile .

