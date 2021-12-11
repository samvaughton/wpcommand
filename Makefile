JOBDATE		?= $(shell date -u +%Y-%m-%dT%H%M%SZ)
GIT_REVISION	= $(shell git rev-parse --short HEAD)
VERSION		?= $(shell git describe --tags --abbrev=0)
POSTGRES_TEST_DSN = postgres://app:password@localhost:5445/app_test?sslmode=disable
POSTGRES_TEST_DOCKER_DSN = postgres://app:password@wpcmd_postgres_test:5432/app_test?sslmode=disable
MIGRATIONS_PATH = ./migrations

#LDFLAGS		+= -linkmode external -extldflags -static
LDFLAGS		+= -X github.com/samvaughton/wpcommand/version.Version=$(VERSION)
LDFLAGS		+= -X github.com/samvaughton/wpcommand/version.Revision=$(GIT_REVISION)
LDFLAGS		+= -X github.com/samvaughton/wpcommand/version.BuildDate=$(JOBDATE)

build-binaries:
	go get github.com/mitchellh/gox
	gox -verbose -output="release/{{.Dir}}-{{.OS}}-{{.Arch}}" \
		-ldflags "$(LDFLAGS)" -osarch="linux/amd64"

port-forward-prod:
	kubectl -n wpcommand port-forward pod/wpcommand-db-postgresql-0 5432:5432

dependency:
	@go get -v ./pkg/./...

test-unit: dependency
	go test ./pkg/./...

test-integration: docker-up dependency setup-test-db
	DATABASE_DSN=$(POSTGRES_TEST_DSN) go test -tags=integration ./pkg/./...

setup-test-db:
	docker run --network container:wpcmd_postgres_test postgres:latest psql postgres://app:password@localhost:5432/postgres?sslmode=disable -c 'DROP DATABASE IF EXISTS app_test'
	docker run --network container:wpcmd_postgres_test postgres:latest psql postgres://app:password@localhost:5432/postgres?sslmode=disable -c 'CREATE DATABASE app_test'
	./bin/migrate -database $(POSTGRES_TEST_DSN) -path $(MIGRATIONS_PATH) up
	go run ./test/load_test_fixtures.go --config=config.test.yaml

# This is the command to fully run
run-ci-tests: docker-up-ci
	docker-compose run wpcmd_test

# This is the entry point of the test docker-compose service
test-integration-ci: setup-test-ci-db dependency
	mkdir -p /tmp/test-reports
	DATABASE_DSN=$(POSTGRES_TEST_DOCKER_DSN) gotestsum --junitfile /tmp/test-reports/unit-tests.xml -- -tags=integration ./pkg/./...

setup-test-ci-db:
	./bin/migrate -database $(POSTGRES_TEST_DOCKER_DSN) -path $(MIGRATIONS_PATH) up
	DATABASE_DSN=$(POSTGRES_TEST_DOCKER_DSN) go run ./test/load_test_fixtures.go

docker-up-ci:
	docker-compose --profile test build --build-arg WPCMD_CONFIG=config.docker.yaml
	docker-compose up -d # do not bring up the test container as will run that manually
	sleep 5
	docker exec wpcmd_postgres_test psql postgres://app:password@localhost:5432/postgres?sslmode=disable -c 'DROP DATABASE IF EXISTS app_test'
	docker exec wpcmd_postgres_test psql postgres://app:password@localhost:5432/postgres?sslmode=disable -c 'CREATE DATABASE app_test'

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

build:
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags "$(LDFLAGS) -w -s" .

install:
	# CGO_ENABLED=0 GOOS=linux go install -ldflags "$(LDFLAGS)" github.com/samvaughton/wpcommand/wpcmd
	GOOS=linux go install -ldflags "$(LDFLAGS)" github.com/samvaughton/wpcommand/cmd/wpcmd

push-image:
	docker push samrentivo/wpcommand:$(VERSION)

image:
	docker build -t samrentivo/wpcommand:$(VERSION) -f Dockerfile .

release: image push-image

run:
	go install github.com/samrentivo/wpcommand/wpcmd
	wpcmd

dev-ci:
	/home/sam/go/bin/CompileDaemon --exclude-dir=docker --command="./dev.sh"  --build="./build.sh"