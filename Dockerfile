FROM golang:1.16.7

COPY . /var/www
WORKDIR /var/www
RUN make build-binaries

ENTRYPOINT ["/var/www/release/v2-linux-amd64", "-config=/var/www/config.docker.yaml"]

EXPOSE 8999
