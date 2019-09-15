# Build Geth in a stock Go builder container
FROM golang:1.12-alpine as builder

RUN apk add --no-cache musl-dev linux-headers

RUN mkdir -p /go/src/github.com/DaoCloud/dap-manager
ADD . /go/src/github.com/pinlan/training-server
RUN cd /go/src/github.com/pinlan/training-server/cmd/training-server && go build -v

# Pull Geth into a second stage deploy alpine container
FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=builder /go/src/github.com/pinlan/training-server/cmd/training-server/training-server  /
CMD ["/bin/sh","-c","/training-server --kubeconfig /etc/config/kubeconfig.yaml --template_path=/etc/config/pod.template.yaml"]
