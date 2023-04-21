package main

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"

	"github.com/cambiahealth/janus-go/log"
	"github.com/spf13/viper"
)

type PrimeIntegrationService struct {
	logger *log.JanusLogger
	chiMux *chi.Mux
}

func main() {
	logger := log.New()
	srv := chi.NewRouter()

	ps := &PrimeIntegrationService{
		logger: logger,
		chiMux: srv,
	}

	ps.logger.Infof(context.Background(), "DYNAMODB_ENDPOINT %s", viper.GetString("DYNAMODB_ENDPOINT"))
	ps.logger.Infof(context.Background(), "PRIME_APIS_TABLE_NAME is %s", viper.GetString("PRIME_APIS_TABLE_NAME"))

	srv.HandleFunc("/test", ps.RowsTest)
	err := seedThoseAPIs()
	if err != nil {
		ps.logger.Panicf(context.Background(), "Unable to load seed data")
	}

	ps.logger.Info(context.Background(), "Starting prime-integration-service server")

	_ = http.ListenAndServe(":8080", srv)
}

func (ps *PrimeIntegrationService) RowsTest(w http.ResponseWriter, r *http.Request) {
	ps.logger.Info(r.Context(), "YO")
	primeAPIs, err := AllPrimeAPIs(r.Context())

	if err != nil {
		ps.logger.Error(r.Context(), "------------------- UNABLE TO FETCH ROWS -----------")
		w.WriteHeader(500)
		return
	}

	ej, err := json.Marshal(primeAPIs)
	_, _ = w.Write(ej)

}
