FROM google/golang
 
WORKDIR /gopath/src/github.com/curt-labs/GoQueue
ADD . /gopath/src/github.com/curt-labs/GoQueue/
 
# go get all of the dependencies
RUN go get

RUN go install

CMD []
ENTRYPOINT ["/gopath/bin/GoQueue"]