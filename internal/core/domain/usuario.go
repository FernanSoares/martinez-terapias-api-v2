package domain

import (
	"time"

	"github.com/google/uuid"
)

// TipoPerfil define os possíveis perfis de usuários do sistema
type TipoPerfil string

const (
	PerfilCliente       TipoPerfil = "cliente"
	PerfilMassoterapeuta TipoPerfil = "massoterapeuta"
	PerfilAdmin         TipoPerfil = "admin"
)

// Usuario representa um usuário do sistema com credenciais de acesso
type Usuario struct {
	ID           uuid.UUID  `json:"id"`
	NomeCompleto string     `json:"nome_completo"`
	Email        string     `json:"email"`
	SenhaHash    string     `json:"senha_hash,omitempty"` // Não incluir no JSON de resposta
	Telefone     string     `json:"telefone"`
	Perfil       TipoPerfil `json:"perfil"`
	Ativo        bool       `json:"ativo"`
	DataCadastro time.Time  `json:"data_cadastro"`
}

// CredenciaisLogin representa os dados necessários para login
type CredenciaisLogin struct {
	Email    string `json:"email"`
	Password string `json:"senha"`
}

// TokenResponse representa a resposta da API após um login bem-sucedido
type TokenResponse struct {
	Token   string `json:"token"`
	Usuario struct {
		ID    uuid.UUID  `json:"id"`
		Nome  string     `json:"nome"`
		Perfil TipoPerfil `json:"perfil"`
	} `json:"usuario"`
}
