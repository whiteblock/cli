FROM ubuntu:latest

# tells kernel to not expect any input from the frontend
# this bypasses the need for tzdata nonsense
ENV DEBIAN_FRONTEND noninteractive
ENV GOPATH ${HOME}
ENV GOBIN ${GOPATH}/bin

# installs dependencies
RUN apt-get update && \
    apt-get install -y vim apt-utils iputils-ping expect git git-extras software-properties-common tmux \
    inetutils-tools wget ca-certificates curl build-essential libssl-dev golang-go

# installs solc package for solidity compiler
RUN add-apt-repository ppa:ethereum/ethereum && apt-get update && apt-get install -y solc

# copies directory with CLI source code from host machine to container
ADD . /cli

# sets PWD to appropriate directory and compiles go binaries for CLI application
WORKDIR /cli/whiteblock
RUN go get && go build

# sets default workdirectory to root and configures paths
WORKDIR /
ENV PATH /cli/whiteblock:${GOBIN}:${PATH}

# enters container using bash
ENTRYPOINT ["/bin/bash"]