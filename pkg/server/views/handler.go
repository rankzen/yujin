package views

import (
	"context"
	"net/http"

	"github.com/rank1zen/yujin/pkg/database"
)

type Env interface {
	GetDatabase() database.DB
	// GetGolioClient() database.RiotClient
}

func NewHandler(ctx context.Context, router *http.ServeMux, env Env) (*http.ServeMux, error) {
	db := env.GetDatabase()

	handler := profilesHandler{db: db}

	router.HandleFunc("/", nil)
	router.HandleFunc("/profile/{puuid}", handler.profile())

	return router, nil
}
