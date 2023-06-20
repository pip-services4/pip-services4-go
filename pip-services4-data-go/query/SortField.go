package query

// SortField Defines a field name and order used to sort query results.
//	see SortParams
//
//	Example:
//		filter := NewFilterParamsFromTuples("type", "Type1")
//		paging := NewPagingParams(0, 100)
//		sorting := NewSortingParams(NewSortField("create_time", true))
//
//		err, page = myDataClient.GetDataByFilter(context.Background(), filter, paging, sorting)
type SortField struct {
	Name      string `json:"name"`
	Ascending bool   `json:"ascending"`
}

// NewEmptySortField creates a new empty instance.
// Returns SortField
func NewEmptySortField() SortField {
	return SortField{}
}

// NewSortField creates a new instance and assigns its values.
//	Parameters:
//		- name string the name of the field to sort by.
//		- ascending: bool true to sort in ascending order, and false to sort in descending order.
//	Returns: SortField
func NewSortField(name string, ascending bool) SortField {
	return SortField{
		Name:      name,
		Ascending: ascending,
	}
}
