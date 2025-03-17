package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"lanchonete/database"
	"lanchonete/models"
)

// @Summary Lista todos os itens
// @Description Retorna uma lista de todos os itens (bebidas e ingredientes)
// @Tags itens
// @Accept json
// @Produce json
// @Success 200 {array} models.ItemResponse
// @Failure 404 {object} string "Nenhum item encontrado"
// @Router /itens/todos [get]
func GetAllItens(c *gin.Context) {
	var itens []models.Item
	if err := database.DB.Find(&itens).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar itens"})
		return
	}

	if len(itens) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Nenhum item encontrado"})
		return
	}

	var response []models.ItemResponse
	for _, item := range itens {
		response = append(response, models.ItemResponse{
			ID:        item.ID,
			Tipo:      string(item.Tipo),
			Descricao: item.Descricao,
			Preco:     item.Preco,
			Extra:     item.Extra,
		})
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Busca um item por código
// @Description Retorna um item específico (bebida ou ingrediente) buscando por código
// @Tags itens
// @Accept json
// @Produce json
// @Param codigo path string true "Código do item"
// @Success 200 {object} models.ItemResponse
// @Failure 400 {object} string "Código inválido"
// @Failure 404 {object} string "Item não encontrado"
// @Router /itens/{codigo} [get]
func GetItem(c *gin.Context) {
	codigo := c.Param("codigo")
	if codigo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Código do item é obrigatório"})
		return
	}

	id, err := strconv.Atoi(codigo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Código inválido"})
		return
	}

	var item models.Item
	if err := database.DB.First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item não encontrado"})
		return
	}

	c.JSON(http.StatusOK, models.ItemResponse{
		ID:        item.ID,
		Tipo:      string(item.Tipo),
		Descricao: item.Descricao,
		Preco:     item.Preco,
		Extra:     item.Extra,
	})
}

// @Summary Lista todas as bebidas
// @Description Retorna uma lista de todas as bebidas disponíveis
// @Tags itens
// @Accept json
// @Produce json
// @Success 200 {array} models.ItemResponse
// @Failure 404 {object} string "Nenhuma bebida encontrada"
// @Router /itens/bebidas [get]
func GetBebidas(c *gin.Context) {
	var bebidas []models.Item
	if err := database.DB.Where("tipo = ?", models.TipoBebida).Find(&bebidas).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar bebidas"})
		return
	}

	if len(bebidas) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Nenhuma bebida encontrada"})
		return
	}

	var response []models.ItemResponse
	for _, bebida := range bebidas {
		response = append(response, models.ItemResponse{
			ID:        bebida.ID,
			Tipo:      string(bebida.Tipo),
			Descricao: bebida.Descricao,
			Preco:     bebida.Preco,
			Extra:     bebida.Extra,
		})
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Lista todos os ingredientes
// @Description Retorna uma lista de todos os ingredientes disponíveis
// @Tags itens
// @Accept json
// @Produce json
// @Success 200 {array} models.ItemResponse
// @Failure 404 {object} string "Nenhum ingrediente encontrado"
// @Router /itens/ingredientes [get]
func GetIngredientes(c *gin.Context) {
	var ingredientes []models.Item
	if err := database.DB.Where("tipo = ?", models.TipoIngrediente).Find(&ingredientes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar ingredientes"})
		return
	}

	if len(ingredientes) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Nenhum ingrediente encontrado"})
		return
	}

	var response []models.ItemResponse
	for _, ingrediente := range ingredientes {
		response = append(response, models.ItemResponse{
			ID:        ingrediente.ID,
			Tipo:      string(ingrediente.Tipo),
			Descricao: ingrediente.Descricao,
			Preco:     ingrediente.Preco,
			Extra:     ingrediente.Extra,
		})
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Cria um novo item
// @Description Cria um novo item (bebida ou ingrediente) com os dados fornecidos
// @Tags itens
// @Accept json
// @Produce json
// @Param item body models.ItemRequest true "Dados do Item"
// @Success 201 {object} models.ItemResponse
// @Failure 400 {object} string "Erro na validação dos dados"
// @Router /itens [post]
func CreateItem(c *gin.Context) {
	var request models.ItemRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validação do tipo
	if request.Tipo != string(models.TipoBebida) && request.Tipo != string(models.TipoIngrediente) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tipo de item inválido. Use 'BEBIDA' ou 'INGREDIENTE'"})
		return
	}

	// Verifica se o item já existe
	var count int64
	if err := database.DB.Model(&models.Item{}).Where("id = ?", request.ID).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao verificar existência do item"})
		return
	}

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Item já existe"})
		return
	}

	item := models.Item{
		ID:        request.ID,
		Tipo:      models.TipoItem(request.Tipo),
		Descricao: request.Descricao,
		Preco:     request.Preco,
		Extra:     request.Extra,
	}

	if err := database.DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar item"})
		return
	}

	c.JSON(http.StatusCreated, models.ItemResponse{
		ID:        item.ID,
		Tipo:      string(item.Tipo),
		Descricao: item.Descricao,
		Preco:     item.Preco,
		Extra:     item.Extra,
	})
}

// @Summary Atualiza um item existente
// @Description Atualiza um item (bebida ou ingrediente) existente
// @Tags itens
// @Accept json
// @Produce json
// @Param codigo path string true "Código do item"
// @Param item body models.ItemUpdateRequest true "Dados do Item"
// @Success 200 {object} models.ItemResponse
// @Failure 400 {object} string "Código inválido"
// @Failure 404 {object} string "Item não encontrado"
// @Router /itens/{codigo} [put]
func UpdateItem(c *gin.Context) {
	codigo := c.Param("codigo")
	if codigo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Código do item é obrigatório"})
		return
	}

	id, err := strconv.Atoi(codigo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Código inválido"})
		return
	}

	var updateRequest models.ItemUpdateRequest
	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verifica se o item existe
	var item models.Item
	if err := database.DB.First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item não encontrado"})
		return
	}

	// Verifica se o item está em algum pedido não finalizado
	var count int64
	if item.Tipo == models.TipoBebida {
		if err := database.DB.Table("pedido_bebidas").
			Joins("JOIN pedidos ON pedidos.id = pedido_bebidas.pedido_id").
			Where("pedido_bebidas.item_id = ? AND pedidos.status != ?", id, models.StatusFinalized).
			Count(&count).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao verificar pedidos"})
			return
		}
	} else {
		if err := database.DB.Table("hamburguer_ingredientes").
			Joins("JOIN hamburguers ON hamburguers.id = hamburguer_ingredientes.hamburguer_id").
			Joins("JOIN pedido_hamburgueres ON pedido_hamburgueres.hamburguer_id = hamburguers.id").
			Joins("JOIN pedidos ON pedidos.id = pedido_hamburgueres.pedido_id").
			Where("hamburguer_ingredientes.item_id = ? AND pedidos.status != ?", id, models.StatusFinalized).
			Count(&count).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao verificar pedidos"})
			return
		}
	}

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Não é possível atualizar um item que está em pedidos não finalizados"})
		return
	}

	// Atualiza o item
	item.Descricao = updateRequest.Descricao
	item.Preco = updateRequest.Preco
	item.Extra = *updateRequest.Extra

	if err := database.DB.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar item"})
		return
	}

	c.JSON(http.StatusOK, models.ItemResponse{
		ID:        item.ID,
		Tipo:      string(item.Tipo),
		Descricao: item.Descricao,
		Preco:     item.Preco,
		Extra:     item.Extra,
	})
}

// @Summary Deleta um item existente
// @Description Remove um item (bebida ou ingrediente) existente pelo código
// @Tags itens
// @Accept json
// @Produce json
// @Param codigo path string true "Código do item"
// @Success 200 {object} string "Item removido com sucesso"
// @Failure 400 {object} string "Código inválido"
// @Failure 404 {object} string "Item não encontrado"
// @Router /itens/{codigo} [delete]
func DeleteItem(c *gin.Context) {
	codigo := c.Param("codigo")
	if codigo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Código do item é obrigatório"})
		return
	}

	id, err := strconv.Atoi(codigo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Código inválido"})
		return
	}

	// Verifica se o item existe
	var item models.Item
	if err := database.DB.First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item não encontrado"})
		return
	}

	// Verifica se o item está em algum pedido não finalizado
	var count int64
	if item.Tipo == models.TipoBebida {
		if err := database.DB.Table("pedido_bebidas").
			Joins("JOIN pedidos ON pedidos.id = pedido_bebidas.pedido_id").
			Where("pedido_bebidas.item_id = ? AND pedidos.status != ?", id, models.StatusFinalized).
			Count(&count).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao verificar pedidos"})
			return
		}
	} else {
		if err := database.DB.Table("hamburguer_ingredientes").
			Joins("JOIN hamburguers ON hamburguers.id = hamburguer_ingredientes.hamburguer_id").
			Joins("JOIN pedido_hamburgueres ON pedido_hamburgueres.hamburguer_id = hamburguers.id").
			Joins("JOIN pedidos ON pedidos.id = pedido_hamburgueres.pedido_id").
			Where("hamburguer_ingredientes.item_id = ? AND pedidos.status != ?", id, models.StatusFinalized).
			Count(&count).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao verificar pedidos"})
			return
		}
	}

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Não é possível deletar um item que está em pedidos não finalizados"})
		return
	}

	// Deleta o item
	if err := database.DB.Delete(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item removido com sucesso"})
}
