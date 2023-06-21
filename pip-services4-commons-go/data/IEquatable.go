package data

type IEquatable[T any] interface {
	Equals(T) bool
}
