FROM golang:1.20 as builder

WORKDIR /app/

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/security-go/

FROM scratch
COPY --from=builder /app/main .
COPY --from=builder /app/config.yaml .

EXPOSE 50051
CMD ["./main"]