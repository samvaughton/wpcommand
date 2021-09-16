FROM golang:1.16.7

COPY . /go/src/github.com/samvaugton/wpcommand
WORKDIR /go/src/github.com/samvaugton/wpcommand
RUN make build-binaries

ENTRYPOINT ["/go/src/github.com/samvaughton/wpcommand/wpcmd", "--config=config.docker.yaml"]

EXPOSE 80
