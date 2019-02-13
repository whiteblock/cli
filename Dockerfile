FROM golang:1.11.5-stretch as built

ENV GOPATH ${HOME}
ENV GOBIN ${GOPATH}/bin


# copies directory with CLI source code from host machine to container
ADD . /cli
# sets PWD to appropriate directory and compiles go binaries for CLI application
WORKDIR /cli/whiteblock
RUN go get && go build

FROM ubuntu:latest as final

COPY --from=built /cli/whiteblock/whiteblock /cli/whiteblock/whiteblock
COPY --from=built /cli/etc/ /cli/etc
RUN  ln -s /cli/whiteblock/whiteblock /usr/local/bin/whiteblock
# tells kernel to not expect any input from the frontend
# this bypasses the need for tzdata nonsense
#ENV DEBIAN_FRONTEND noninteractive
#ENV GOPATH ${HOME}
#ENV GOBIN ${GOPATH}/bin

# installs dependencies
#RUN apt-get update && \
#    apt-get install -y vim iputils-ping expect git git-extras software-properties-common tmux \
#    inetutils-tools wget ca-certificates curl build-essential libssl-dev 

# installs solc package for solidity compiler
#RUN add-apt-repository ppa:ethereum/ethereum && apt-get update && apt-get install -y solc

# sets default workdirectory to root and configures paths
WORKDIR /

# enters container using bash
ENTRYPOINT ["/bin/bash"]
