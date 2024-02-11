package user

import (
	"encoding/json"
	"strconv"
	"time"

	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
	"github.com/katana/fortuna/backend-go/internal/config/logger"
	"github.com/katana/fortuna/backend-go/internal/dto"
	"github.com/katana/fortuna/backend-go/pkg/model"
	"github.com/katana/fortuna/backend-go/pkg/service/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func createUser(service user.UserServiceInterface) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user := &model.Usuario{}

		err := json.NewDecoder(r.Body).Decode(user)

		if err != nil {
			logger.Error("erro ao converte para json", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"MSG": "Error to parse User to JSON", "codigo": 500}`))
			return
		}

		if !user.ValidarRoler(user.Role) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"MSG": "Roler  invalida", "codigo": 400}`))
			return

		}

		userExist, err := service.GetByEmail(r.Context(), user.Email)

		if userExist != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"MSG": "User with username already exists", "codigo": 400}`))
			return
		}

		result, err := service.Create(r.Context(), user)
		if err != nil {
			logger.Error("erro ao criar usuario", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"MSG": "Error to Insert the User", "codigo": 500}`))
			return
		}

		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(result)
		if err != nil {
			logger.Error("erro ao converte esultada para json", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"MSG": "Error to parse User to JSON", "codigo": 500}`))
			return
		}
	})
}

func updateUser(service user.UserServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userToChange := &model.Usuario{}
		idp := chi.URLParam(r, "id")
		err := json.NewDecoder(r.Body).Decode(userToChange)

		if err != nil {
			logger.Error("error decoding request body", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		_, err = service.GetByID(r.Context(), idp)
		if err != nil {
			http.Error(w, "suario não encontrado", http.StatusNotFound)
			return
		}

		if !userToChange.ValidarRoler(userToChange.Role) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"MSG": "Roler  invalida", "codigo": 400}`))
			return

		}

		id, err := primitive.ObjectIDFromHex(idp)
		if err != nil {
			http.Error(w, "erro ao converter id", http.StatusBadRequest)

			return
		}

		userToChange.ID = id
		_, err = service.Update(r.Context(), idp, *&userToChange)
		if err != nil {
			logger.Error("erro ao acessar a camada de service do prd no upd", err)
			http.Error(w, "Error ao atualizar meio de pagamento", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"MSG": "Success", "codigo": 1})
	}
}

func getByIdUsuario(service user.UserServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		idp := chi.URLParam(r, "id")
		logger.Info("PEGANDO O PARAMENTRO NA CONSULTA")
		result, err := service.GetByID(r.Context(), idp)
		if err != nil {
			logger.Error("erro ao acessar a camada de service do prd no por id", err)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"MSG": "Meio de pagamento não encontrado", "codigo": 404}`))
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

func getAllUsuario(service user.UserServiceInterface) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		filters := model.FilterUsuario{
			Nome:   chi.URLParam(r, "nome"),
			Enable: chi.URLParam(r, "enable"),
		}

		limit, _ := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
		page, _ := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)

		result, err := service.GetAll(r.Context(), filters, limit, page)
		if err != nil {
			logger.Error("erro ao acessar a camada de service do prd no upd", err)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"MSG": "User not found", "codigo": 404}`))
			return
		}
		err = json.NewEncoder(w).Encode(result)
		if err != nil {
			logger.Error("erro ao converto para json", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"MSG": "Error to parse User to JSON", "codigo": 500}`))
			return
		}
	})
}

func PegarJwt(service user.UserServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		jwt := r.Context().Value("jwt").(*jwtauth.JWTAuth)

		jwtExpiresIn := r.Context().Value("JWTTokenExp").(int)
		logger.Info("CRIOU O JWTEXPERI")

		var user dto.GetJwtInput
		err := json.NewDecoder(r.Body).Decode(&user)

		if err != nil {
			logger.Error("erro ao converte para json o toke user", err)
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		userExist, err := service.GetByEmail(r.Context(), user.Email)
		if err != nil {

			logger.Error("erro ao localizar user por email", err)
			w.WriteHeader(http.StatusBadRequest)
		}

		if userExist == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"MSG": "Usuáiro  não cadastrado", "codigo": 400}`))
			return
		}

		if !userExist.CheckPassword(user.Senha) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"MSG": "Usuáiro não autorizado ", "codigo": 400}`))
			return
		}

		// Adiciona a informação da Role ao mapa
		tokenClaims := map[string]interface{}{
			"sub":  userExist.ID.String(),
			"exp":  time.Now().Add(time.Second * time.Duration(jwtExpiresIn)).Unix(),
			"role": userExist.Role,
		}

		// Gera o token com as informações incluídas
		_, tokenString, _ := jwt.Encode(tokenClaims)
		accessToken := dto.GetJWTOutput{AccessToken: tokenString}

		/*_, tokenString, _ := jwt.Encode(map[string]interface{}{
			"sub": userExist.ID.String(),
			"exp": time.Now().Add(time.Second * time.Duration(jwtExpiresIn)).Unix(),
		})*/

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(accessToken)

	}
}
