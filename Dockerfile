FROM golang

WORKDIR /go/src/github.com/tlwr/monzo_exporter

COPY $PWD/go.mod go.mod
COPY $PWD/go.sum go.sum

ENV GO111MODULE=on

RUN go mod download

COPY $PWD/*.go ./

RUN go build -o /bin/monzo_exporter

ENTRYPOINT ["/bin/monzo_exporter"]
