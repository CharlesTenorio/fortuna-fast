package user

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/katana/fortuna/backend-go/pkg/service/user"
)

func RegisterUsuarioAPIHandlers(r chi.Router, service user.UserServiceInterface) {
	r.Route("/api/v1/usuario", func(r chi.Router) {
		r.Post("/add", createUser(service))
		r.Post("/getjwt", PegarJwt(service))
		r.Put("/update/{id}/{nome}", updateUser(service))
		r.Get("/getbyid/{id}", getByIdUsuario(service)) // Adicionado a barra no início
		r.Get("/all", func(w http.ResponseWriter, r *http.Request) {
			handler := getAllUsuario(service)
			handler.ServeHTTP(w, r)
		})
	})
}
