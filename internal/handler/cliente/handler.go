package cliente

import (
	"encoding/json"
	"strconv"

	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/katana/back-end/orcafacil-go/internal/config/logger"
	"github.com/katana/back-end/orcafacil-go/pkg/service/cliente"
	"github.com/katana/back-end/orcafacil-go/pkg/service/validation"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/katana/back-end/orcafacil-go/pkg/model"
)

func createCliente(service cliente.ClienteServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Cliente := &model.Cliente{}

		err := json.NewDecoder(r.Body).Decode(&Cliente)

		if err != nil {
			logger.Error("error decoding request body", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		switch Cliente.Tipo {
		case "Fisica":
			if !validarCPF(w, *Cliente) {
				return
			}
		case "Juridica":
			if !validarCNPJ(w, *Cliente) {
				return
			}
		default:
			http.Error(w, "Tipo de cliente inválido", http.StatusBadRequest)
			return
		}

		if service.GetByDocumento(r.Context(), Cliente.Documento) {
			http.Error(w, "Documento já cadastrado", http.StatusBadRequest)
			return
		}
		_, err = service.Create(r.Context(), *Cliente)
		if err != nil {
			logger.Error("erro ao acessar a camada de service do mpg", err)
			http.Error(w, "Error ou salvar Cliente"+err.Error(), http.StatusInternalServerError)
			return
		}

		type Response struct {
			Message string `json:"message"`
		}

		// Crie uma instância da estrutura com a mensagem desejada.
		msg := Response{
			Message: "Dados gravados com sucesso",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(msg)
	}
}

func updateCliente(service cliente.ClienteServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		idp := chi.URLParam(r, "id")
		logger.Info("PEGANDO O PARAMENTRO")

		_, err := service.GetByID(r.Context(), idp)
		if err != nil {
			http.Error(w, "Cliente nao encontrada", http.StatusNotFound)
			return
		}

		mpg := &model.Cliente{}
		nome := chi.URLParam(r, "nome")
		logger.Info("PEGANDO O NOME")
		logger.Info(nome)
		if nome == "" {
			http.Error(w, "o Nome do curso e obrigatório", http.StatusBadRequest)
			return
		}

		mpg.Nome = nome
		id, err := primitive.ObjectIDFromHex(idp)
		if err != nil {
			http.Error(w, "erro ao converter id", http.StatusBadRequest)

			return
		}

		mpg.ID = id
		_, err = service.Update(r.Context(), idp, *&mpg)
		if err != nil {
			logger.Error("erro ao acessar a camada de service do mpg no upd", err)
			http.Error(w, "Error ao atualizar meio de pagamento", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"MSG": "Success", "codigo": 1})
	}
}

func getByIdCliente(service cliente.ClienteServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		idp := chi.URLParam(r, "id")
		logger.Info("PEGANDO O PARAMENTRO NA CONSULTA")
		result, err := service.GetByID(r.Context(), idp)
		if err != nil {
			logger.Error("erro ao acessar a camada de service da Cliente no por id", err)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"MSG": "Cliente não encontrada", "codigo": 404}`))
			return
		}

		err = json.NewEncoder(w).Encode(result)
		if err != nil {
			logger.Error("erro ao converter em json", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"MSG": "Error to parse Bot to JSON", "codigo": 500}`))
			return
		}
	}
}

func getAllCliente(service cliente.ClienteServiceInterface) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		filters := model.FilterCliente{
			Nome:    chi.URLParam(r, "nome"),
			Enabled: chi.URLParam(r, "enable"),
		}

		limit, _ := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
		page, _ := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)

		result, err := service.GetAll(r.Context(), filters, limit, page)
		if err != nil {
			logger.Error("erro ao acessar a camada de service do mpg no upd", err)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"MSG": "User not found", "codigo": 404}`))
			return
		}

		// Configurando o cabeçalho para resposta JSON usando o middleware
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		// Escrevendo a resposta JSON
		err = json.NewEncoder(w).Encode(result)
		if err != nil {
			logger.Error("erro ao converter para json", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"MSG": "Error to parse User to JSON", "codigo": 500}`))
			return
		}
	})
}

func validarCPF(w http.ResponseWriter, cli model.Cliente) bool {
	if !validation.IsCPFValid(cli.Documento) {
		http.Error(w, "CPF inválido", http.StatusBadRequest)
		return false
	}
	return true
}

func validarCNPJ(w http.ResponseWriter, cli model.Cliente) bool {
	if !validation.IsCNPJValid(cli.Documento) {
		http.Error(w, "CNPJ inválido", http.StatusBadRequest)
		return false
	}
	return true
}
