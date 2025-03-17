package models

type TipoItem string

const (
	TipoBebida     TipoItem = "BEBIDA"
	TipoIngrediente TipoItem = "INGREDIENTE"
)

type Item struct {
	ID        uint    `gorm:"primaryKey;autoIncrement:false" json:"id"`
	Tipo      TipoItem `gorm:"not null" json:"tipo"`
	Descricao string  `gorm:"not null" json:"descricao"`
	Preco     float64 `gorm:"not null" json:"preco"`
	Extra     bool    `json:"extra"` // true para "Com açúcar" em bebidas ou "Adicional" em ingredientes
}

func (Item) TableName() string {
	return "items"
}

type ItemResponse struct {
	ID        uint    `json:"id"`
	Tipo      string  `json:"tipo"`
	Descricao string  `json:"descricao"`
	Preco     float64 `json:"preco"`
	Extra     bool    `json:"extra"`
}

type ItemRequest struct {
	ID        uint    `json:"id" binding:"required"`
	Tipo      string  `json:"tipo" binding:"required"`
	Descricao string  `json:"descricao" binding:"required"`
	Preco     float64 `json:"preco" binding:"required"`
	Extra     bool    `json:"extra" binding:"required"`
}

type ItemUpdateRequest struct {
	Descricao string  `json:"descricao" binding:"required"`
	Preco     float64 `json:"preco" binding:"required"`
	Extra     *bool   `json:"extra" binding:"required"`
}