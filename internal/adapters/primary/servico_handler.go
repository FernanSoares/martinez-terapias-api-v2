package primary

import (
	"net/http"

	"martinezterapias/api/internal/core/domain"
	"martinezterapias/api/internal/core/ports"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ServicoHandler implementa os handlers HTTP para serviços
type ServicoHandler struct {
	servicoService ports.ServicoService
}

// NewServicoHandler cria uma nova instância de ServicoHandler
func NewServicoHandler(servicoService ports.ServicoService) *ServicoHandler {
	return &ServicoHandler{
		servicoService: servicoService,
	}
}

// GetAll retorna todos os serviços
// @Summary Listar serviços
// @Description Retorna a lista de todos os serviços oferecidos pela clínica
// @Tags servicos
// @Produce json
// @Success 200 {array} domain.Servico
// @Failure 500 {object} map[string]string
// @Router /servicos [get]
func (h *ServicoHandler) GetAll(c *gin.Context) {
	servicos, err := h.servicoService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, servicos)
}

// GetByID retorna um serviço pelo ID
// @Summary Obter serviço por ID
// @Description Retorna os detalhes de um serviço específico
// @Tags servicos
// @Produce json
// @Param id path string true "ID do serviço"
// @Success 200 {object} domain.Servico
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /servicos/{id} [get]
func (h *ServicoHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	servico, err := h.servicoService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Serviço não encontrado"})
		return
	}

	c.JSON(http.StatusOK, servico)
}

// Create cria um novo serviço
// @Summary Cadastrar serviço
// @Description Cria um novo registro de serviço oferecido pela clínica
// @Tags servicos
// @Accept json
// @Produce json
// @Param servico body domain.Servico true "Dados do serviço"
// @Success 201 {object} domain.Servico
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /servicos [post]
func (h *ServicoHandler) Create(c *gin.Context) {
	var servico domain.Servico
	if err := c.ShouldBindJSON(&servico); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validações básicas
	if servico.NomeServico == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nome do serviço é obrigatório"})
		return
	}

	if servico.DuracaoMinutos <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Duração deve ser maior que zero"})
		return
	}

	if servico.Valor <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Valor deve ser maior que zero"})
		return
	}

	if err := h.servicoService.Create(&servico); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, servico)
}

// Update atualiza um serviço existente
// @Summary Atualizar serviço
// @Description Atualiza os dados de um serviço existente
// @Tags servicos
// @Accept json
// @Produce json
// @Param id path string true "ID do serviço"
// @Param servico body domain.Servico true "Dados do serviço"
// @Success 200 {object} domain.Servico
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /servicos/{id} [put]
func (h *ServicoHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Verificar se o serviço existe
	_, err = h.servicoService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Serviço não encontrado"})
		return
	}

	var servico domain.Servico
	if err := c.ShouldBindJSON(&servico); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Garantir que o ID na URL corresponde ao ID no corpo da requisição
	servico.ID = id

	if err := h.servicoService.Update(&servico); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, servico)
}

// Delete remove um serviço
// @Summary Excluir serviço
// @Description Remove um serviço do catálogo da clínica
// @Tags servicos
// @Produce json
// @Param id path string true "ID do serviço"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /servicos/{id} [delete]
func (h *ServicoHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := h.servicoService.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
