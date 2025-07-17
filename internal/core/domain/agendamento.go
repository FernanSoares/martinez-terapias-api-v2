package domain

import (
	"time"

	"github.com/google/uuid"
)

// StatusAgendamento define os possíveis estados de um agendamento
type StatusAgendamento string

const (
	StatusAgendado  StatusAgendamento = "agendado"
	StatusConfirmado StatusAgendamento = "confirmado"
	StatusRealizado  StatusAgendamento = "realizado"
	StatusCancelado  StatusAgendamento = "cancelado"
)

// Agendamento representa a entidade de agendamento de sessões no domínio
type Agendamento struct {
	ID              uuid.UUID        `json:"id"`
	DataHora        time.Time        `json:"data_hora"`
	ClienteID       uuid.UUID        `json:"cliente_id"`
	ServicoID       uuid.UUID        `json:"servico_id"`
	MassoterapeutaID uuid.UUID        `json:"massoterapeuta_id"`
	ValorCobrado    float64          `json:"valor_cobrado"`
	Status          StatusAgendamento `json:"status"`
	Observacoes     string           `json:"observacoes,omitempty"`
}
