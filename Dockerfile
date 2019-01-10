FROM golang:1.11.3-stretch

ADD . /cli

RUN cd /cli/whiteblock &&\
    go get || \
    go build

WORKDIR /cli/whiteblock/

ENTRYPOINT ["/cli/whiteblock/whiteblock"]