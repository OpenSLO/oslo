FROM golang:1.24-alpine3.20 AS builder

WORKDIR /app

COPY ./go.mod ./go.sum ./
COPY ./cmd/oslo ./cmd/oslo
COPY ./internal ./internal

ARG LDFLAGS

RUN CGO_ENABLED=0 go build \
  -ldflags "${LDFLAGS}" \
  -o /artifacts/oslo \
  "${PWD}/cmd/oslo"


FROM gcr.io/distroless/static-debian12

COPY --from=builder /artifacts/oslo /usr/bin/oslo

ENTRYPOINT ["oslo"]
