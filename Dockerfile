FROM golang:latest

WORKDIR $GOPATH/src/github.com/D7024E-Distributed-Systems/Kademlia

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build src/main.go
CMD ["go", "run", "src/main.go"]
