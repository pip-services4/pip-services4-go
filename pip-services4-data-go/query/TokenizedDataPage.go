package query

// TokenizedDataPage is a data transfer object that is used to pass results of paginated queries.
// It contains items of retrieved page and optional total number of items.
// Most often this object type is used to send responses to paginated queries.
// Pagination parameters are defined by TokenizedPagingParams object. The token parameter in the TokenizedPagingParams
// there determines a starting point for a new search. It is received in the TokenizedDataPage from the previous search.
// The takes parameter sets number of items to return to the page.
// And the optional total parameter tells to return total number of items in the query.
// Remember: not all implementations support the total parameter because its generation may lead to severe
// performance implications.
//	see TokenizedPagingParams
//
//	Example:
//		page, err := myDataClient.GetDataByFilter(
//			context.Background(),
//			"123",
//			NewFilterParamsFromTuples("completed": true),
//			NewTokenizedPagingParams("", 100, true),
//		);
//
//		if err != nil {
//			panic(err)
//		}
//		for item range page.Data {
//			fmt.Println(item);
//		}
type TokenizedDataPage[T any] struct {
	Token string `json:"token"`
	Data  []T    `json:"data"`
}

const EmptyTokenValue string = ""

// NewEmptyTokenizedDataPage creates a new empty instance of data page.
//	Returns: *TokenizedDataPage[T]
func NewEmptyTokenizedDataPage[T any]() *TokenizedDataPage[T] {
	return &TokenizedDataPage[T]{
		Token: EmptyTokenValue,
	}
}

// NewTokenizedDataPage creates a new instance of data page and assigns its values.
//	Parameters:
//		- token a token that defines a starting point for next search
//		- data []T a list of items from the retrieved page.
//	Returns: *TokenizedDataPage[T]
func NewTokenizedDataPage[T any](token string, data []T) *TokenizedDataPage[T] {
	return &TokenizedDataPage[T]{Token: token, Data: data}
}

// HasData method check if data exists
func (d *TokenizedDataPage[T]) HasData() bool {
	return len(d.Data) > 0
}

// HasToken method check if token exists
func (d *TokenizedDataPage[T]) HasToken() bool {
	return len(d.Token) > 0
}
