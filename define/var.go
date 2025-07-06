package define

//query para
type (
	QueryPara struct {
		Key                string
		AttributesToSearch []string
		Distinct           string      //distinct field
		Filter             interface{} //filter condition
		Sort               []string
		Facets             []string //agg fields
		Page, PageSize     int
	}
)
