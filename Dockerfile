FROM golang:1.7.3
ADD . /go/src/github.com/netice9/swarm-intelligence
WORKDIR /go/src/github.com/netice9/swarm-intelligence
RUN go install .
ENV PORT=5000
CMD swarm-intelligence

