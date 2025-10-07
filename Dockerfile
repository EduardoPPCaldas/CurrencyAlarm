ARG GO_VERSION=1
FROM golang:${GO_VERSION}-bookworm AS builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .

# Build a static Linux binary for the worker located in cmd/api
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
RUN go build -v -o /run-app ./cmd/api

FROM debian:bookworm-slim

# Install CA certificates so HTTPS certificate verification succeeds
# (some minimal images don't include them by default)
RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy the built binary into the runtime image
COPY --from=builder /run-app /usr/local/bin/run-app
RUN chmod +x /usr/local/bin/run-app

# Keep PATH sane
ENV PATH="/usr/local/bin:${PATH}"

# By default the app will look for a .env file in the working directory (/app).
# It's recommended to mount your .env at runtime or pass env vars with -e
# rather than baking secrets into the image.

ENTRYPOINT ["/usr/local/bin/run-app"]
