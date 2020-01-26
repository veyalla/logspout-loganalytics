FROM alpine as builder

ARG logspout_version=3.2.6

RUN set -e \
    && apk add --update go build-base git mercurial ca-certificates

COPY ./CGI-WebGateway2.pem /usr/local/share/ca-certificates/CGI-WebGateway2.pem
RUN update-ca-certificates

RUN git clone https://github.com/gliderlabs/logspout /src \
    && cd /src \
    && git checkout tags/v${logspout_version}

WORKDIR /src

COPY ./modules.go ./

RUN go get github.com/kbeaugrand/logspout-loganalytics \
    && go build -ldflags "-X main.Version=${logspout_version}" -o /bin/logspout

FROM scratch

ENTRYPOINT ["/bin/logspout"]

COPY --from=builder /bin/logspout /bin/logspout

# backwards compatibility
# RUN ln -fs /tmp/docker.sock /var/run/docker.sock