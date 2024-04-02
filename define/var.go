package define

//query para
type (
	QueryPara struct {
		Key                string
		AttributesToSearch []string
		Filter             interface{}
		Sort               []string
		Facets             []string
		Page, PageSize     int
	}
)
