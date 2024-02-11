package validation

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Address struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func GetAddressByCEP(cep string) (Address, error) {
	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return Address{}, err
	}

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return Address{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Address{}, fmt.Errorf("failed to get address for CEP: %s", cep)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Address{}, err
	}

	var address Address
	err = json.Unmarshal(body, &address)
	if err != nil {
		return Address{}, err
	}

	return address, nil
}
func ExtractNumbers(input string) string {
	re := regexp.MustCompile("[^0-9]+")
	return re.ReplaceAllString(input, "")
}

func IsCPFValid(cpf string) bool {
	// Remove caracteres não numéricos
	cpf = ExtractNumbers(cpf)

	// Verifica o tamanho do CPF
	if len(cpf) != 11 {
		return false
	}

	// Verifica se todos os dígitos são iguais (caso especial de CPF inválido)
	isAllDigitsEqual := true
	for i := 1; i < 11; i++ {
		if cpf[i] != cpf[0] {
			isAllDigitsEqual = false
			break
		}
	}
	if isAllDigitsEqual {
		return false
	}

	// Calcula o primeiro dígito verificador
	sum := 0
	for i := 0; i < 9; i++ {
		sum += int(cpf[i]-'0') * (10 - i)
	}
	remainder := sum % 11
	digit1 := 0
	if remainder >= 2 {
		digit1 = 11 - remainder
	}

	// Calcula o segundo dígito verificador
	sum = 0
	for i := 0; i < 10; i++ {
		sum += int(cpf[i]-'0') * (11 - i)
	}
	remainder = sum % 11
	digit2 := 0
	if remainder >= 2 {
		digit2 = 11 - remainder
	}

	// Verifica se os dígitos verificadores são válidos
	return digit1 == int(cpf[9]-'0') && digit2 == int(cpf[10]-'0')
}

func IsCNPJValid(cnpj string) bool {
	// Remove caracteres não numéricos
	cnpj = ExtractNumbers(cnpj)

	// Verifica o tamanho do CNPJ
	if len(cnpj) != 14 {
		return false
	}

	// Calcula o primeiro dígito verificador
	sum := 0
	for i := 0; i < 4; i++ {
		sum += int(cnpj[i]-'0') * (5 - i)
	}
	for i := 0; i < 8; i++ {
		sum += int(cnpj[i+4]-'0') * (9 - i)
	}
	remainder := sum % 11
	digit1 := 0
	if remainder >= 2 {
		digit1 = 11 - remainder
	}

	// Calcula o segundo dígito verificador
	sum = 0
	for i := 0; i < 5; i++ {
		sum += int(cnpj[i]-'0') * (6 - i)
	}
	for i := 0; i < 8; i++ {
		sum += int(cnpj[i+5]-'0') * (9 - i)
	}
	remainder = sum % 11
	digit2 := 0
	if remainder >= 2 {
		digit2 = 11 - remainder
	}

	return digit1 == int(cnpj[12]-'0') && digit2 == int(cnpj[13]-'0')
}

func CareString(word string) string {
	caser := cases.Title(language.BrazilianPortuguese)
	newWord := strings.TrimSpace(word)
	newWord = caser.String(newWord)
	return newWord
}
