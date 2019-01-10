FROM ubuntu:latest

RUN apt-get update && \
    apt-get install -y vim apt-utils iputils-ping expect git git-extras software-properties-common tmux \
    inetutils-tools wget ca-certificates curl build-essential libssl-dev golang-go 

RUN add-apt-repository ppa:ethereum/ethereum && apt-get update && apt-get install -y solc

ADD . /cli

WORKDIR /cli/whiteblock
RUN go get && go build

WORKDIR /

ENTRYPOINT ["/cli/whiteblock/whiteblock"]