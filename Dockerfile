FROM golang:latest

WORKDIR $GOPATH/src/github.com/D7024E-Distributed-Systems/Kademlia
RUN apt update
RUN apt install reptyr
RUN yes yes | apt install screen
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build src/main.go
CMD  ["./main"]
# CMD ["go", "run", "src/main.go"]
