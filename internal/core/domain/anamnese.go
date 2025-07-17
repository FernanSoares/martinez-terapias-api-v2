package domain

import (
	"time"

	"github.com/google/uuid"
)

// Anamnese representa a ficha de histórico de saúde do cliente
type Anamnese struct {
	ID                      uuid.UUID `json:"id"`
	ClienteID               uuid.UUID `json:"cliente_id"`
	DataPreenchimento       time.Time `json:"data_preenchimento"`
	QueixaPrincipal         string    `json:"queixa_principal"`
	HistoricoDoencas        string    `json:"historico_doencas"`
	CirurgiasPrevias        string    `json:"cirurgias_previas"`
	MedicamentosEmUso       string    `json:"medicamentos_em_uso"`
	Alergias                string    `json:"alergias"`
	HabitosDiarios          string    `json:"habitos_diarios"`
	ObservacoesMassoterapeuta string    `json:"observacoes_massoterapeuta"`
}
