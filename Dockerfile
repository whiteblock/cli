FROM golang:1.11.5-alpine as built

ENV GOPATH ${HOME}
ENV GOBIN ${GOPATH}/bin


# copies directory with CLI source code from host machine to container
ADD . /cli
# sets PWD to appropriate directory and compiles go binaries for CLI application
WORKDIR /cli/whiteblock
RUN apk add git
RUN go get && go build

FROM alpine:latest as final

COPY --from=built /cli/whiteblock/whiteblock /cli/whiteblock/whiteblock
COPY --from=built /cli/etc/ /cli/etc
RUN  ln -s /cli/whiteblock/whiteblock /usr/local/bin/whiteblock

# sets default workdirectory to root and configures paths
WORKDIR /

# enters container using bash
ENTRYPOINT ["/bin/bash"]
