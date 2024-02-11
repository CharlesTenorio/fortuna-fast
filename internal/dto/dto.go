package dto

type GetJwtInput struct {
	Email string `json:"email"`
	Senha string `json:"senha"`
	Role  string `json:"role"`
}

type GetJWTOutput struct {
	AccessToken string `json:"access_token"`
}

type FornecedoresEmPrd struct {
	ID         string  `json:"id"`
	PrecoVenda float64 `json:"preco_venda"`
}
