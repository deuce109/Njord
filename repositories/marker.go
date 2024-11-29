package repositories

import (
	"deuce109/7dtd-map-server/v2/models"

	"github.com/globalsign/mgo/bson"
)

type MarkerRepositoryHandler interface {
	Upsert(id string, marker models.Marker) error
	GetByUserId(userId string) []models.Marker
	RemoveById(id string) error
	GetWorldMarkers() []models.Marker
}

type collection interface {
	Find(selector interface{}) cursor
	RemoveId(id interface{}) error
	Insert(insertion interface{}) error
	Upsert(selector interface{}, update interface{}) (interface{}, error)
}

type iter interface {
	Next(v interface{}) bool
}

type cursor interface {
	Iter() iter
}

type MarkerRepository struct {
	collection collection
}

func (r *MarkerRepository) Upsert(id string, marker models.Marker) error {
	_, err := r.collection.Upsert(bson.M{
		"_id": id,
	}, marker)

	return err
}

func (r *MarkerRepository) GetByUserId(userId string) []models.Marker {
	cur := r.collection.Find(bson.M{
		"userId": userId,
	})
	return getResultsFromCursor[models.Marker](cur)
}

func (r *MarkerRepository) RemoveById(id string) error {
	return r.collection.RemoveId(id)
}

func (r *MarkerRepository) GetWorldMarkers() []models.Marker {
	cur := r.collection.Find(bson.M{
		"global": true,
	})
	return getResultsFromCursor[models.Marker](cur)

}
