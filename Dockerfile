FROM golang:1.17 AS builder

WORKDIR /opt/app

COPY go.* ./
RUN go mod tidy
COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gcpcc .

FROM scratch
COPY --from=builder /opt/app/gcpcc /gcpcc
CMD ["/gcpcc"]
