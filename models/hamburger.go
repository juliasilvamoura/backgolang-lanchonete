package models

type HamburguerIngrediente struct {
	HamburguerID uint `gorm:"primaryKey"`
	ItemID       uint `gorm:"primaryKey"`
	Quantidade   int  `gorm:"not null;default:1"`
	Item         Item `gorm:"foreignKey:ItemID"`
}

type Hamburguer struct {
	ID          uint    `gorm:"primaryKey; autoIncrement:false" json:"id"`
	Descricao   string  `gorm:"not null" json:"descricao" binding:"required"`
	Preco       float64 `gorm:"not null" json:"preco" binding:"required,gt=0"`
	Ingredientes []Item  `gorm:"many2many:hamburguer_ingredientes;foreignKey:ID;joinForeignKey:hamburguer_id;References:ID;joinReferences:item_id" json:"-"`
	HamburguerIngredientes []HamburguerIngrediente `gorm:"foreignKey:HamburguerID" json:"ingredientes"`
}

// HamburguerRequest é o modelo para criar um novo hambúrguer
type HamburguerRequest struct {
	ID            uint    `json:"id" binding:"required"`
	Descricao     string  `json:"descricao" binding:"required"`
	Preco         float64 `json:"preco" binding:"required,gt=0"`
	Ingredientes  []IngredienteRequest `json:"ingredientes" binding:"required,min=1"`
}

type IngredienteRequest struct {
	ID         uint `json:"id" binding:"required"`
	Quantidade int  `json:"quantidade" binding:"required,min=1"`
}

// HamburguerUpdateRequest é o modelo para atualizar um hambúrguer existente
type HamburguerUpdateRequest struct {
	Descricao    string  `json:"descricao" binding:"required"`
	Preco        float64 `json:"preco" binding:"required,gt=0"`
	Ingredientes []IngredienteRequest `json:"ingredientes" binding:"required,min=1"`
}