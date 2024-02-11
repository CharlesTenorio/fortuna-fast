package model

import (
	"encoding/json"
	"time"

	"github.com/katana/fortuna/backend-go/internal/config/logger"
	"github.com/katana/fortuna/backend-go/pkg/service/validation"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Cliente struct {
	ID          primitive.ObjectID `bson:"_id" json:"_id"`
	DataType    string             `bson:"data_type" json:"-"`
	TipoCliente string             `bson:"tipo" json:"tipo_cliente"`
	IDUsuario   primitive.ObjectID `bson:"user_id " json:"id_usr"`
	Nome        string             `bson:"nome" json:"nome"`
	Email       string             `bson:"email" json:"email"`
	Sexo        string             `bson:"sexo" json:"sexo"`
	Telefone    string             `bson:"telefone" json:"telefone"`
	Tipo        string             `bson:"tipo" json:"tipo"`
	Documento   string             `bson:"cpf_cnpj" json:"cpf_cnpj"`
	Enabled     bool               `bson:"enabled" json:"enabled"`
	CreatedAt   string             `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt   string             `bson:"updated_at" json:"updated_at,omitempty"`
}

func (c Cliente) ClienteConvet() string {
	data, err := json.Marshal(c)

	if err != nil {
		logger.Error("error to convert Client to JSON", err)

		return ""
	}

	return string(data)
}

type FilterCliente struct {
	Nome      string             `json:"nome"`
	IDUsuario primitive.ObjectID `bson:"user_id " json:"id_usr"`
	Documento string             `bson:"documento" json:"documento"`
	Enabled   string             `json:"enabled"`
}

func NewCliente(cliente_request Cliente) *Cliente {
	return &Cliente{
		ID:        primitive.NewObjectID(),
		DataType:  "cliente",
		IDUsuario: cliente_request.IDUsuario,
		Nome:      validation.CareString(cliente_request.Nome),
		Enabled:   true,
		CreatedAt: time.Now().String(),
	}
}
