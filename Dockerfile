# Specify the version of Go to use
FROM golang:1.19

ENV GO111MODULE=on \
    CGO_ENABLED=0

WORKDIR /conformance
COPY . .

RUN go mod download && go build -o /conformance/run

ENTRYPOINT ["/conformance/run"]
