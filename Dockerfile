# Stage 1: Build da aplicação
FROM golang:1.24-alpine AS builder

# Instalar dependências necessárias para compilação com CGO
RUN apk add --no-cache gcc musl-dev

# Definir o diretório de trabalho
WORKDIR /src

# Copiar arquivos de dependências primeiro para aproveitar o cache
COPY go.mod go.sum ./
RUN go mod download

# Copiar o restante do código fonte
COPY . .

# Compilar a aplicação com suporte a SQLite (CGO necessário)
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o app ./cmd/api

# Stage 2: Imagem final otimizada
FROM alpine:latest

# Adicionar certificados CA e dependências de runtime necessárias
RUN apk add --no-cache ca-certificates tzdata sqlite libc6-compat

# Criar um usuário não-root para executar a aplicação
RUN adduser -D -h /app appuser
WORKDIR /app

# Copiar o binário compilado do estágio anterior
COPY --from=builder /src/app /app/

# Criar diretório para dados e garantir permissões
RUN mkdir -p /data && chown -R appuser:appuser /data

# Configurações de ambiente (usar valores via variáveis de ambiente em produção)
ENV DATABASE_URL=/data/martinezterapias.db \
    SERVER_PORT=8080 \
    TZ=America/Sao_Paulo \
    GIN_MODE=release

# A variável JWT_SECRET deve ser definida externamente via variável de ambiente

# Mudar para o usuário não-root
USER appuser

# Expor a porta da aplicação
EXPOSE 8080

# Verificação de saúde
HEALTHCHECK --interval=30s --timeout=3s --retries=3 \
  CMD wget -q --spider http://localhost:8080/health || exit 1

# Comando para iniciar a aplicação
CMD ["/app/app"]
