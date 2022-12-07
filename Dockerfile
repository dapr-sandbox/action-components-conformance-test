# Specify the version of Go to use
FROM golang:1.19

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /conformance
COPY go.mod .
COPY go.sum .
COPY conformance_test.go .
COPY main.go .

RUN go mod download

RUN go build -o /conformance/run

ENTRYPOINT ["/conformance/run"]
