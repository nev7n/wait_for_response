# Build the app
FROM --platform=$BUILDPLATFORM golang:1.26.1-trixie AS builder
ARG TARGETOS
ARG TARGETARCH
ARG VERSION
WORKDIR /app
# Enable go modules even inside GOPATH
ENV GO111MODULE=on
# Mark neecosem modules as private
# RUN go env -w GOPRIVATE=github.com/newcosem/*

# Copy `go.mod` for definitions and `go.sum` to invalidate the next layer
# in case of a change in the dependencies
#COPY go.mod go.sum ./
# Download dependencies
#RUN go mod download

COPY ./ ./
RUN CGO_ENABLED=0 go test ./...
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -a -ldflags="-X main.version=${VERSION}" -o wait_for_response /app/main/main.go

# Create a minimal docker container and copy the app into it
FROM debian:13.4-slim
RUN apt update && apt install -y ca-certificates iproute2

USER 65534:65534
WORKDIR /app
COPY --from=builder --chmod=0755 /app/wait_for_response .
COPY --chmod=0755 entrypoint.sh .

ENTRYPOINT ["/app/entrypoint.sh"]
