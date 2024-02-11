package model

import (
	"encoding/json"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserInterface interface {
	String() string
}

type Usuario struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Nome      string             `bson:"name" json:"name"`
	Email     string             `bson:"email" json:"email"`
	Senha     string             `bson:"senha" json:"senha"`
	Enable    bool               `bson:"enable" json:"enable"`
	Role      string             `bson:"role" json:"role"`
	CreatedAt string             `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt string             `bson:"updated_at" json:"updated_at,omitempty"`
}

type FilterUsuario struct {
	Nome   string
	Email  string
	Enable string
}

func (u *Usuario) String() string {
	data, err := json.Marshal(u)

	if err != nil {
		log.Println("Error convert User to JSON")
		log.Println(err.Error())
		return ""
	}

	return string(data)
}

func (u *Usuario) CheckPassword(senha string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Senha), []byte(senha))
	return err == nil

}

func (u *Usuario) ValidarRoler(role string) bool {
	roleMap := map[string]string{
		"cliente":    "cliente",
		"fornecedor": "fornecedor",
		"parceiro":   "parceiro",
		"admin":      "admin",
	}
	_, existe := roleMap[role]
	return existe
}

func NewUsuario(nome, senha, email, role string) (*Usuario, error) {
	dt := time.Now().Format(time.RFC3339)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Erro to SetPassWord", err.Error())
		return nil, err
	}

	tmp_user := &Usuario{
		Nome:      nome,
		Senha:     string(hashedPassword),
		Email:     email,
		Enable:    false,
		Role:      role,
		CreatedAt: dt,
		UpdatedAt: dt,
	}

	return tmp_user, nil
}
