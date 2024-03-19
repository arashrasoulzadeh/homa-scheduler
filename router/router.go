package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/arashrasoulzadeh/homa-scheduler/models"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func Init(db *gorm.DB) {
	r := chi.NewRouter()
	port := ":3000"
	r.Use(middleware.Logger)
	scheduleRoute(r, db)
	singleScheduledItemRoute(r, db)
	fmt.Println("listening on port " + port)
	http.ListenAndServe(port, r)

}
func scheduleRoute(r *chi.Mux, db *gorm.DB) {
	r.Post("/api/v1/schedule", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("accept", "applicatin/json")
		w.Header().Add("content-type", "applicatin/json")
		item := models.Command{
			Id:        uuid.New(),
			Command:   "resize",
			Args:      datatypes.JSONMap{"image": "http"},
			Status:    "pending",
			CreatedAt: time.Now(),
		}
		item.MarkAsDev()
		db.Save(item)
		json.NewEncoder(w).Encode(item)
	})
}
func singleScheduledItemRoute(r *chi.Mux, db *gorm.DB) {
	r.Post("/api/v1/schedule/{uuid}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("accept", "applicatin/json")
		w.Header().Add("content-type", "applicatin/json")
		id := chi.URLParam(r, "uuid")
		item := models.Command{}
		if err := db.Where("id = ?", id).First(&item).Error; err != nil {
			w.WriteHeader(404)
			return
		}
		json.NewEncoder(w).Encode(item)
	})
}
