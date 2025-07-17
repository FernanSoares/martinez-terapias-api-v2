package domain

import (
	"github.com/google/uuid"
)

// Servico representa a entidade de serviço (tratamento) no domínio
type Servico struct {
	ID             uuid.UUID `json:"id"`
	NomeServico    string    `json:"nome_servico"`
	Descricao      string    `json:"descricao"`
	DuracaoMinutos int       `json:"duracao_minutos"`
	Valor          float64   `json:"valor"`
	ImagemURL      string    `json:"imagem_url,omitempty"`
}
