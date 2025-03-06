FROM golang:latest as builder


WORKDIR /app
ENV GOOS=linux
ENV CGO_ENABLED=1
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN  go build -C ./cmd  -o ../bin/server
RUN  go build -C ./client  -o ../bin/client


FROM debian:bookworm
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /app
COPY --from=builder /app/bin ./bin
COPY --from=builder /app/.env .bin/server
EXPOSE 8181
EXPOSE 8182


CMD ["./bin/server"]
