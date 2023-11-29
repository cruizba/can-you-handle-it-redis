FROM golang:1.21.4 AS builder

ARG TARGETOS
ARG TARGETPLATFORM
ARG TARGETARCH
RUN echo building execute for "$TARGETPLATFORM"

WORKDIR /workspace

# Cache dependencies
RUN go env -w GOCACHE=/go-cache
RUN go env -w GOMODCACHE=/gomod-cache

COPY ./go.* ./
RUN --mount=type=cache,target=/gomod-cache \
  go mod download

COPY . .

RUN --mount=type=cache,target=/gomod-cache --mount=type=cache,target=/go-cache \
    GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0 \
    go build

FROM alpine

COPY --from=builder /workspace/redis-load-test /app/

ENTRYPOINT ["/app/redis-load-test"]