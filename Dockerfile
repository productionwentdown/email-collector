FROM golang:1.10-alpine as go

ARG repo=github.com/productionwentdown/email-collector

WORKDIR /go/src/${repo}
COPY . .
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
RUN go build -ldflags '-extldflags "-static"' -o email-collector


FROM scratch

EXPOSE 8080
COPY --from=go /go/src/${repo}/email-collector email-collector

ENTRYPOINT ["/email-collector"]
