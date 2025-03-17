package routes

import (
	"lanchonete/controller"
	_ "lanchonete/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title API Lanchonete
// @version 1.0
// @description API para gerenciamento de pedidos de uma lanchonete
// @host localhost:8080
// @BasePath /
func HandleRequests(r *gin.Engine) {
	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello, World!"})
	})
	
	// Rotas de itens
	r.GET("/itens/todos", controller.GetAllItens)      // Lista todos os itens
	r.GET("/itens/:codigo", controller.GetItem)        // Busca item por código
	r.POST("/itens", controller.CreateItem)            // Cria novo item
	r.PUT("/itens/:codigo", controller.UpdateItem)     // Atualiza item existente
	r.DELETE("/itens/:codigo", controller.DeleteItem)  // Remove item existente
	r.GET("/itens/bebidas", controller.GetBebidas)     // Lista todas as bebidas
	r.GET("/itens/ingredientes", controller.GetIngredientes) // Lista todos os ingredientes

	// Rotas de hamburguers
	r.GET("/hamburguers", controller.GetAllHamburguers)
	r.GET("/hamburguers/:id", controller.GetHamburguerByID)
	r.GET("/hamburguers/nome/:nome", controller.GetHamburguerByName)
	r.POST("/hamburguers", controller.CreateHamburguer)
	r.PUT("/hamburguers/:id", controller.UpdateHamburguer)
	r.DELETE("/hamburguers/:id", controller.DeleteHamburguer)

	// Rotas de pedidos
	r.GET("/pedidos", controller.GetAllPedidos)
	r.GET("/pedidos/:id", controller.GetPedidoByID)
	r.POST("/pedidos", controller.CreatePedido)
	r.PUT("/pedidos/:id", controller.UpdatePedido)
	r.DELETE("/pedidos/:id", controller.DeletePedido)

	r.Run(":8080")
}
