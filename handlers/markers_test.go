package handlers

import (
	"bytes"
	"deuce109/7dtd-map-server/v2/models"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testRepo[T interface{}] struct {
	FailRemove bool
	FailUpsert bool
}

func (r *testRepo[T]) GetWorldMarkers() []T {
	return []T{}
}

func (r *testRepo[T]) GetByUserId(_ string) []T {
	return []T{}

}

func (r *testRepo[T]) RemoveById(_ string) error {
	if r.FailRemove {
		return errors.New("marker not found")
	}
	return nil
}

func (r *testRepo[T]) Upsert(_ string, _ T) error {
	if r.FailUpsert {
		return errors.New("unable to insert marker")
	}
	return nil
}

var repo = &testRepo[models.Marker]{
	FailRemove: false,
	FailUpsert: false,
}

var handler = &MarkerHandler{
	Repository: repo,
}

func TestMain(m *testing.M) {
	m.Run()
}

func TestGetMarkers(t *testing.T) {

	w := httptest.NewRecorder()

	req := httptest.NewRequest("GET", "/markers", nil)

	req.Header.Set("Game-Id", "id")

	handler.GetMarkers(w, req)
	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestDeleteMarker(t *testing.T) {
	w := httptest.NewRecorder()

	req := httptest.NewRequest("DELETE", "/markers/123", nil)

	req.SetPathValue("id", "132")

	handler.DeleteMarker(w, req)
	if w.Code != 204 {
		t.Errorf("Expected status 204, got %d", w.Code)
	}
}

func TestDeleteMarkerWithFailDeleteSet(t *testing.T) {
	w := httptest.NewRecorder()

	req := httptest.NewRequest("DELETE", "/markers/123", nil)

	req.SetPathValue("id", "132")

	repo.FailRemove = true

	handler.DeleteMarker(w, req)
	if w.Code != 500 {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

func TestUpdateMarker(t *testing.T) {

	json, _ := models.ToJson(models.DefaultMarker)

	body := new(bytes.Buffer)

	body.Write(json)

	req := httptest.NewRequest("GET", "/markers", body)

	w := httptest.NewRecorder()

	handler.UpsertMarker(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

}

func TestUpdateMarkerWithPoorlyFormattedBody(t *testing.T) {

	json := []byte{0}

	body := new(bytes.Buffer)

	body.Write(json)

	req := httptest.NewRequest("GET", "/markers", body)

	w := httptest.NewRecorder()

	repo.FailUpsert = true

	handler.UpsertMarker(w, req)

	if w.Code != 500 {
		t.Errorf("Expected status %d, got %d", 500, w.Code)
	}

}

func TestUpdateMarkerWithFailInsertSet(t *testing.T) {

	json, _ := models.ToJson(models.DefaultMarker)

	body := new(bytes.Buffer)

	body.Write(json)

	req := httptest.NewRequest("GET", "/markers", body)

	w := httptest.NewRecorder()

	repo.FailUpsert = true

	handler.UpsertMarker(w, req)

	if w.Code != 500 {
		t.Errorf("Expected status %d, got %d", 500, w.Code)
	}

}
