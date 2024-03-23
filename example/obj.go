package main

type TestDoc struct {
	Id    int64    `json:"id"`
	Title string   `json:"title"`
	Tags  []string `json:"tags"`
}

//construct
func NewTestDoc() *TestDoc {
	this := &TestDoc{
		Tags: []string{},
	}
	return this
}
