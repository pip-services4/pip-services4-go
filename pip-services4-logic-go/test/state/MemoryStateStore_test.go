package test_state

import (
	"testing"

	"github.com/pip-services4/pip-services4-go/pip-services4-logic-go/state"
)

func newMemoryStateStoreFixture() *StateStoreFixture {
	stater := state.NewEmptyMemoryStateStore[any]()
	fixture := NewStateStoreFixture(stater)
	return fixture
}

func TestSaveAndLoad(t *testing.T) {
	fixture := newMemoryStateStoreFixture()
	fixture.TestSaveAndLoad(t)
}

func TestDelete(t *testing.T) {
	fixture := newMemoryStateStoreFixture()
	fixture.TestDelete(t)
}
