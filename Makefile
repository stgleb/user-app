ifndef PKGS
PKGS := $(shell go list ./... 2>&1 | grep -v 'vendor' | grep -v 'sanity')
endif

all: get-tools vendor-sync lint vet

compose:
	docker-compose up

docker-build:
	docker build -t stgleb/user-app -f Dockerfile.server .
	docker build -t stgleb/mysql -f Dockerfile.database .

get-tools:
	go get -u golang.org/x/lint/golint

lint:
	for file in $(GO_FILES); do \
		golint $${file}; \
		if [ -n "$$(golint $${file})" ]; then \
			exit 1; \
		fi; \
	done

proto:
	protoc --go_out=plugins=grpc:api/ api/api.proto

vet:
	go vet $(PKGS)

vendor-sync:
	go mod tidy
	go mod download
	go mod vendor