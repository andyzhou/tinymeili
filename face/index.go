package face

import (
	"errors"
	"fmt"
	"github.com/andyzhou/tinymeili/conf"
	"github.com/meilisearch/meilisearch-go"
	"log"
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
	//check
	if fields == nil || len(fields) <= 0 {
		return errors.New("invalid parameter")
	}
	//update fields
	task, err := f.index.UpdateFilterableAttributes(&fields)
	if err != nil {
		return err
	}

	//wait for task
	finalTask, _ := f.client.WaitForTask(task.TaskUID)
	if finalTask.Status != "succeeded" {
		return fmt.Errorf(finalTask.Error.Code)
	}
	return nil
}

//update primary key
func (f *Index) UpdatePrimaryKey(key string) error {
	//check
	if key == "" {
		return errors.New("invalid parameter")
	}
	//update key
	task, err := f.index.UpdateIndex(key)
	if err != nil {
		return err
	}

	//wait for task
	finalTask, _ := f.client.WaitForTask(task.TaskUID)
	if finalTask.Status != "succeeded" {
		return fmt.Errorf(finalTask.Error.Code)
	}
	return nil
}

//inter init
func (f *Index) interInit() {
	//init index config
	indexCfg := &meilisearch.IndexConfig{
		Uid: f.indexConf.IndexName,
		PrimaryKey: f.indexConf.PrimaryKey,
	}

	//create index
	task, err := f.client.CreateIndex(indexCfg)
	if err != nil {
		panic(any(err))
	}

	//wait for task
	finalTask, _ := f.client.WaitForTask(task.TaskUID)
	if finalTask.Status != "succeeded" && finalTask.Error.Code != "index_already_exists" {
		err = fmt.Errorf("create index failed, err:%v", finalTask.Status)
		log.Printf("init index %v failed, err:%v\n", f.indexConf.IndexName, err.Error())
		panic(any(err))
	}
	index, _ := f.client.GetIndex(f.indexConf.IndexName)

	//set index primary key
	index.PrimaryKey = f.indexConf.PrimaryKey

	//sync index obj
	f.index = index

	//set filterable fields
	if f.indexConf.FilterableFields != nil && len(f.indexConf.FilterableFields) > 0 {
		err = f.UpdateFilterableAttributes(f.indexConf.FilterableFields)
		if err != nil {
			panic(any(err))
		}
	}

	//init doc obj
	f.doc = NewDoc(f.client, f.index, f.workers)
}