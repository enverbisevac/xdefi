FROM golang:1.17.5-alpine as builder
RUN apk update && apk add --no-cache bash git ca-certificates tzdata && update-ca-certificates
RUN adduser -D -g '' appuser
WORKDIR /app

# Fetch dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o server ./main.go

FROM scratch
# Import from builder.
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd

# Copy our static executable
COPY --from=builder /app/server /usr/local/bin/server
# Use an unprivileged user.
USER appuser
# Run the binary.
ENTRYPOINT ["/usr/local/bin/server"]