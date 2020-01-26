FROM golang:alpine as builder

ARG logspout_version=3.2.6

# Create appuser.
ENV USER=logspout
ENV UID=10001 

RUN adduser \    
    --disabled-password \    
    --gecos "" \    
    --home "/nonexistent" \    
    --shell "/sbin/nologin" \    
    --no-create-home \    
    --uid "${UID}" \    
    "${USER}"

RUN apk update \
    && apk add --no-cache go build-base git ca-certificates \
    && update-ca-certificates

RUN git clone https://github.com/gliderlabs/logspout /src \
    && cd /src \
    && git checkout tags/v${logspout_version}

WORKDIR /src

COPY ./modules.go ./

RUN go get github.com/kbeaugrand/logspout-loganalytics \
    && go build -ldflags "-X main.Version=${logspout_version} -w -s" -o /bin/logspout

FROM scratch
LABEL maintainer "Kevin BEAUGRAND <kbeaugrand@gmail.com>"

VOLUME /mnt/routes

# Import from the builder.
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy our static executable.
COPY --from=builder /bin/logspout /bin/logspout

# Use an unprivileged user.
USER logspout:logspout

ENTRYPOINT ["/bin/logspout"]