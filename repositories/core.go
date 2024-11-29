package repositories

func getResultsFromCursor[T interface{}](cur cursor) []T {

	ts := make([]T, 0)
	var t T
	iter := cur.Iter()
	for iter.Next(&t) {
		ts = append(ts, t)
	}

	return ts
}
