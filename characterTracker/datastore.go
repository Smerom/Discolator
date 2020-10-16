package characterTracker

import (
	"context"

	"cloud.google.com/go/datastore"
)

type datastoreTracker struct {
	memoryTracker
}

func NewDatastoreTracker() Tracker {

	client, err := datastore.NewClient(context.Background(), "")
	if err != nil {

	}

	// load count
	_ = client

	// start periodic save

	return &datastoreTracker{}
}
