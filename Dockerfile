FROM golang:1.20 as build

WORKDIR /go/src/oslo
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -o /go/bin/oslo ./cmd/sloctl/main.go


FROM gcr.io/distroless/static-debian11

COPY --from=build /go/bin/oslo /
ENTRYPOINT ["/oslo"]
