FROM golang:1.17.2 AS build-env

WORKDIR /go/src/app
COPY go ./go
COPY go.mod .
COPY main.go .

ENV CGO_ENABLED=0

RUN go get -d -v ./...
RUN go build -o /go/bin/app

FROM gcr.io/distroless/base
COPY --from=build-env /go/bin/app /

EXPOSE 8080

CMD ["/app"]