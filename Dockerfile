FROM golang:1.17.3 AS build-env

WORKDIR /go/src/app

COPY api ./api
COPY app ./app
COPY helpers ./helpers
COPY model ./model
COPY security ./security

COPY README.md .
COPY go.mod .

COPY main.go .

ENV CGO_ENABLED=0
ENV GO111MODULE=on

RUN go get -d -v ./...
RUN go mod tidy
RUN go vet -v
RUN go test -v

RUN go build -o /go/bin/app

# We use /base because of OpenSSL, libSSL and glibc
FROM gcr.io/distroless/base

COPY --from=build-env /go/bin/app /

EXPOSE 8080

CMD ["/app"]