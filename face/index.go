package face

import (
	"github.com/meilisearch/meilisearch-go"
)

/*
 * one index face
 */

//face info
type Index struct {
	indexTag string
	client *meilisearch.Client //reference
	index *meilisearch.Index //reference
	doc *Doc
}

//construct
func NewIndex(indexTag string, index *meilisearch.Index) *Index {
	this := &Index{
		indexTag: indexTag,
		index: index,
		doc: NewDoc(index),
	}
	return this
}

//get doc face
func (f *Index) GetDoc() *Doc {
	return f.doc
}

//update filterable fields
func (f *Index) UpdateFilterableAttributes(fields ...string) error {
	_, err := f.index.UpdateFilterableAttributes(&fields)
	return err
}