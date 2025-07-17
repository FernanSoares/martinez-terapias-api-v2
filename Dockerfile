FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download


COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -a -o app ./cmd/api

FROM alpine:latest

RUN apk add --no-cache libc6-compat sqlite

WORKDIR /root/


COPY --from=builder /app/app .


RUN mkdir -p /data

# Configurações de ambiente (substitua por valores reais em produção)
ENV DATABASE_URL=/data/martinezterapias.db
ENV SERVER_PORT=8080
ENV JWT_SECRET=seu_segredo_jwt_aqui

# Expor a porta 8080
EXPOSE 8080

# Comando para iniciar a aplicação
CMD ["./app"]
