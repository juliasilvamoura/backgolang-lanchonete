package database

import (
	"lanchonete/models"
	"log"
)

func CleanDB() {
	// Limpa todas as tabelas
	DB.Exec("TRUNCATE TABLE pedido_hamburgueres CASCADE")
	DB.Exec("TRUNCATE TABLE pedido_bebidas CASCADE")
	DB.Exec("TRUNCATE TABLE hamburguer_ingredientes CASCADE")
	DB.Exec("TRUNCATE TABLE pedidos CASCADE")
	DB.Exec("TRUNCATE TABLE hamburguers CASCADE")
	DB.Exec("TRUNCATE TABLE items CASCADE")
}

func SeedDB() {
	log.Println("Iniciando seed do banco de dados...")
	
	// Criando itens (bebidas e ingredientes)
	itens := []models.Item{
		// Bebidas
		{ID: 1, Tipo: models.TipoBebida, Descricao: "Coca-Cola 350ml", Preco: 5.00, Extra: true},
		{ID: 2, Tipo: models.TipoBebida, Descricao: "Coca-Cola Zero 350ml", Preco: 5.00, Extra: false},
		{ID: 3, Tipo: models.TipoBebida, Descricao: "Guaraná Antarctica 350ml", Preco: 4.50, Extra: true},
		{ID: 4, Tipo: models.TipoBebida, Descricao: "Água Mineral 500ml", Preco: 3.00, Extra: false},
		
		// Ingredientes
		{ID: 5, Tipo: models.TipoIngrediente, Descricao: "Pão Brioche", Preco: 2.00, Extra: false},
		{ID: 6, Tipo: models.TipoIngrediente, Descricao: "Hambúrguer 180g", Preco: 8.00, Extra: false},
		{ID: 7, Tipo: models.TipoIngrediente, Descricao: "Queijo Cheddar", Preco: 3.00, Extra: false},
		{ID: 8, Tipo: models.TipoIngrediente, Descricao: "Bacon", Preco: 4.00, Extra: true},
		{ID: 9, Tipo: models.TipoIngrediente, Descricao: "Alface", Preco: 1.00, Extra: false},
		{ID: 10, Tipo: models.TipoIngrediente, Descricao: "Tomate", Preco: 1.00, Extra: false},
		{ID: 11, Tipo: models.TipoIngrediente, Descricao: "Cebola Caramelizada", Preco: 2.00, Extra: true},
		{ID: 12, Tipo: models.TipoIngrediente, Descricao: "Ovo", Preco: 2.50, Extra: true},
		{ID: 13, Tipo: models.TipoIngrediente, Descricao: "Molho Especial", Preco: 1.50, Extra: false},
	}

	for _, item := range itens {
		if err := DB.Create(&item).Error; err != nil {
			log.Printf("Erro ao criar item %s: %v\n", item.Descricao, err)
		}
	}

	// Criando hambúrgueres
	hamburgueres := []models.Hamburguer{
		{
			ID: 1,
			Descricao: "Classic Burger",
			Preco: 25.90,
		},
		{
			ID: 2,
			Descricao: "Bacon Burger",
			Preco: 29.90,
		},
		{
			ID: 3,
			Descricao: "Mega Burger",
			Preco: 34.90,
		},
		{
			ID: 4,
			Descricao: "Duplo Burger Duplo Queijo",
			Preco: 39.90,
		},
	}

	// Mapeamento dos ingredientes para cada hambúrguer
	hamburguerIngredientes := map[uint][]models.HamburguerIngrediente{
		1: { // Classic Burger
			{ItemID: 5, Quantidade: 1}, // Pão Brioche
			{ItemID: 6, Quantidade: 1}, // Hambúrguer 180g
			{ItemID: 7, Quantidade: 1}, // Queijo Cheddar
			{ItemID: 9, Quantidade: 1}, // Alface
			{ItemID: 10, Quantidade: 1}, // Tomate
			{ItemID: 13, Quantidade: 1}, // Molho Especial
		},
		2: { // Bacon Burger
			{ItemID: 5, Quantidade: 1}, // Pão Brioche
			{ItemID: 6, Quantidade: 1}, // Hambúrguer 180g
			{ItemID: 7, Quantidade: 1}, // Queijo Cheddar
			{ItemID: 8, Quantidade: 2}, // Bacon (2 fatias)
			{ItemID: 11, Quantidade: 1}, // Cebola Caramelizada
			{ItemID: 13, Quantidade: 1}, // Molho Especial
		},
		3: { // Mega Burger
			{ItemID: 5, Quantidade: 1}, // Pão Brioche
			{ItemID: 6, Quantidade: 1}, // Hambúrguer 180g
			{ItemID: 7, Quantidade: 1}, // Queijo Cheddar
			{ItemID: 8, Quantidade: 2}, // Bacon (2 fatias)
			{ItemID: 12, Quantidade: 1}, // Ovo
			{ItemID: 11, Quantidade: 1}, // Cebola Caramelizada
			{ItemID: 13, Quantidade: 1}, // Molho Especial
		},
		4: { // Duplo Burger Duplo Queijo
			{ItemID: 5, Quantidade: 1}, // Pão Brioche
			{ItemID: 6, Quantidade: 2}, // Hambúrguer 180g (2 unidades)
			{ItemID: 7, Quantidade: 2}, // Queijo Cheddar (2 fatias)
			{ItemID: 8, Quantidade: 2}, // Bacon (2 fatias)
			{ItemID: 13, Quantidade: 1}, // Molho Especial
		},
	}

	for _, hamburguer := range hamburgueres {
		if err := DB.Create(&hamburguer).Error; err != nil {
			log.Printf("Erro ao criar hambúrguer %s: %v\n", hamburguer.Descricao, err)
			continue
		}

		// Adiciona os ingredientes com suas quantidades
		for _, ingrediente := range hamburguerIngredientes[hamburguer.ID] {
			ingrediente.HamburguerID = hamburguer.ID
			if err := DB.Create(&ingrediente).Error; err != nil {
				log.Printf("Erro ao adicionar ingrediente ao hambúrguer %s: %v\n", hamburguer.Descricao, err)
			}
		}
	}

	// Filtrando bebidas para criar pedidos
	var bebidasDisponiveis []models.Item
	DB.Where("tipo = ?", models.TipoBebida).Find(&bebidasDisponiveis)

	// Criando pedidos
	pedidos := []models.Pedido{
		{
			Descricao: "Pedido para João",
			Status: models.StatusStarted,
			Nome: "João Silva",
			Endereco: "Rua das Flores, 123",
			Telefone: "11999999999",
			Observacoes: "Sem cebola, por favor",
			ValorTotal: hamburgueres[0].Preco + bebidasDisponiveis[0].Preco,
		},
		{
			Descricao: "Pedido para Maria",
			Status: models.StatusDelivery,
			Nome: "Maria Santos",
			Endereco: "Av. Principal, 456",
			Telefone: "11988888888",
			Observacoes: "Bacon bem passado",
			ValorTotal: hamburgueres[1].Preco + hamburgueres[2].Preco + bebidasDisponiveis[1].Preco + bebidasDisponiveis[3].Preco,
		},
		{
			Descricao: "Pedido para Pedro",
			Status: models.StatusStarted,
			Nome: "Pedro Oliveira",
			Endereco: "Rua dos Pinheiros, 789",
			Telefone: "11966666666",
			Observacoes: "Todos os hambúrgueres sem tomate",
			ValorTotal: (hamburgueres[0].Preco * 2) + (bebidasDisponiveis[0].Preco * 2),
		},
	}

	for _, pedido := range pedidos {
		if err := DB.Create(&pedido).Error; err != nil {
			log.Printf("Erro ao criar pedido %s: %v\n", pedido.Descricao, err)
			continue
		}

		// Criar relacionamentos com hambúrgueres
		switch pedido.Descricao {
		case "Pedido para João":
			pedidoHamburguer := models.PedidoHamburguer{
				PedidoID: pedido.ID,
				HamburguerID: hamburgueres[0].ID,
				Quantidade: 1,
			}
			if err := DB.Create(&pedidoHamburguer).Error; err != nil {
				log.Printf("Erro ao criar relação pedido-hamburguer: %v\n", err)
			}

			pedidoBebida := models.PedidoBebida{
				PedidoID: pedido.ID,
				ItemID: bebidasDisponiveis[0].ID,
				Quantidade: 1,
			}
			if err := DB.Create(&pedidoBebida).Error; err != nil {
				log.Printf("Erro ao criar relação pedido-bebida: %v\n", err)
			}

		case "Pedido para Maria":
			// Primeiro hambúrguer
			pedidoHamburguer1 := models.PedidoHamburguer{
				PedidoID: pedido.ID,
				HamburguerID: hamburgueres[1].ID,
				Quantidade: 1,
			}
			if err := DB.Create(&pedidoHamburguer1).Error; err != nil {
				log.Printf("Erro ao criar relação pedido-hamburguer: %v\n", err)
			}

			// Segundo hambúrguer
			pedidoHamburguer2 := models.PedidoHamburguer{
				PedidoID: pedido.ID,
				HamburguerID: hamburgueres[2].ID,
				Quantidade: 1,
			}
			if err := DB.Create(&pedidoHamburguer2).Error; err != nil {
				log.Printf("Erro ao criar relação pedido-hamburguer: %v\n", err)
			}

			// Primeira bebida
			pedidoBebida1 := models.PedidoBebida{
				PedidoID: pedido.ID,
				ItemID: bebidasDisponiveis[1].ID,
				Quantidade: 1,
			}
			if err := DB.Create(&pedidoBebida1).Error; err != nil {
				log.Printf("Erro ao criar relação pedido-bebida: %v\n", err)
			}

			// Segunda bebida
			pedidoBebida2 := models.PedidoBebida{
				PedidoID: pedido.ID,
				ItemID: bebidasDisponiveis[3].ID,
				Quantidade: 1,
			}
			if err := DB.Create(&pedidoBebida2).Error; err != nil {
				log.Printf("Erro ao criar relação pedido-bebida: %v\n", err)
			}

		case "Pedido para Pedro":
			// Dois hambúrgueres iguais
			pedidoHamburguer := models.PedidoHamburguer{
				PedidoID: pedido.ID,
				HamburguerID: hamburgueres[0].ID,
				Quantidade: 2,
			}
			if err := DB.Create(&pedidoHamburguer).Error; err != nil {
				log.Printf("Erro ao criar relação pedido-hamburguer: %v\n", err)
			}

			// Duas bebidas iguais
			pedidoBebida := models.PedidoBebida{
				PedidoID: pedido.ID,
				ItemID: bebidasDisponiveis[0].ID,
				Quantidade: 2,
			}
			if err := DB.Create(&pedidoBebida).Error; err != nil {
				log.Printf("Erro ao criar relação pedido-bebida: %v\n", err)
			}
		}
	}

	log.Println("Seed do banco de dados concluído!")
} 