package main

import (
	"log"
	"net/http"

	"github.com/katana/fortuna/backend-go/internal/config"
	"github.com/katana/fortuna/backend-go/internal/config/logger"

	hand_cliente "github.com/katana/fortuna/backend-go/internal/handler/cliente"

	hand_sorteio "github.com/katana/fortuna/backend-go/internal/handler/sorteio"
	hand_usr "github.com/katana/fortuna/backend-go/internal/handler/user"

	"github.com/katana/fortuna/backend-go/pkg/adapter/mongodb"

	"github.com/katana/fortuna/backend-go/pkg/server"

	service_usr "github.com/katana/fortuna/backend-go/pkg/service/user"

	service_cliente "github.com/katana/fortuna/backend-go/pkg/service/cliente"
	service_sorteio "github.com/katana/fortuna/backend-go/pkg/service/sorteio"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

var (
	VERSION = "0.1.0-dev"
	COMMIT  = "ABCDEFG-dev"
)

func main() {
	// fila := []rabbitmq.Fila{
	// 	{
	// 		Name:    "QUEUE_PRDS_PARA_COTACAO",
	// 		Durable: true,
	// 	},
	// }
	logger.Info("start Application Fortuna fast API")
	conf := config.NewConfig()

	mogDbConn := mongodb.New(conf)
	//rbtMQConn := rabbitmq.NewRabbitMQ(fila, conf)
	//rdisConn := redisdb.NewRedisClient(conf)
	usr_service := service_usr.NewUsuarioservice(mogDbConn)

	cli_service := service_cliente.NewClienteervice(mogDbConn)

	sor_service := service_sorteio.NewSorteioService(mogDbConn)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.WithValue("jwt", conf.TokenAuth))
	r.Use(middleware.WithValue("JWTTokenExp", conf.JWTTokenExp))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Get("/", healthcheck)
	hand_usr.RegisterUsuarioAPIHandlers(r, usr_service)

	hand_cliente.RegisterClientePIHandlers(r, cli_service)
	hand_sorteio.RegisterSorteioPIHandlers(r, sor_service)

	srv := server.NewHTTPServer(r, conf)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	log.Printf("Server Run on [Port: %s], [Mode: %s], [Version: %s], [Commit: %s]", conf.PORT, conf.Mode, VERSION, COMMIT)

	select {}
}

func healthcheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"MSG": "Server Ok", "codigo": 200}`))
}
