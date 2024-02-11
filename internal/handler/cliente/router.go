package cliente

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/katana/back-end/orcafacil-go/pkg/service/cliente"
)

func RegisterClientePIHandlers(r chi.Router, service cliente.ClienteServiceInterface) {
	r.Route("/api/v1/cliente", func(r chi.Router) {
		r.Post("/add", createCliente(service))
		r.Put("/update/{id}/{nome}", updateCliente(service))
		r.Get("/getbyid/{id}", getByIdCliente(service))
		r.Get("/all", func(w http.ResponseWriter, r *http.Request) {
			handler := getAllCliente(service)
			handler.ServeHTTP(w, r)
		})
	})
}
