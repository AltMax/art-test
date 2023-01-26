MOCKERY_VERSION=v2.16.0

all: build

build: build_server build_migrations

build_server:
	GOPRIVATE=github.com/AltMax CGO_ENABLED=0 go build -o unit_service

build_migrations:
	GOPRIVATE=github.com/AltMax CGO_ENABLED=0 go build -o migrate_common ./migrations/common/*.go

build_docker:
	docker build -t unit_service --platform=linux/amd64 . 

test:
	GOPRIVATE=github.com/AltMax go test -cover -race -p 1 ./...

generate_pb:
	mkdir -p services
	rm -rf services/*.pb.go
	docker run --rm -v $(PWD):$(PWD) -w $(PWD) --platform=linux/amd64 protogen -I=proto --gogofaster_out=plugins=grpc:services `ls proto`

mod:
	GOPRIVATE=github.com/AltMax go mod tidy

generate_mocks:
	docker run --env-file .go-env --rm -it -v $(PWD):$(PWD) -w $(PWD) --platform=linux/amd64 vektra/mockery:$(MOCKERY_VERSION) --output=./units/mocks --name=Units --srcpkg=github.com/AltMax/art-test/units