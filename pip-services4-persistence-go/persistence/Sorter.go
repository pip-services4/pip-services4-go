package persistence

//------------- Sorter -----------------------

// sorter is a helper class for sorting data in MemoryPersistence
// implements sort.Interface
type sorter[T any] struct {
	items    []T
	compFunc func(a, b T) bool
}

// Len calculate length
//	Returns: length of items array
func (s sorter[T]) Len() int {
	return len(s.items)
}

// Swap two items in array
//	Parameters:
//		- i,j int indexes of array for swap
func (s sorter[T]) Swap(i, j int) {
	s.items[i], s.items[j] = s.items[j], s.items[i]
}

// Less compare less function
//	Parameters:
//		- i,j int indexes of array for compare
// Returns: bool true if items[i] < items[j] and false otherwise
func (s sorter[T]) Less(i, j int) bool {
	if s.compFunc == nil {
		panic("Sort.Less Error compare function is nil!")
	}
	return s.compFunc(s.items[i], s.items[j])
}
