package query

import (
	//"errors"
	"github.com/pip-services4/pip-services4-go/pip-services4-commons-go/convert"
)

// DataPage is a transfer object that is used to pass results of paginated queries. It contains items of retrieved
// page and optional total number of items.
// Most often this object type is used to send responses to paginated queries.
// Pagination parameters are defined by PagingParams object. The skip parameter in the PagingParams
// there means how many items to skip. The takes parameter sets number of items to return in the page.
// And the optional total parameter tells to return total number of items in the query.
// Remember: not all implementations support the total parameter because its generation may lead to severe
// performance implications.
//
//		see PagingParams
//
//		Example:
//			page, err := myDataClient.GetDataByFilter(
//				context.Background(),
//	     	"123",
//	     	FilterParams.fromTuples("completed": true),
//	     	NewPagingParams(0, 100, true),
//	 	);
//
//	 	if err != nil {
//	 		panic(err)
//	 	}
//	 	for item range page.Data {
//	         fmt.Println(item);
//	     }
//	 );
type DataPage[T any] struct {
	Total int `json:"total"`
	Data  []T `json:"data"`
}

const EmptyTotalValue int = -1

// NewEmptyDataPage creates a new empty instance of data page.
//
//	Returns: *DataPage
func NewEmptyDataPage[T any]() *DataPage[T] {
	return &DataPage[T]{
		Total: EmptyTotalValue,
		Data:  nil,
	}
}

// NewDataPage creates a new instance of data page and assigns its values.
//
//	Parameters:
//		- value data a list of items from the retrieved page.
//		- total int
//	Returns: *DataPage
func NewDataPage[T any](data []T, total int) *DataPage[T] {
	dataPage := DataPage[T]{
		Data:  data,
		Total: total,
	}

	return &dataPage
}

// HasData method check if data exists
func (d *DataPage[T]) HasData() bool {
	return len(d.Data) > 0
}

// HasTotal method check if total exists and valid
func (d *DataPage[T]) HasTotal() bool {
	return d.Total >= len(d.Data)
}

func (d DataPage[T]) MarshalJSON() ([]byte, error) {
	result := map[string]any{
		"data": d.Data,
	}
	if d.HasTotal() {
		result["total"] = d.Total
	}
	buf, err := convert.JsonConverter.ToJson(result)
	return []byte(buf), err
}
