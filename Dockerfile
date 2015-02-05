FROM golang:1.4

RUN mkdir -p /home/deployer/gosrc/src/github.com/curt-labs/GoQueue
ADD . /home/deployer/gosrc/src/github.com/curt-labs/GoQueue
WORKDIR /home/deployer/gosrc/src/github.com/curt-labs/GoQueue
RUN export GOPATH=/home/deployer/gosrc && go get
RUN export GOPATH=/home/deployer/gosrc && go build -o GoQueue ./index.go

ENTRYPOINT /home/deployer/gosrc/src/github.com/curt-labs/GoQueue/GoQueue