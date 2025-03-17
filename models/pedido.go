package models

import (
	"time"
	"github.com/google/uuid"
)

type StatusPedido string

const (
	StatusStarted   StatusPedido = "STARTED"
	StatusDelivery  StatusPedido = "DELIVERY"
	StatusFinalized StatusPedido = "FINALIZED"
)

type PedidoHamburguer struct {
	PedidoID     uuid.UUID  `gorm:"type:uuid;primaryKey"`
	HamburguerID uint       `gorm:"primaryKey"`
	Quantidade   int        `gorm:"not null;default:1"`
	Hamburguer   Hamburguer `gorm:"foreignKey:HamburguerID"`
}

func (PedidoHamburguer) TableName() string {
	return "pedido_hamburgueres"
}

type PedidoBebida struct {
	PedidoID   uuid.UUID `gorm:"type:uuid;primaryKey"`
	ItemID     uint      `gorm:"primaryKey"`
	Quantidade int       `gorm:"not null;default:1"`
	Bebida     Item      `gorm:"foreignKey:ItemID"`
}

type Pedido struct {
	ID           uuid.UUID     `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Data         time.Time     `gorm:"not null;default:CURRENT_TIMESTAMP" json:"data"`
	Descricao    string        `gorm:"not null" json:"descricao" binding:"required"`
	Status       StatusPedido  `gorm:"not null;default:'STARTED'" json:"status"`
	Nome         string        `gorm:"not null" json:"nome" binding:"required"`
	Endereco     string        `gorm:"not null" json:"endereco" binding:"required"`
	Telefone     string        `gorm:"not null" json:"telefone" binding:"required,len=11"`
	Hamburgueres []Hamburguer  `gorm:"many2many:pedido_hamburgueres;foreignKey:ID;joinForeignKey:pedido_id;References:ID;joinReferences:hamburguer_id" json:"-"`
	Bebidas      []Item        `gorm:"many2many:pedido_bebidas;foreignKey:ID;joinForeignKey:pedido_id;References:ID;joinReferences:item_id" json:"-"`
	PedidoHamburgueres []PedidoHamburguer `gorm:"foreignKey:PedidoID" json:"hamburgueres"`
	PedidoBebidas      []PedidoBebida     `gorm:"foreignKey:PedidoID" json:"bebidas"`
	Observacoes  string        `json:"observacoes"`
	ValorTotal   float64       `gorm:"not null" json:"valor_total"`
}

type PedidoResponse struct {
	ID           uuid.UUID          `json:"id"`
	Data         time.Time          `json:"data"`
	Descricao    string            `json:"descricao"`
	Status       StatusPedido       `json:"status"`
	Nome         string            `json:"nome"`
	Endereco     string            `json:"endereco"`
	Telefone     string            `json:"telefone"`
	Hamburgueres []PedidoHamburguer `json:"hamburgueres"`
	Bebidas      []PedidoBebida     `json:"bebidas"`
	Observacoes  string            `json:"observacoes"`
	ValorTotal   float64           `json:"valor_total"`
}

type PedidoRequest struct {
	Descricao      string         `json:"descricao" binding:"required"`
	Nome           string         `json:"nome" binding:"required"`
	Endereco       string         `json:"endereco" binding:"required"`
	Telefone       string         `json:"telefone" binding:"required,len=11"`
	Hamburgueres   []PedidoItemRequest `json:"hamburgueres" binding:"required,min=1"`
	Bebidas        []PedidoItemRequest `json:"bebidas"`
	Observacoes    string         `json:"observacoes"`
}

type PedidoItemRequest struct {
	ID         uint `json:"id" binding:"required"`
	Quantidade int  `json:"quantidade" binding:"required,min=1"`
}

type PedidoUpdateRequest struct {
	Descricao      string             `json:"descricao"`
	Status         StatusPedido       `json:"status"`
	Nome           string             `json:"nome"`
	Endereco       string             `json:"endereco"`
	Telefone       string             `json:"telefone"`
	Hamburgueres   []PedidoItemRequest `json:"hamburgueres"`
	Bebidas        []PedidoItemRequest `json:"bebidas"`
	Observacoes    string             `json:"observacoes"`
}
