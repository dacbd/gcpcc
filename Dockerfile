FROM golang:1.17 AS builder

WORKDIR /opt/app

COPY go.* ./
RUN go mod download
COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gcpcc .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /opt/app/gcpcc /gcpcc
CMD ["/gcpcc"]
