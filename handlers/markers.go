package handlers

import (
	"deuce109/7dtd-map-server/v2/logging"
	"deuce109/7dtd-map-server/v2/models"
	"deuce109/7dtd-map-server/v2/repositories"
	"fmt"
	"net/http"
)

type MarkerHandler struct {
	Repository repositories.MarkerRepositoryHandler
}

func (h *MarkerHandler) GetMarkers(w http.ResponseWriter, r *http.Request) {

	worldMarkers := h.Repository.GetWorldMarkers()

	userMarkers := h.Repository.GetByUserId(r.Header["Game-Id"][0])

	markers := map[string]interface{}{
		"result": append(worldMarkers, userMarkers...),
	}

	jsonData, err := models.ToJson(markers)

	if err != nil {
		logging.Error(fmt.Sprintf("Error marshalling Markers in GetMarkers\n%s", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
}

func (h *MarkerHandler) DeleteMarker(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")

	err := h.Repository.RemoveById(id)

	if err != nil {
		logging.Error(fmt.Sprintf("Error attempting to remove id '%s'\n%s", id, err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *MarkerHandler) UpsertMarker(w http.ResponseWriter, r *http.Request) {
	marker, err := processBodyToItem[models.Marker](r)

	if err != nil {
		logging.Error(fmt.Sprintf("Error attempting to create marker\n%s", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.Repository.Upsert(string(marker.Id), marker)

	if err != nil {
		logging.Error(fmt.Sprintf("Error attempting to insert marker\n%s", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
