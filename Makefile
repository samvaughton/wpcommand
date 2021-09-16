JOBDATE		?= $(shell date -u +%Y-%m-%dT%H%M%SZ)
GIT_REVISION	= $(shell git rev-parse --short HEAD)
VERSION		?= $(shell git describe --tags --abbrev=0)

#LDFLAGS		+= -linkmode external -extldflags -static
LDFLAGS		+= -X github.com/samvaughton/wpcommand/version.Version=$(VERSION)
LDFLAGS		+= -X github.com/samvaughton/wpcommand/version.Revision=$(GIT_REVISION)
LDFLAGS		+= -X github.com/samvaughton/wpcommand/version.BuildDate=$(JOBDATE)

build-binaries:
	go get github.com/mitchellh/gox
	gox -verbose -output="release/wpcmd-{{.Dir}}-{{.OS}}-{{.Arch}}" \
		-ldflags "$(LDFLAGS)" -osarch="linux/amd64"

test:
	go get github.com/mfridman/tparse
	go test -json -v `go list ./... | egrep -v /tests` -cover | tparse -all -smallscreen

build:
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags "$(LDFLAGS) -w -s" .

install:
	# CGO_ENABLED=0 GOOS=linux go install -ldflags "$(LDFLAGS)" github.com/samvaughton/wpcommand/wpcmd
	GOOS=linux go install -ldflags "$(LDFLAGS)" github.com/samvaughton/wpcommand/cmd/wpcmd

image:
	docker build -t samvaughton/wpcommand:$(VERSION) -f Dockerfile .

run:
	go install github.com/samvaughton/wpcommand/wpcmd
	wpcmd