FROM node:9.9.0 AS frontend-builder
RUN mkdir frontend
WORKDIR frontend
COPY frontend/package.json .
COPY frontend/yarn.lock .
RUN yarn install
COPY frontend .
RUN yarn build && ls build

FROM golang:1.10.0 AS go-builder
RUN go get github.com/jteeuwen/go-bindata/...
RUN go get github.com/elazarl/go-bindata-assetfs/...
COPY . /go/src/github.com/netice9/swarm-intelligence
WORKDIR /go/src/github.com/netice9/swarm-intelligence
COPY --from=frontend-builder /frontend/build frontend/build
RUN go generate ./frontend
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/swarm-intelligence .

FROM alpine:3.6
COPY --from=go-builder /go/bin/swarm-intelligence /
CMD ["/swarm-intelligence"]
EXPOSE 8080
