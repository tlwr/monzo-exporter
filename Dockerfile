ARG golang_version
FROM golang:$golang_version

WORKDIR /go/src/github.com/tlwr/monzo-exporter

COPY $PWD/go.mod go.mod
COPY $PWD/go.sum go.sum

ENV GO111MODULE=on

RUN go mod download

COPY $PWD/*.go ./

RUN go build -o /bin/monzo-exporter

ENTRYPOINT ["/bin/monzo-exporter"]
