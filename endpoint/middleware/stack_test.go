package middleware

import (
	"testing"

	"github.com/pakkasys/fluidapi/core/api"
	"github.com/stretchr/testify/assert"
)

func TestInsertAfterID_Success(t *testing.T) {
	// Define initial middleware stack
	mw1 := api.MiddlewareWrapper{ID: "auth"}
	mw2 := api.MiddlewareWrapper{ID: "logging"}
	mwStack := Stack{mw1, mw2}

	// Define new middleware to insert
	newMiddleware := api.MiddlewareWrapper{ID: "metrics"}

	// Insert after 'auth'
	inserted := mwStack.InsertAfterID("auth", newMiddleware)

	// Assertions
	assert.True(t, inserted, "Middleware should be inserted")
	assert.Equal(t, 3, len(mwStack), "Middleware stack should have 3 elements")
	assert.Equal(t, "auth", mwStack[0].ID)
	assert.Equal(t, "metrics", mwStack[1].ID, "Middleware not in 2nd position")
	assert.Equal(t, "logging", mwStack[2].ID)
}

func TestInsertAfterID_AppendToEnd(t *testing.T) {
	// Define initial middleware stack
	mw1 := api.MiddlewareWrapper{ID: "auth"}
	mw2 := api.MiddlewareWrapper{ID: "logging"}
	mwStack := Stack{mw1, mw2}

	// Define new middleware to insert
	newMiddleware := api.MiddlewareWrapper{ID: "metrics"}

	// Insert after 'logging'
	inserted := mwStack.InsertAfterID("logging", newMiddleware)

	// Assertions
	assert.True(t, inserted, "Middleware should be inserted")
	assert.Equal(t, 3, len(mwStack), "Middleware stack should have 3 elements")
	assert.Equal(t, "auth", mwStack[0].ID)
	assert.Equal(t, "logging", mwStack[1].ID)
	assert.Equal(t, "metrics", mwStack[2].ID, "New middleware not in the end")
}

func TestInsertAfterID_IDNotFound(t *testing.T) {
	// Define initial middleware stack
	mw1 := api.MiddlewareWrapper{ID: "auth"}
	mw2 := api.MiddlewareWrapper{ID: "logging"}
	mwStack := Stack{mw1, mw2}

	// Define new middleware to insert
	newMiddleware := api.MiddlewareWrapper{ID: "metrics"}

	// Try inserting after a non-existent middleware ID
	inserted := mwStack.InsertAfterID("non-existent-id", newMiddleware)

	// Assertions
	assert.False(t, inserted, "Middleware should not be inserted")
	assert.Equal(t, 2, len(mwStack), "Middleware stack should have 2 elements")
	assert.Equal(t, "auth", mwStack[0].ID)
	assert.Equal(t, "logging", mwStack[1].ID)
}
