# Martinez Terapias API

API backend em Go para o sistema de agendamento e gerenciamento de clientes da Martinez Terapias.

## Tecnologias

- Go 1.24
- Gin (framework web)
- GORM (ORM)
- SQLite (banco de dados)
- JWT (autenticação)

## Estrutura do Projeto

A API segue uma arquitetura hexagonal (ports and adapters):

- `cmd/api/`: Ponto de entrada da aplicação
- `internal/`: Código principal da aplicação
  - `adapters/`: Adaptadores de entrada e saída
    - `primary/`: Handlers HTTP
    - `secondary/`: Repositórios de dados
  - `core/`: Lógica de negócios
    - `domain/`: Entidades de domínio
    - `ports/`: Interfaces para adapters
    - `services/`: Serviços da aplicação
  - `config/`: Configurações da aplicação
  - `middleware/`: Middlewares HTTP

## Funcionalidades

- Autenticação JWT
- Gerenciamento de clientes
- Gerenciamento de agendamentos
- Gerenciamento de serviços
- Gerenciamento de fichas de anamnese


## Como executar

### Via Docker

```bash
# Construir a imagem Docker
docker build -t martinez-api .

# Executar o container
docker run -p 8080:8080 -v ./data:/data martinez-api
```

### Localmente

```bash
# Baixar dependências
go mod download

# Executar a aplicação
go run cmd/api/main.go
```

## Variáveis de Ambiente

- `SERVER_PORT`: Porta do servidor (padrão: 8080)
- `DATABASE_URL`: Caminho para o banco de dados SQLite (padrão: martinezterapias.db)
- `JWT_SECRET`: Segredo para assinatura dos tokens JWT


## Endpoints da API

### Autenticação
- `POST /api/registrar`: Registrar novo usuário
- `POST /api/login`: Autenticar usuário
- `POST /api/registrar-admin`: Registrar usuário administrador

### Clientes
- `GET /api/clientes`: Listar todos os clientes
- `GET /api/clientes/:id`: Obter detalhes de um cliente
- `POST /api/clientes`: Criar novo cliente
- `PUT /api/clientes/:id`: Atualizar cliente existente
- `DELETE /api/clientes/:id`: Excluir cliente

### Agendamentos
- `GET /api/agendamentos`: Listar todos os agendamentos
- `GET /api/agendamentos/:id`: Obter detalhes de um agendamento
- `POST /api/agendamentos`: Criar novo agendamento
- `PUT /api/agendamentos/:id`: Atualizar agendamento existente
- `PATCH /api/agendamentos/:id/status`: Atualizar status de um agendamento
- `DELETE /api/agendamentos/:id`: Excluir agendamento
- `POST /api/agendamentos/:id/solicitar-reagendamento`: Solicitar reagendamento

### Serviços
- `GET /api/servicos`: Listar todos os serviços
- `GET /api/servicos/:id`: Obter detalhes de um serviço
- `POST /api/servicos`: Criar novo serviço
- `PUT /api/servicos/:id`: Atualizar serviço existente
- `DELETE /api/servicos/:id`: Excluir serviço

### Anamnese
- `GET /api/anamnese`: Listar todas as fichas de anamnese
- `GET /api/anamnese/:id`: Obter detalhes de uma ficha de anamnese
- `PUT /api/anamnese/:id`: Atualizar ficha de anamnese
- `GET /api/clientes-anamnese/:cliente_id`: Obter ficha de anamnese de um cliente
- `POST /api/clientes-anamnese/:cliente_id`: Criar ficha de anamnese para um cliente

### Perfil do Usuário
- `GET /api/me/perfil`: Obter perfil do usuário autenticado
- `GET /api/me/agendamentos`: Obter agendamentos do usuário autenticado
