package face

import (
	"errors"
	"fmt"
	"log"

	"github.com/andyzhou/tinymeili/conf"
	"github.com/meilisearch/meilisearch-go"
)

/*
 * one index face
 * @author <AndyZhou>
 * @mail <diudiu8848@163.com>
 */

//face info
type Index struct {
	client    *meilisearch.Client //reference
	indexConf *conf.IndexConf
	index     *meilisearch.Index
	doc       *Doc
	workers   int
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

	//reset filterable fields
	task, err := f.index.ResetFilterableAttributes()
	if err != nil {
		return err
	}
	//wait for task
	finalTask, _ := f.client.WaitForTask(task.TaskUID)
	if finalTask.Status != "succeeded" {
		return fmt.Errorf(finalTask.Error.Code)
	}

	//update filterable fields
	task, err = f.index.UpdateFilterableAttributes(&fields)
	if err != nil {
		return err
	}

	//wait for task
	finalTask, _ = f.client.WaitForTask(task.TaskUID)
	if finalTask.Status != "succeeded" {
		return fmt.Errorf(finalTask.Error.Code)
	}
	return nil
}

//update sortable fields
func (f *Index) UpdateSortableFields(fields []string) error {
	//check
	if fields == nil || len(fields) <= 0 {
		return errors.New("invalid parameter")
	}

	//update sortable fields
	task, err := f.index.UpdateSortableAttributes(&fields)
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

//rebuild index
func (f *Index) ReCreateIndex() error {
	//remove index first
	err := f.DeleteIndex(f.indexConf.IndexName)
	if err != nil {
		return err
	}

	//begin init new index
	err = f.interInit(true)
	return err
}

//delete index
func (f *Index) DeleteIndex(indexName string) error {
	//check
	if indexName == "" {
		return errors.New("invalid parameter")
	}

	//remove index first
	task, err := f.client.DeleteIndex(indexName)
	if err != nil {
		return err
	}

	//wait for task
	finalTask, subErr := f.client.WaitForTask(task.TaskUID)
	if subErr != nil {
		return subErr
	}
	if finalTask.Status != "succeeded" {
		err = fmt.Errorf("create index failed, err:%v", finalTask.Status)
		log.Printf("init index %v failed, err:%v\n", f.indexConf.IndexName, err.Error())
		return err
	}
	return nil
}

//inter init
func (f *Index) interInit(onlyReturns ...bool) error {
	var (
		index *meilisearch.Index
		err error
		onlyReturn bool
	)
	if onlyReturns != nil && len(onlyReturns) > 0 {
		onlyReturn = onlyReturns[0]
	}

	//init index config
	indexCfg := &meilisearch.IndexConfig{
		Uid: f.indexConf.IndexName,
		PrimaryKey: f.indexConf.PrimaryKey,
	}

	//create or init index
	if f.indexConf.CreateIndex {
		//create index
		task, subErr := f.client.CreateIndex(indexCfg)
		if subErr != nil {
			if onlyReturn {
				return subErr
			}else{
				panic(any(subErr))
			}
		}

		//wait for task
		finalTask, subErrTwo := f.client.WaitForTask(task.TaskUID)
		if subErrTwo != nil {
			if onlyReturn {
				return subErrTwo
			}else{
				panic(any(subErrTwo))
			}
		}
		if finalTask.Status != "succeeded" && finalTask.Error.Code != "index_already_exists" {
			err = fmt.Errorf("create index failed, err:%v", finalTask.Status)
			log.Printf("init index %v failed, err:%v\n", f.indexConf.IndexName, err.Error())
			if onlyReturn {
				return err
			}else{
				panic(any(err))
			}
		}
		index, _ = f.client.GetIndex(f.indexConf.IndexName)
	}else{
		//init index
		index = f.client.Index(f.indexConf.IndexName)
	}

	if f.indexConf.PrimaryKey != "" {
		//set index primary key
		index.PrimaryKey = f.indexConf.PrimaryKey
	}

	//sync index obj
	f.index = index

	//set filterable fields
	if f.indexConf.FilterableFields != nil && len(f.indexConf.FilterableFields) > 0 {
		if f.indexConf.UpdateFields {
			err = f.UpdateFilterableAttributes(f.indexConf.FilterableFields)
			if err != nil {
				if onlyReturn {
					return err
				}else{
					panic(any(err))
				}
			}
		}
	}

	//set sortable fields
	if f.indexConf.SortableFields != nil && len(f.indexConf.SortableFields) > 0 {
		if f.indexConf.UpdateFields {
			err = f.UpdateSortableFields(f.indexConf.SortableFields)
			if err != nil {
				if onlyReturn {
					return err
				}else{
					panic(any(err))
				}
			}
		}
	}

	//init doc obj
	f.doc = NewDoc(f.client, f.index, f.workers)
	return nil
}