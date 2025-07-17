package primary

import (
	"net/http"

	"martinezterapias/api/internal/core/domain"
	"martinezterapias/api/internal/core/ports"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AnamneseHandler implementa os handlers HTTP para fichas de anamnese
type AnamneseHandler struct {
	anamneseService ports.AnamneseService
}

// NewAnamneseHandler cria uma nova instância de AnamneseHandler
func NewAnamneseHandler(anamneseService ports.AnamneseService) *AnamneseHandler {
	return &AnamneseHandler{
		anamneseService: anamneseService,
	}
}

// GetAll retorna todas as fichas de anamnese
// @Summary Listar fichas de anamnese
// @Description Retorna a lista de todas as fichas de anamnese
// @Tags anamnese
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} domain.Anamnese
// @Failure 401 {object} map[string]string "Não autorizado"
// @Failure 500 {object} map[string]string
// @Router /anamnese [get]
func (h *AnamneseHandler) GetAll(c *gin.Context) {
	anamneses, err := h.anamneseService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, anamneses)
}

// GetByID retorna uma ficha de anamnese pelo ID
// @Summary Obter ficha de anamnese por ID
// @Description Retorna os detalhes de uma ficha de anamnese específica
// @Tags anamnese
// @Produce json
// @Param id path string true "ID da ficha de anamnese"
// @Security ApiKeyAuth
// @Success 200 {object} domain.Anamnese
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string "Não autorizado"
// @Failure 404 {object} map[string]string
// @Router /anamnese/{id} [get]
func (h *AnamneseHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	anamnese, err := h.anamneseService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ficha de anamnese não encontrada"})
		return
	}

	c.JSON(http.StatusOK, anamnese)
}

// GetByClienteID retorna todas as fichas de anamnese de um cliente
// @Summary Obter fichas de anamnese por cliente
// @Description Retorna todas as fichas de anamnese de um cliente específico
// @Tags anamnese
// @Produce json
// @Param cliente_id path string true "ID do cliente"
// @Security ApiKeyAuth
// @Success 200 {array} domain.Anamnese
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string "Não autorizado"
// @Failure 500 {object} map[string]string
// @Router /clientes/{cliente_id}/anamnese [get]
func (h *AnamneseHandler) GetByClienteID(c *gin.Context) {
	clienteID, err := uuid.Parse(c.Param("cliente_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do cliente inválido"})
		return
	}

	anamneses, err := h.anamneseService.GetByClienteID(clienteID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, anamneses)
}

// Create cria uma nova ficha de anamnese para um cliente
// @Summary Criar ficha de anamnese
// @Description Cria uma nova ficha de anamnese para um cliente
// @Tags anamnese
// @Accept json
// @Produce json
// @Param cliente_id path string true "ID do cliente"
// @Param anamnese body domain.Anamnese true "Dados da ficha de anamnese"
// @Security ApiKeyAuth
// @Success 201 {object} domain.Anamnese
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string "Não autorizado"
// @Failure 500 {object} map[string]string
// @Router /clientes/{cliente_id}/anamnese [post]
func (h *AnamneseHandler) Create(c *gin.Context) {
	clienteID, err := uuid.Parse(c.Param("cliente_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do cliente inválido"})
		return
	}

	var anamnese domain.Anamnese
	if err := c.ShouldBindJSON(&anamnese); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Garantir que o cliente_id da URL seja usado na anamnese
	anamnese.ClienteID = clienteID

	if err := h.anamneseService.Create(&anamnese); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, anamnese)
}

// Update atualiza uma ficha de anamnese existente
// @Summary Atualizar ficha de anamnese
// @Description Atualiza os dados de uma ficha de anamnese existente
// @Tags anamnese
// @Accept json
// @Produce json
// @Param id path string true "ID da ficha de anamnese"
// @Param anamnese body domain.Anamnese true "Dados da ficha de anamnese"
// @Security ApiKeyAuth
// @Success 200 {object} domain.Anamnese
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string "Não autorizado"
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /anamnese/{id} [put]
func (h *AnamneseHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Verificar se a ficha existe
	existente, err := h.anamneseService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ficha de anamnese não encontrada"})
		return
	}

	var anamnese domain.Anamnese
	if err := c.ShouldBindJSON(&anamnese); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Garantir que o ID na URL corresponde ao ID no corpo da requisição
	anamnese.ID = id
	
	// Manter o ID do cliente e data de preenchimento original
	anamnese.ClienteID = existente.ClienteID
	if !existente.DataPreenchimento.IsZero() {
		anamnese.DataPreenchimento = existente.DataPreenchimento
	}

	if err := h.anamneseService.Update(&anamnese); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, anamnese)
}
