FROM golang:1.12.5-stretch as built


# copies directory with CLI source code from host machine to container
ADD . /go/src/github.com/whiteblock/cli
# sets PWD to appropriate directory and compiles go binaries for CLI application
WORKDIR /go/src/github.com/whiteblock/cli/whiteblock
RUN sed -i "s/DEFAULT_VERSION/compiled $(date) commit-$(git rev-parse HEAD)/g" cmd/version.go
RUN go get
#RUN go build -ldflags "-linkmode external -extldflags -static" -a .
RUN go build
FROM ubuntu:latest as final

COPY --from=built /go/src/github.com/whiteblock/cli/whiteblock/whiteblock /cli/whiteblock/whiteblock
COPY --from=built /go/src/github.com/whiteblock/cli/etc/ /cli/etc
COPY --from=built /go/src/github.com/whiteblock/cli/etc/ /etc
RUN  ln -s /cli/whiteblock/whiteblock /usr/local/bin/whiteblock
# tells kernel to not expect any input from the frontend
# this bypasses the need for tzdata nonsense

# sets default workdirectory to root and configures paths
WORKDIR /

# enters container using bash
ENTRYPOINT ["/bin/bash"]
