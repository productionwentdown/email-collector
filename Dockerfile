FROM golang:1.10-alpine as go

WORKDIR /go/src/email-collector
COPY . .
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
RUN go build -ldflags '-extldflags "-static"' -o email-collector


FROM scratch

EXPOSE 8080
COPY --from=go /go/src/email-collector/email-collector email-collector

ENTRYPOINT ["/email-collector"]
