package face

import (
	"github.com/andyzhou/tinymeili/conf"
	"github.com/meilisearch/meilisearch-go"
)

/*
 * one index face
 */

//face info
type Index struct {
	client *meilisearch.Client //reference
	indexConf *conf.IndexConf
	workers int
	index *meilisearch.Index
	doc *Doc
}

//construct
func NewIndex(
	client *meilisearch.Client,
	indexConf *conf.IndexConf,
	workers int) *Index {
	this := &Index{
		client: client,
		indexConf: indexConf,
		workers: workers,
	}
	this.interInit()
	return this
}

//quit
func (f *Index) Quit() {
	f.doc.Quit()
}

//get doc face
func (f *Index) GetDoc() *Doc {
	return f.doc
}

//get status info
func (f *Index) GetStatus() (*meilisearch.StatsIndex, error) {
	return f.index.GetStats()
}

//update filterable fields
func (f *Index) UpdateFilterableAttributes(fields []string) error {
	_, err := f.index.UpdateFilterableAttributes(&fields)
	return err
}

//inter init
func (f *Index) interInit() {
	//init index
	index := f.client.Index(f.indexConf.IndexName)

	//update filterable fields
	if f.indexConf.FilterableFields != nil && len(f.indexConf.FilterableFields) > 0 {
		_, err := index.UpdateFilterableAttributes(&f.indexConf.FilterableFields)
		if err != nil {
			panic(any(err))
		}
	}

	//create index
	indexCfg := &meilisearch.IndexConfig{
		Uid: f.indexConf.IndexName,
		PrimaryKey: f.indexConf.PrimaryKey,
	}
	_, err := f.client.CreateIndex(indexCfg)
	if err != nil {
		panic(any(err))
	}

	//init doc obj
	f.index = index
	f.doc = NewDoc(index, f.workers)
}