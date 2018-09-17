FROM golang:1.10-alpine as build

# args
ARG version="1.0.0"
ARG repo="github.com/productionwentdown/email-collector"

# dependencies
RUN apk add --no-cache ca-certificates

# source
WORKDIR $GOPATH/src/${repo}
COPY . .

# build
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
RUN go build -ldflags "-s -w" -o /email-collector


FROM scratch

ARG version

# labels
LABEL org.label-schema.vcs-url="https://github.com/productionwentdown/email-collector"
LABEL org.label-schema.version=${version}
LABEL org.label-schema.schema-version="1.0"

# copy binary and ca certs
COPY --from=build /email-collector /email-collector
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

EXPOSE 8080

ENTRYPOINT ["/email-collector"]
