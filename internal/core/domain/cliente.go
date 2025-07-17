package domain

import (
	"time"

	"github.com/google/uuid"
)

// Cliente representa a entidade cliente no domínio
type Cliente struct {
	ID            uuid.UUID `json:"id"`
	NomeCompleto  string    `json:"nome_completo"`
	CPF           string    `json:"cpf,omitempty"`
	Email         string    `json:"email"`
	Telefone      string    `json:"telefone"`
	DataNascimento time.Time `json:"data_nascimento"`
	DataCadastro  time.Time `json:"data_cadastro"`
	Ativo         bool      `json:"ativo"`
}
