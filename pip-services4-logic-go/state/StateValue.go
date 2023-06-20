package state

// StateValue a data object that holds a retrieved state value with its key.
type StateValue[T any] struct {
	Key   string `json:"key" bson:"key"`     // A unique state key
	Value T      `json:"value" bson:"value"` // A stored state value
}
