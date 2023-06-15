package data

type IMap interface {
	Get(key string) (any, bool)
	Put(key string, value any)
	Remove(key string)
	Contains(key string) bool
	Len() int
}
