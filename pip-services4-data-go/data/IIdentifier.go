package data

type IIdentifier[K any] interface {
	Empty() bool
	Equals(K) bool
}
