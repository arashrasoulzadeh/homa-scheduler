package invokers

import (
	"context"
	"net/http"

	"github.com/arashrasoulzadeh/homa-scheduler/providers"
	"github.com/arashrasoulzadeh/homa-scheduler/router"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func RunHttpServer(lc fx.Lifecycle, logger *zap.SugaredLogger, data providers.Data) *chi.Mux {
	r := chi.NewRouter()

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			port := ":3000"
			r.Use(middleware.Logger)
			router.Init(r, logger, data)
			logger.Infow("http server is running", "port", port)
			go http.ListenAndServe(port, r)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
	return r
}
