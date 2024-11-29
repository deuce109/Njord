package repositories

import (
	"deuce109/7dtd-map-server/v2/models"
	"deuce109/7dtd-map-server/v2/readers"
	"errors"
	"testing"
)

type testCollection struct {
	Parameters map[string]map[string]interface{}
}

type testCursor struct {
}

type testIter struct {
	callCounter int
}

func (c *testCollection) Insert(v interface{}) error {
	return nil
}

func (c *testCollection) RemoveId(v interface{}) error {
	return nil
}

func (c *testCollection) Upsert(selector interface{}, v interface{}) (interface{}, error) {

	if selector == nil {
		return nil, errors.New("selector should not be nil")
	} else if v == nil {
		return nil, errors.New("v should not be nil")
	}

	return nil, nil
}

func (i *testIter) Next(v interface{}) bool {
	i.callCounter += 1
	return i.callCounter == 1
}

func (c *testCursor) Iter() iter {
	return &testIter{callCounter: 0}
}

func (c *testCollection) Find(_ interface{}) cursor {
	return &testCursor{}
}

var repo MarkerRepository

func TestMain(t *testing.M) {
	collection := &testCollection{}
	repo = MarkerRepository{
		collection: collection,
	}
	readers.GetUserInfo("../test_data/good", "../test_data/good")
	t.Run()
}

func TestGetByUserId(t *testing.T) {
	markers := repo.GetByUserId("test")

	if markers[0].UserId != "" {
		t.Fail()
	}
}

func TestGetWorldMArkers(t *testing.T) {
	markers := repo.GetWorldMarkers()

	if markers[0].UserId != "" {
		t.Fail()
	}
}

func TestRemoveId(t *testing.T) {
	err := repo.RemoveById("testModel")

	if err != nil {
		t.Fatal(err)
	}
}

func TestUpdate(t *testing.T) {
	err := repo.Upsert("testModel", models.DefaultMarker)

	if err != nil {
		t.Fatal(err)
	}
}

func TestInsert(t *testing.T) {
	err := repo.Upsert("", models.DefaultMarker)

	if err != nil {
		t.Fatal(err)
	}
}
