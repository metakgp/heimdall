FROM golang:1.21

# Copy metaploy configuration
COPY metaploy/heimdall.metaploy.conf /
COPY metaploy/postinstall.sh /

# Copy source files
WORKDIR /
COPY go.mod go.sum main.go mail.go /

# Build go package
RUN go build

# Run postinstall script and the binary
CMD ["/postinstall.sh", "/heimdall"]