FROM golang:1.4

RUN mkdir -p /go/src/github.com/curt-labs/GoQueue
ADD . /app

RUN go install github.com/curt-labs/GoQueue

ENTRYPOINT /go/bin/GoQueue
