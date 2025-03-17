package main

import (
	"fmt"
	"lanchonete/database"
)

func main() {
	fmt.Println("Conectando ao banco de dados...")
	database.ConnectDB()

	fmt.Println("Limpando dados existentes...")
	database.CleanDB()

	fmt.Println("Populando banco de dados com dados de teste...")
	database.SeedDB()

	fmt.Println("Processo concluído!")
} 