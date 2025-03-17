package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"lanchonete/database"
	"lanchonete/models"
)

// @Summary Lista todos os pedidos
// @Description Retorna uma lista de todos os pedidos cadastrados, com opção de filtrar por status não finalizado
// @Tags pedidos
// @Accept json
// @Produce json
// @Param status query string false "Filtrar por status não finalizado (true/false)"
// @Success 200 {array} models.PedidoResponse
// @Router /pedidos [get]
func GetAllPedidos(c *gin.Context) {
	var pedidos []models.Pedido
	database.DB.Preload("PedidoHamburgueres.Hamburguer").Preload("PedidoBebidas.Bebida").Find(&pedidos)

	c.JSON(http.StatusOK, pedidos)
}

// @Summary Busca um pedido por ID
// @Description Retorna um pedido específico baseado no ID
// @Tags pedidos
// @Accept json
// @Produce json
// @Param id path string true "ID do Pedido"
// @Success 200 {object} models.PedidoResponse
// @Failure 404 {object} string "Pedido não encontrado"
// @Router /pedidos/{id} [get]
func GetPedidoByID(c *gin.Context) {
	id := c.Param("id")
	var pedido models.Pedido

	if err := database.DB.Preload("PedidoHamburgueres.Hamburguer").Preload("PedidoBebidas.Bebida").First(&pedido, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pedido não encontrado"})
		return
	}

	c.JSON(http.StatusOK, pedido)
}

// @Summary Cria um novo pedido
// @Description Cria um novo pedido com os dados fornecidos
// @Tags pedidos
// @Accept json
// @Produce json
// @Param pedido body models.PedidoRequest true "Dados do Pedido"
// @Success 201 {object} models.PedidoResponse
// @Failure 400 {object} string "Erro na validação dos dados"
// @Router /pedidos [post]
func CreatePedido(c *gin.Context) {
	var request models.PedidoRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Criar o pedido base
	pedido := models.Pedido{
		Descricao:   request.Descricao,
		Nome:        request.Nome,
		Endereco:    request.Endereco,
		Telefone:    request.Telefone,
		Observacoes: request.Observacoes,
		Status:      models.StatusStarted,
	}

	// Iniciar uma transação
	tx := database.DB.Begin()

	// Criar o pedido
	if err := tx.Create(&pedido).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar pedido"})
		return
	}

	valorTotal := 0.0

	// Adicionar hambúrgueres
	for _, hamburguerReq := range request.Hamburgueres {
		var hamburguer models.Hamburguer
		if err := tx.First(&hamburguer, hamburguerReq.ID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Hambúrguer não encontrado"})
			return
		}

		pedidoHamburguer := models.PedidoHamburguer{
			PedidoID:     pedido.ID,
			HamburguerID: hamburguer.ID,
			Quantidade:   hamburguerReq.Quantidade,
		}

		if err := tx.Create(&pedidoHamburguer).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao adicionar hambúrguer ao pedido"})
			return
		}

		valorTotal += hamburguer.Preco * float64(hamburguerReq.Quantidade)
	}

	// Adicionar bebidas
	for _, bebidaReq := range request.Bebidas {
		var bebida models.Item
		if err := tx.First(&bebida, bebidaReq.ID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bebida não encontrada"})
			return
		}

		pedidoBebida := models.PedidoBebida{
			PedidoID:   pedido.ID,
			ItemID:     bebida.ID,
			Quantidade: bebidaReq.Quantidade,
		}

		if err := tx.Create(&pedidoBebida).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao adicionar bebida ao pedido"})
			return
		}

		valorTotal += bebida.Preco * float64(bebidaReq.Quantidade)
	}

	// Atualizar o valor total do pedido
	pedido.ValorTotal = valorTotal
	if err := tx.Save(&pedido).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar valor total do pedido"})
		return
	}

	// Commit da transação
	tx.Commit()

	// Carregar os relacionamentos para retornar
	database.DB.Preload("PedidoHamburgueres.Hamburguer").Preload("PedidoBebidas.Bebida").First(&pedido, pedido.ID)

	c.JSON(http.StatusCreated, pedido)
}

// @Summary Atualiza um pedido existente
// @Description Atualiza um pedido existente com os dados fornecidos
// @Tags pedidos
// @Accept json
// @Produce json
// @Param id path string true "ID do Pedido"
// @Param pedido body models.PedidoUpdateRequest true "Dados do Pedido"
// @Success 200 {object} models.PedidoResponse
// @Failure 400 {object} string "Erro na validação dos dados"
// @Failure 404 {object} string "Pedido não encontrado"
// @Router /pedidos/{id} [put]
func UpdatePedido(c *gin.Context) {
	id := c.Param("id")
	var pedido models.Pedido
	var request models.PedidoUpdateRequest

	if err := database.DB.First(&pedido, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pedido não encontrado"})
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := database.DB.Begin()

	// Atualizar campos básicos se fornecidos
	if request.Descricao != "" {
		pedido.Descricao = request.Descricao
	}
	if request.Status != "" {
		pedido.Status = request.Status
	}
	if request.Nome != "" {
		pedido.Nome = request.Nome
	}
	if request.Endereco != "" {
		pedido.Endereco = request.Endereco
	}
	if request.Telefone != "" {
		pedido.Telefone = request.Telefone
	}
	if request.Observacoes != "" {
		pedido.Observacoes = request.Observacoes
	}

	valorTotal := 0.0

	// Atualizar hambúrgueres se fornecidos
	if len(request.Hamburgueres) > 0 {
		// Remover relacionamentos existentes
		if err := tx.Where("pedido_id = ?", pedido.ID).Delete(&models.PedidoHamburguer{}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar hambúrgueres"})
			return
		}

		// Adicionar novos relacionamentos
		for _, hamburguerReq := range request.Hamburgueres {
			var hamburguer models.Hamburguer
			if err := tx.First(&hamburguer, hamburguerReq.ID).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{"error": "Hambúrguer não encontrado"})
				return
			}

			pedidoHamburguer := models.PedidoHamburguer{
				PedidoID:     pedido.ID,
				HamburguerID: hamburguer.ID,
				Quantidade:   hamburguerReq.Quantidade,
			}

			if err := tx.Create(&pedidoHamburguer).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao adicionar hambúrguer ao pedido"})
				return
			}

			valorTotal += hamburguer.Preco * float64(hamburguerReq.Quantidade)
		}
	} else {
		// Se não foram fornecidos novos hambúrgueres, calcular o valor total com os existentes
		var pedidoHamburgueres []models.PedidoHamburguer
		tx.Preload("Hamburguer").Where("pedido_id = ?", pedido.ID).Find(&pedidoHamburgueres)
		for _, ph := range pedidoHamburgueres {
			valorTotal += ph.Hamburguer.Preco * float64(ph.Quantidade)
		}
	}

	// Atualizar bebidas se fornecidas
	if len(request.Bebidas) > 0 {
		// Remover relacionamentos existentes
		if err := tx.Where("pedido_id = ?", pedido.ID).Delete(&models.PedidoBebida{}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar bebidas"})
			return
		}

		// Adicionar novos relacionamentos
		for _, bebidaReq := range request.Bebidas {
			var bebida models.Item
			if err := tx.First(&bebida, bebidaReq.ID).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{"error": "Bebida não encontrada"})
				return
			}

			pedidoBebida := models.PedidoBebida{
				PedidoID:   pedido.ID,
				ItemID:     bebida.ID,
				Quantidade: bebidaReq.Quantidade,
			}

			if err := tx.Create(&pedidoBebida).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao adicionar bebida ao pedido"})
				return
			}

			valorTotal += bebida.Preco * float64(bebidaReq.Quantidade)
		}
	} else {
		// Se não foram fornecidas novas bebidas, calcular o valor total com as existentes
		var pedidoBebidas []models.PedidoBebida
		tx.Preload("Bebida").Where("pedido_id = ?", pedido.ID).Find(&pedidoBebidas)
		for _, pb := range pedidoBebidas {
			valorTotal += pb.Bebida.Preco * float64(pb.Quantidade)
		}
	}

	// Atualizar o valor total
	pedido.ValorTotal = valorTotal

	if err := tx.Save(&pedido).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar pedido"})
		return
	}

	tx.Commit()

	// Carregar os relacionamentos atualizados
	database.DB.Preload("PedidoHamburgueres.Hamburguer").Preload("PedidoBebidas.Bebida").First(&pedido, pedido.ID)

	c.JSON(http.StatusOK, pedido)
}

// @Summary Deleta um pedido
// @Description Deleta um pedido existente
// @Tags pedidos
// @Accept json
// @Produce json
// @Param id path string true "ID do Pedido"
// @Success 204 "No Content"
// @Failure 400 {object} string "Erro ao deletar pedido"
// @Failure 404 {object} string "Pedido não encontrado"
// @Router /pedidos/{id} [delete]
func DeletePedido(c *gin.Context) {
	id := c.Param("id")
	var pedido models.Pedido

	if err := database.DB.First(&pedido, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pedido não encontrado"})
		return
	}

	tx := database.DB.Begin()

	// Remover relacionamentos com hambúrgueres
	if err := tx.Where("pedido_id = ?", pedido.ID).Delete(&models.PedidoHamburguer{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar relacionamentos com hambúrgueres"})
		return
	}

	// Remover relacionamentos com bebidas
	if err := tx.Where("pedido_id = ?", pedido.ID).Delete(&models.PedidoBebida{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar relacionamentos com bebidas"})
		return
	}

	// Deletar o pedido
	if err := tx.Delete(&pedido).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar pedido"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"message": "Pedido deletado com sucesso"})
}

