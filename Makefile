JOBDATE		?= $(shell date -u +%Y-%m-%dT%H%M%SZ)
GIT_REVISION	= $(shell git rev-parse --short HEAD)
VERSION		?= $(shell git describe --tags --abbrev=0)
POSTGRES_TEST_DSN = postgres://app:password@localhost:5445/app_test?sslmode=disable
MIGRATIONS_PATH = ./migrations

#LDFLAGS		+= -linkmode external -extldflags -static
LDFLAGS		+= -X github.com/samvaughton/wpcommand/version.Version=$(VERSION)
LDFLAGS		+= -X github.com/samvaughton/wpcommand/version.Revision=$(GIT_REVISION)
LDFLAGS		+= -X github.com/samvaughton/wpcommand/version.BuildDate=$(JOBDATE)

build-binaries:
	go get github.com/mitchellh/gox
	gox -verbose -output="release/{{.Dir}}-{{.OS}}-{{.Arch}}" \
		-ldflags "$(LDFLAGS)" -osarch="linux/amd64"

dependency:
	@go get -v ./pkg/./...

test:
	go test ./pkg/./...

test-integration: docker-up dependency setup-test-db
	go test -tags=integration ./pkg/./...

setup-test-db:
	PGPASSWORD=password docker exec wpcommand_postgres_test_1 psql -U app -d postgres -c 'DROP DATABASE IF EXISTS app_test'
	PGPASSWORD=password docker exec wpcommand_postgres_test_1 psql -U app -d postgres -c 'CREATE DATABASE app_test'
	./bin/migrate -database $(POSTGRES_TEST_DSN) -path $(MIGRATIONS_PATH) up
	go run ./test/load_test_fixtures.go

docker-up:
	@docker-compose up -d

docker-down:
	@docker-compose down

build:
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags "$(LDFLAGS) -w -s" .

install:
	# CGO_ENABLED=0 GOOS=linux go install -ldflags "$(LDFLAGS)" github.com/samvaughton/wpcommand/wpcmd
	GOOS=linux go install -ldflags "$(LDFLAGS)" github.com/samvaughton/wpcommand/cmd/wpcmd

push-image:
	docker push samrentivo/wpcommand:$(VERSION)

image:
	docker build -t samrentivo/wpcommand:$(VERSION) -f Dockerfile .

run:
	go install github.com/samrentivo/wpcommand/wpcmd
	wpcmd

dev-ci:
	/home/sam/go/bin/CompileDaemon --exclude-dir=docker --command="./dev.sh"  --build="./build.sh"