FROM golang:1.16.7

ARG WPCMD_CONFIG="config.k8.yaml"

COPY . /var/www
COPY ${WPCMD_CONFIG} /var/www/config.build.yaml

WORKDIR /var/www
RUN make build-binaries

RUN go get gotest.tools/gotestsum

EXPOSE 8999

ENTRYPOINT ["/var/www/release/v2-linux-amd64", "-config=/var/www/config.build.yaml"]

