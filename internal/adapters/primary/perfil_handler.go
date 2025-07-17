package primary

import (
	"net/http"

	"martinezterapias/api/internal/core/ports"
	"martinezterapias/api/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// PerfilHandler implementa os handlers HTTP para acesso às informações do próprio usuário
type PerfilHandler struct {
	usuarioService ports.UsuarioService
}

// NewPerfilHandler cria uma nova instância de PerfilHandler
func NewPerfilHandler(usuarioService ports.UsuarioService) *PerfilHandler {
	return &PerfilHandler{
		usuarioService: usuarioService,
	}
}

// GetPerfil retorna o perfil do usuário logado
// @Summary Obter perfil do usuário
// @Description Retorna os dados do usuário logado
// @Tags me
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} domain.Usuario
// @Failure 401 {object} map[string]string "Não autorizado"
// @Failure 500 {object} map[string]string
// @Router /me/perfil [get]
func (h *PerfilHandler) GetPerfil(c *gin.Context) {
	// Obter ID do usuário do token JWT
	usuarioID, exists := c.Get(middleware.UsuarioIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}
	
	clienteID, ok := usuarioID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar ID do usuário"})
		return
	}

	// Buscar perfil
	perfil, err := h.usuarioService.GetPerfilCliente(clienteID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, perfil)
}

// GetAgendamentos retorna os agendamentos do usuário logado
// @Summary Obter agendamentos do usuário
// @Description Retorna a lista de agendamentos do usuário logado
// @Tags me
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} domain.Agendamento
// @Failure 401 {object} map[string]string "Não autorizado"
// @Failure 500 {object} map[string]string
// @Router /me/agendamentos [get]
func (h *PerfilHandler) GetAgendamentos(c *gin.Context) {
	// Obter ID do usuário do token JWT
	usuarioID, exists := c.Get(middleware.UsuarioIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}
	
	clienteID, ok := usuarioID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar ID do usuário"})
		return
	}

	// Buscar agendamentos - passando o mesmo ID como solicitante e cliente
	agendamentos, err := h.usuarioService.GetAgendamentosCliente(clienteID, clienteID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, agendamentos)
}
