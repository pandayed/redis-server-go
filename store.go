package main

import "github.com/pandayed/redis-server-go/store"

var storeInstance store.DataStore

// newStore preserves the constructor used by existing tests.
func newStore() store.DataStore {
	return store.New()
}
