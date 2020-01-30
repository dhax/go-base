FROM golang:alpine AS builder

# Add source code
ADD ./ /go/src/github.com/dhax/go-base/

RUN cd /go/src/github.com/dhax/go-base && \       
    go build && \
    mv ./go-base /usr/bin/go-base

# Multi-Stage production build
FROM alpine

RUN apk add --update ca-certificates

# Retrieve the binary from the previous stage
COPY --from=builder /usr/bin/go-base /usr/local/bin/go-base

# Set the binary as the entrypoint of the container
CMD ["go-base", "serve"]