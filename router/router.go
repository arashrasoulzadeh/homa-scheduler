package router

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/arashrasoulzadeh/homa-scheduler/models"
	"github.com/arashrasoulzadeh/homa-scheduler/providers"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/datatypes"
)

func Init(r *chi.Mux, logger *zap.SugaredLogger, data providers.Data) {
	scheduleRoute(r, data)
	singleScheduledItemRoute(r, data)
}
func scheduleRoute(r *chi.Mux, data providers.Data) {
	r.Post("/api/v1/schedule", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("accept", "applicatin/json")
		w.Header().Add("content-type", "applicatin/json")
		item := models.Command{
			ID:        uuid.New(),
			Command:   "resize",
			Args:      datatypes.JSONMap{"image": "http"},
			Status:    "pending",
			CreatedAt: time.Now(),
			Channel:   "general",
		}
		item.MarkAsDev()
		data.Bus <- item
		data.Connection.Save(item)
		json.NewEncoder(w).Encode(item)
	})
}
func singleScheduledItemRoute(r *chi.Mux, data providers.Data) {
	r.Post("/api/v1/schedule/{uuid}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("accept", "applicatin/json")
		w.Header().Add("content-type", "applicatin/json")
		id := chi.URLParam(r, "uuid")
		item := models.Command{}
		if err := data.Connection.Where("id = ?", id).First(&item).Error; err != nil {
			w.WriteHeader(404)
			return
		}
		json.NewEncoder(w).Encode(item)
	})
}
