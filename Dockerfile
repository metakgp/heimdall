FROM golang:1.21

# Copy metaploy configuration
COPY metaploy/heimdall.metaploy.conf /
COPY metaploy/postinstall.sh /

# Copy source files
WORKDIR /tmp/heimdall
COPY go.mod go.sum main.go /tmp/heimdall

# Build go package
RUN go build

# Copy the binary to /
RUN cp /tmp/heimdall/heimdall /

# Run postinstall script and the binary
CMD ["/postinstall.sh", "/heimdall"]