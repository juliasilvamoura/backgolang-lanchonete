package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"lanchonete/database"
	"lanchonete/models"
)

// @Summary Lista todos os hamburgueres
// @Description Retorna uma lista de todos os hamburgueres disponíveis
// @Tags hamburgueres
// @Accept json
// @Produce json
// @Success 200 {array} models.Hamburguer
// @Router /hamburguers [get]
func GetAllHamburguers(c *gin.Context) {
	var hamburguers []models.Hamburguer
	database.DB.Preload("HamburguerIngredientes.Item").Find(&hamburguers)
	c.JSON(http.StatusOK, hamburguers)
}

// @Summary Busca um hamburguer por ID
// @Description Retorna um hamburguer específico baseado no ID
// @Tags hamburgueres
// @Accept json
// @Produce json
// @Param id path int true "ID do Hamburguer"
// @Success 200 {object} models.Hamburguer
// @Failure 404 {object} string "Hamburguer não encontrado"
// @Router /hamburguers/{id} [get]
func GetHamburguerByID(c *gin.Context) {
	id := c.Param("id")
	var hamburguer models.Hamburguer
	
	// Verifica se existe um hambúrguer com este ID
	var count int64
	if err := database.DB.Model(&models.Hamburguer{}).Where("id = ?", id).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar hambúrguer"})
		return
	}
	
	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hambúrguer não encontrado"})
		return
	}

	// Busca o hambúrguer com seus ingredientes
	if err := database.DB.Preload("HamburguerIngredientes.Item").First(&hamburguer, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar detalhes do hambúrguer"})
		return
	}

	c.JSON(http.StatusOK, hamburguer)
}

// @Summary Busca um hamburguer por nome
// @Description Retorna um hamburguer específico baseado no nome
// @Tags hamburgueres
// @Accept json
// @Produce json
// @Param nome path string true "Nome do Hamburguer"
// @Success 200 {object} models.Hamburguer
// @Failure 404 {object} string
// @Router /hamburguers/nome/{nome} [get]
func GetHamburguerByName(c *gin.Context) {
	name := c.Param("Descricao")
	var hamburguers []models.Hamburguer
	database.DB.Preload("HamburguerIngredientes.Item").Where("descricao ILIKE ?", "%"+name+"%").Find(&hamburguers)
	c.JSON(http.StatusOK, hamburguers)
}

// @Summary Cria um novo hamburguer
// @Description Cria um novo hamburguer com os dados fornecidos
// @Tags hamburgueres
// @Accept json
// @Produce json
// @Param hamburguer body models.HamburguerRequest true "Dados do Hamburguer"
// @Success 201 {object} models.Hamburguer
// @Failure 400 {object} string "Erro na validação dos dados"
// @Failure 404 {object} string "Ingrediente não encontrado"
// @Router /hamburguers [post]
func CreateHamburguer(c *gin.Context) {
	var request models.HamburguerRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
		return
	}

	// Verifica se o hambúrguer já existe
	var count int64
	if err := database.DB.Model(&models.Hamburguer{}).Where("id = ?", request.ID).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao verificar hambúrguer existente"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Já existe um hambúrguer com este ID"})
		return
	}

	// Inicia uma transação
	tx := database.DB.Begin()

	// Cria o hambúrguer
	hamburguer := models.Hamburguer{
		ID:        request.ID,
		Descricao: request.Descricao,
		Preco:     request.Preco,
	}

	if err := tx.Create(&hamburguer).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar hambúrguer"})
		return
	}

	// Adiciona os ingredientes com suas quantidades
	for _, ingrediente := range request.Ingredientes {
		var item models.Item
		if err := tx.First(&item, ingrediente.ID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"error": "Ingrediente não encontrado: " + string(ingrediente.ID)})
			return
		}

		// Verifica se é um ingrediente
		if item.Tipo != models.TipoIngrediente {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "O item " + string(ingrediente.ID) + " não é um ingrediente"})
			return
		}

		hamburguerIngrediente := models.HamburguerIngrediente{
			HamburguerID: hamburguer.ID,
			ItemID:       ingrediente.ID,
			Quantidade:   ingrediente.Quantidade,
		}

		if err := tx.Create(&hamburguerIngrediente).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao adicionar ingrediente ao hambúrguer"})
			return
		}
	}

	// Commit da transação
	tx.Commit()

	// Carrega os relacionamentos para retornar
	database.DB.Preload("HamburguerIngredientes.Item").First(&hamburguer, hamburguer.ID)
	c.JSON(http.StatusCreated, hamburguer)
}

// @Summary Atualiza um hamburguer existente
// @Description Atualiza um hamburguer existente com os dados fornecidos
// @Tags hamburgueres
// @Accept json
// @Produce json
// @Param id path int true "ID do Hamburguer"
// @Param hamburguer body models.HamburguerUpdateRequest true "Dados do Hamburguer"
// @Success 200 {object} models.Hamburguer
// @Failure 400 {object} string "Erro na validação dos dados"
// @Failure 404 {object} string "Hamburguer não encontrado"
// @Router /hamburguers/{id} [put]
func UpdateHamburguer(c *gin.Context) {
	id := c.Param("id")

	// Verifica se o hambúrguer existe
	var count int64
	if err := database.DB.Model(&models.Hamburguer{}).Where("id = ?", id).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao verificar hambúrguer"})
		return
	}
	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hambúrguer não encontrado"})
		return
	}

	// Verifica se o hambúrguer está em algum pedido não finalizado
	if err := database.DB.Table("pedido_hamburgueres").
		Joins("JOIN pedidos ON pedidos.id = pedido_hamburgueres.pedido_id").
		Where("pedido_hamburgueres.hamburguer_id = ? AND pedidos.status != ?", id, models.StatusFinalized).
		Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao verificar pedidos"})
		return
	}

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Não é possível atualizar um hambúrguer que está em pedidos não finalizados"})
		return
	}

	var request models.HamburguerUpdateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
		return
	}

	// Inicia uma transação
	tx := database.DB.Begin()

	// Atualiza o hambúrguer
	var hamburguer models.Hamburguer
	if err := tx.First(&hamburguer, id).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar hambúrguer"})
		return
	}

	hamburguer.Descricao = request.Descricao
	hamburguer.Preco = request.Preco

	if err := tx.Save(&hamburguer).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao salvar hambúrguer"})
		return
	}

	// Remove os ingredientes antigos
	if err := tx.Where("hamburguer_id = ?", hamburguer.ID).Delete(&models.HamburguerIngrediente{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao remover ingredientes antigos"})
		return
	}

	// Adiciona os novos ingredientes
	for _, ingrediente := range request.Ingredientes {
		var item models.Item
		if err := tx.First(&item, ingrediente.ID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"error": "Ingrediente não encontrado: " + string(ingrediente.ID)})
			return
		}

		// Verifica se é um ingrediente
		if item.Tipo != models.TipoIngrediente {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "O item " + string(ingrediente.ID) + " não é um ingrediente"})
			return
		}

		hamburguerIngrediente := models.HamburguerIngrediente{
			HamburguerID: hamburguer.ID,
			ItemID:       ingrediente.ID,
			Quantidade:   ingrediente.Quantidade,
		}

		if err := tx.Create(&hamburguerIngrediente).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao adicionar ingrediente ao hambúrguer"})
			return
		}
	}

	// Commit da transação
	tx.Commit()

	// Carrega os relacionamentos para retornar
	database.DB.Preload("HamburguerIngredientes.Item").First(&hamburguer, hamburguer.ID)
	c.JSON(http.StatusOK, hamburguer)
}

// @Summary Deleta um hamburguer
// @Description Deleta um hamburguer existente
// @Tags hamburgueres
// @Accept json
// @Produce json
// @Param id path int true "ID do Hamburguer"
// @Success 204 "No Content"
// @Failure 400 {object} string "Erro ao deletar hamburguer"
// @Failure 404 {object} string "Hamburguer não encontrado"
// @Router /hamburguers/{id} [delete]
func DeleteHamburguer(c *gin.Context) {
	id := c.Param("id")

	// Verifica se o hambúrguer existe
	var count int64
	if err := database.DB.Model(&models.Hamburguer{}).Where("id = ?", id).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao verificar hambúrguer"})
		return
	}
	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hambúrguer não encontrado"})
		return
	}

	// Verifica se o hambúrguer está em algum pedido não finalizado
	if err := database.DB.Table("pedido_hamburgueres").
		Joins("JOIN pedidos ON pedidos.id = pedido_hamburgueres.pedido_id").
		Where("pedido_hamburgueres.hamburguer_id = ? AND pedidos.status != ?", id, models.StatusFinalized).
		Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao verificar pedidos"})
		return
	}

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Não é possível deletar um hambúrguer que está em pedidos não finalizados"})
		return
	}

	// Inicia uma transação
	tx := database.DB.Begin()

	// Remove os ingredientes
	if err := tx.Where("hamburguer_id = ?", id).Delete(&models.HamburguerIngrediente{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao remover ingredientes"})
		return
	}

	// Deleta o hambúrguer
	if err := tx.Delete(&models.Hamburguer{}, id).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar hambúrguer"})
		return
	}

	// Commit da transação
	tx.Commit()

	c.Status(http.StatusNoContent)
}

