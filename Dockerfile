FROM golang:latest as builder


WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build -o mail ./cmd/grpc-mail/main.go

FROM scratch 

WORKDIR /app

COPY --from=builder /app/mail ./
COPY --from=builder /app/.env ./

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 3003
EXPOSE 50053

ENTRYPOINT [ "./mail" ]



