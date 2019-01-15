FROM golang

COPY $PWD/* /go/src/github.com/tlwr/monzo_exporter/

WORKDIR /go/src/github.com/tlwr/monzo_exporter

RUN go get && go build -o /bin/monzo_exporter

ENTRYPOINT ["/bin/monzo_exporter"]
