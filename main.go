package main

import (
	"crypto/rand"
	"deuce109/7dtd-map-server/v2/handlers"
	"deuce109/7dtd-map-server/v2/logging"
	"deuce109/7dtd-map-server/v2/middleware"
	"deuce109/7dtd-map-server/v2/readers"
	"deuce109/7dtd-map-server/v2/repositories"
	"deuce109/7dtd-map-server/v2/utils/crypto"
	"deuce109/7dtd-map-server/v2/validators"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"gopkg.in/validator.v2"
)

func main() {
	validator.SetValidationFunc("isUrl", validators.IsPath)

	random := &crypto.Random{
		Reader: rand.Reader,
	}

	worldName := os.Getenv("WORLD_NAME")
	saveName := os.Getenv("SAVE_NAME")

	playerXmlPath := filepath.Clean(filepath.Join("data", "saves", worldName, saveName))
	adminXmlPath := filepath.Clean(filepath.Join("data", "saves"))

	readers.GetUserInfo(playerXmlPath, adminXmlPath)

	secret, err := random.GetRandomSecret(32)

	if err != nil {
		logging.Fatal(err.Error())
	}

	mw := &middleware.JWTMiddleware{
		Secret: secret,
	}

	r := mux.NewRouter()
	r.HandleFunc("/login", mw.CreateJwt).Methods("POST").Headers("Game-Id")

	markerHandler := &handlers.MarkerHandler{
		Repository: &repositories.MarkerRepository{},
	}

	markerRouter := r.PathPrefix("/markers").Headers("Authorization").Subrouter()

	markerRouter.Use(mw.CheckJwt)
	markerRouter.HandleFunc("/", markerHandler.GetMarkers).Methods("GET").Headers("Game-Id")
	markerRouter.HandleFunc("/", markerHandler.DeleteMarker).Methods("DELETE")
	markerRouter.HandleFunc("/", markerHandler.UpsertMarker).Methods("POST")

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:9001",
	}

	logging.Fatal(srv.ListenAndServe().Error())
}
