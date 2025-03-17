package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"lanchonete/database"
	"lanchonete/routes"
)

func main() {
	fmt.Println("Iniciando o servidor da API")
	
	// Inicializa o router
	r := gin.Default()

	// Configuração do CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// Conectar ao banco
	database.ConnectDB()
	
	// Configurar rotas passando o router
	routes.HandleRequests(r)
}