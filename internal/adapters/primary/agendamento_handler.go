package primary

import (
	"net/http"
	"time"

	"martinezterapias/api/internal/core/domain"
	"martinezterapias/api/internal/core/ports"
	"martinezterapias/api/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AgendamentoHandler implementa os handlers HTTP para agendamentos
type AgendamentoHandler struct {
	agendamentoService ports.AgendamentoService
}

// NewAgendamentoHandler cria uma nova instância de AgendamentoHandler
func NewAgendamentoHandler(agendamentoService ports.AgendamentoService) *AgendamentoHandler {
	return &AgendamentoHandler{
		agendamentoService: agendamentoService,
	}
}

// GetAll retorna todos os agendamentos, com possibilidade de filtragem
// @Summary Listar agendamentos
// @Description Retorna a lista de agendamentos, com opção de filtros
// @Tags agendamentos
// @Produce json
// @Param data_inicio query string false "Data de início (formato: YYYY-MM-DD)"
// @Param data_fim query string false "Data de fim (formato: YYYY-MM-DD)"
// @Param status query string false "Status do agendamento (agendado, confirmado, realizado, cancelado)"
// @Security ApiKeyAuth
// @Success 200 {array} domain.Agendamento
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string "Não autorizado"
// @Failure 500 {object} map[string]string
// @Router /agendamentos [get]
func (h *AgendamentoHandler) GetAll(c *gin.Context) {
	// Verificar parâmetros de consulta
	dataInicioStr := c.Query("data_inicio")
	dataFimStr := c.Query("data_fim")
	status := c.Query("status")

	var agendamentos []domain.Agendamento
	var err error

	// Filtrar por período
	if dataInicioStr != "" && dataFimStr != "" {
		dataInicio, errInicio := time.Parse("2006-01-02", dataInicioStr)
		dataFim, errFim := time.Parse("2006-01-02", dataFimStr)
		
		if errInicio != nil || errFim != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de data inválido. Use YYYY-MM-DD"})
			return
		}
		
		// Definir a hora final para o fim do dia
		dataFim = dataFim.Add(23 * time.Hour).Add(59 * time.Minute).Add(59 * time.Second)
		
		agendamentos, err = h.agendamentoService.GetByPeriodo(dataInicio, dataFim)
	} else if status != "" {
		// Validar status
		var statusAgendamento domain.StatusAgendamento
		switch status {
		case "agendado":
			statusAgendamento = domain.StatusAgendado
		case "confirmado":
			statusAgendamento = domain.StatusConfirmado
		case "realizado":
			statusAgendamento = domain.StatusRealizado
		case "cancelado":
			statusAgendamento = domain.StatusCancelado
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Status inválido"})
			return
		}
		
		agendamentos, err = h.agendamentoService.GetByStatus(statusAgendamento)
	} else {
		agendamentos, err = h.agendamentoService.GetAll()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, agendamentos)
}

// GetByID retorna um agendamento pelo ID
// @Summary Obter agendamento por ID
// @Description Retorna os detalhes de um agendamento específico
// @Tags agendamentos
// @Produce json
// @Param id path string true "ID do agendamento"
// @Security ApiKeyAuth
// @Success 200 {object} domain.Agendamento
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string "Não autorizado"
// @Failure 404 {object} map[string]string
// @Router /agendamentos/{id} [get]
func (h *AgendamentoHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	agendamento, err := h.agendamentoService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Agendamento não encontrado"})
		return
	}

	c.JSON(http.StatusOK, agendamento)
}

// Create cria um novo agendamento
// @Summary Criar agendamento
// @Description Cria um novo agendamento de sessão
// @Tags agendamentos
// @Accept json
// @Produce json
// @Param agendamento body domain.Agendamento true "Dados do agendamento"
// @Security ApiKeyAuth
// @Success 201 {object} domain.Agendamento
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string "Não autorizado"
// @Failure 500 {object} map[string]string
// @Router /agendamentos [post]
func (h *AgendamentoHandler) Create(c *gin.Context) {
	var agendamento domain.Agendamento
	if err := c.ShouldBindJSON(&agendamento); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validações básicas
	if agendamento.ClienteID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do cliente é obrigatório"})
		return
	}

	if agendamento.ServicoID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do serviço é obrigatório"})
		return
	}

	if agendamento.MassoterapeutaID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do massoterapeuta é obrigatório"})
		return
	}

	if agendamento.DataHora.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data e hora são obrigatórias"})
		return
	}

	if err := h.agendamentoService.Create(&agendamento); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, agendamento)
}

// Update atualiza um agendamento existente
// @Summary Atualizar agendamento
// @Description Atualiza os dados de um agendamento existente
// @Tags agendamentos
// @Accept json
// @Produce json
// @Param id path string true "ID do agendamento"
// @Param agendamento body domain.Agendamento true "Dados do agendamento"
// @Security ApiKeyAuth
// @Success 200 {object} domain.Agendamento
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string "Não autorizado"
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /agendamentos/{id} [put]
func (h *AgendamentoHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Verificar se o agendamento existe
	_, err = h.agendamentoService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Agendamento não encontrado"})
		return
	}

	var agendamento domain.Agendamento
	if err := c.ShouldBindJSON(&agendamento); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Garantir que o ID na URL corresponde ao ID no corpo da requisição
	agendamento.ID = id

	if err := h.agendamentoService.Update(&agendamento); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, agendamento)
}

// UpdateStatus atualiza apenas o status de um agendamento
// @Summary Atualizar status do agendamento
// @Description Atualiza apenas o status de um agendamento existente
// @Tags agendamentos
// @Accept json
// @Produce json
// @Param id path string true "ID do agendamento"
// @Param status body map[string]string true "Novo status"
// @Security ApiKeyAuth
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string "Não autorizado"
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /agendamentos/{id}/status [patch]
func (h *AgendamentoHandler) UpdateStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var request struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validar o status
	var statusAgendamento domain.StatusAgendamento
	switch request.Status {
	case "agendado":
		statusAgendamento = domain.StatusAgendado
	case "confirmado":
		statusAgendamento = domain.StatusConfirmado
	case "realizado":
		statusAgendamento = domain.StatusRealizado
	case "cancelado":
		statusAgendamento = domain.StatusCancelado
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status inválido"})
		return
	}

	if err := h.agendamentoService.UpdateStatus(id, statusAgendamento); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status atualizado com sucesso"})
}

// Delete remove um agendamento
// @Summary Excluir agendamento
// @Description Remove um agendamento
// @Tags agendamentos
// @Produce json
// @Param id path string true "ID do agendamento"
// @Security ApiKeyAuth
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string "Não autorizado"
// @Failure 500 {object} map[string]string
// @Router /agendamentos/{id} [delete]
func (h *AgendamentoHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := h.agendamentoService.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// SolicitarReagendamento solicita reagendamento
// @Summary Solicitar reagendamento
// @Description Solicita reagendamento de uma sessão
// @Tags agendamentos
// @Produce json
// @Param id path string true "ID do agendamento"
// @Security ApiKeyAuth
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string "Não autorizado"
// @Failure 403 {object} map[string]string "Proibido"
// @Failure 500 {object} map[string]string
// @Router /agendamentos/{id}/solicitar-reagendamento [post]
func (h *AgendamentoHandler) SolicitarReagendamento(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	
	// Obter ID do cliente a partir do token JWT
	usuarioID, exists := c.Get(middleware.UsuarioIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}
	
	usuarioAtualID, ok := usuarioID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar ID do usuário"})
		return
	}
	
	// Verificar permissões do usuário (admin pode reagendar qualquer agendamento)
	// Obter detalhes do agendamento para verificar se o cliente é o dono
	agendamento, err := h.agendamentoService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Agendamento não encontrado"})
		return
	}
	
	// Extrair perfil do usuário atual do token
	perfilRaw, exists := c.Get(middleware.PerfilKey)
	perfil := ""
	if exists {
		perfil, _ = perfilRaw.(string)
	}
	
	// Permitir a operação se for admin, massoterapeuta ou o próprio cliente dono do agendamento
	podeReagendar := perfil == "admin" || perfil == "massoterapeuta" || agendamento.ClienteID == usuarioAtualID
	
	if !podeReagendar {
		c.JSON(http.StatusForbidden, gin.H{"error": "Você não tem permissão para reagendar este agendamento"})
		return
	}
	
	// Solicitar reagendamento com o ID do cliente do agendamento (não necessariamente o usuário atual)
	err = h.agendamentoService.SolicitarReagendamento(id, agendamento.ClienteID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Sua solicitação de reagendamento foi enviada. O massoterapeuta entrará em contato para confirmar a nova data.",
	})
}
