package database

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"lanchonete/models"
)

var (
	DB  *gorm.DB
	err error
)

func ConnectDB() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=root password=root dbname=lanchonete port=5432 sslmode=disable"
	}

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatal("Erro ao conectar ao banco de dados:", err)
	}

	// Desabilita as foreign keys durante a migração
	DB.Exec("ALTER TABLE IF EXISTS pedido_hamburgueres DROP CONSTRAINT IF EXISTS fk_pedido_hamburgueres_pedido")
	DB.Exec("ALTER TABLE IF EXISTS pedido_hamburgueres DROP CONSTRAINT IF EXISTS fk_pedido_hamburgueres_hamburguer")
	DB.Exec("DROP TABLE IF EXISTS pedido_hamburgueres CASCADE")

	// Auto Migrate na ordem correta
	DB.AutoMigrate(
		&models.Item{},
		&models.Hamburguer{},
		&models.HamburguerIngrediente{},
		&models.Pedido{},
		&models.PedidoHamburguer{},
		&models.PedidoBebida{},
	)

	// Habilita as foreign keys após a migração
	DB.Exec("SET CONSTRAINTS ALL IMMEDIATE")
}