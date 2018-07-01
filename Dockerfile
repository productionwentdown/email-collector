ARG name="email-collector"
ARG repo="github.com/productionwentdown/${name}"


FROM golang:1.10-alpine as go

ARG name
ARG repo

WORKDIR /go/src/${repo}
COPY . .
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
RUN go build -ldflags '-extldflags "-static"' -o ${name}


FROM alpine:3.7

ARG name
ARG repo

EXPOSE 8080
COPY --from=go /go/src/${repo}/${name} /${name}

ENTRYPOINT ["/email-collector"]
