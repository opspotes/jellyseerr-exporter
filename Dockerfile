FROM golang:1.21-alpine3.17 as builder

WORKDIR /src

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/jellyseerr-exporter
RUN update-ca-certificates

# Minimalist image
FROM scratch

COPY --from=builder /go/bin/jellyseerr-exporter /go/bin/jellyseerr-exporter
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ENTRYPOINT ["/go/bin/jellyseerr-exporter"]
