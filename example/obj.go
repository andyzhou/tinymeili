package main

type TestDoc struct {
	Id       int64                  `json:"id"`
	Poster   int64                  `json:"poster"`
	Title    string                 `json:"title"`
	Property map[string]interface{} `json:"property"`
	Tags     []string               `json:"tags"`
	BaseJson
}

//construct
func NewTestDoc() *TestDoc {
	this := &TestDoc{
		Property: map[string]interface{}{},
		Tags: []string{},
	}
	return this
}
