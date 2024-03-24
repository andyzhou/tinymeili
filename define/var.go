package define

//query para
type (
	QueryPara struct {
		Key string
		Filter interface{}
		Sort []string
		Facets []string
		Page, PageSize int
	}
)
