package main

import (
	"errors"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/andyzhou/tinymeili"
	"github.com/andyzhou/tinymeili/conf"
	"github.com/andyzhou/tinymeili/define"
	"github.com/andyzhou/tinymeili/face"
)

/*
 * example demo
 */

const (
	HostTag   = "test"
	Host      = "http://127.0.0.1:7700"
	ApiKey    = "test"
	IndexName = "test2"
)

var (
	mc *tinymeili.MeiLi
)

//init
func init()  {
	//init client
	mc = tinymeili.GetMeiLi()

	//gen and fill client config
	clientCfg := mc.GenClientConfig()
	clientCfg.Tag = HostTag
	clientCfg.Host = Host
	clientCfg.ApiKey = ApiKey
	clientCfg.TimeOut = time.Duration(30) * time.Second
	clientCfg.Workers = define.DefaultWorkers
	clientCfg.IndexesConf = []*conf.IndexConf{
		&conf.IndexConf{
			IndexName: IndexName,
			PrimaryKey: "id",
			FilterableFields: []string{"poster", "tags", "property"},
			CreateIndex: true,
			UpdateFields: true,
		},
	}

	//add client
	err := mc.AddClient(clientCfg)
	if err != nil {
		panic(any(err))
	}
}

//get index obj
func getIndexObj(indexName string) (*face.Index, error) {
	client, subErr := mc.GetClient(HostTag)
	if subErr != nil || client == nil {
		return nil, subErr
	}
	index, subErrTwo := client.GetIndex(indexName)
	return index, subErrTwo
}

//add new doc
func addDoc(beginId int64) error {
	//setup tags
	tags := []string{
		"go",
		"java",
		"编程",
		"技术",
		"mysql",
		"redis",
	}
	randIdx := rand.Intn(len(tags)) + 1
	randTags := tags[0:randIdx]
	//init obj
	obj := NewTestDoc()
	obj.Id = beginId
	obj.Poster = int64(rand.Intn(5) + 1)
	obj.Tags = randTags
	obj.Property = map[string]interface{}{
		"sex":1,
		"age":10,
		"city":"beijing",
	}

	////get embed value
	//vector, _ := GetEmbedding("go")
	//if len(vector) > 0 {
	//	obj.Vectors = vector
	//}

	//get index obj
	indexObj, err := getIndexObj(IndexName)
	if err != nil || indexObj == nil {
		return err
	}

	//add doc
	err = indexObj.GetDoc().AddDoc(*obj)
	return err
}

//del doc
func delDoc() error {
	//get index obj
	indexObj, err := getIndexObj(IndexName)
	if err != nil || indexObj == nil {
		return err
	}

	attributeFields := []string{"tags"}
	err = indexObj.UpdateFilterableAttributes(attributeFields)

	//del doc by ids
	docIds := []string{"1711266923113935000"}
	err = indexObj.GetDoc().DelDoc("", docIds...)

	////del doc by filter
	//filter := "tags = 'china'"
	//err = indexObj.GetDoc().DelDocsByFilter([]string{filter})
	return err
}

//get multi doc
func getMultiDoc() error {
	//get index obj
	indexObj, err := getIndexObj(IndexName)
	if err != nil || indexObj == nil {
		return err
	}

	//update filterable fields
	attributeFields := []string{"id"}
	err = indexObj.UpdateFilterableAttributes(attributeFields)

	//del doc by ids
	docIds := []string{"28", "29"}
	resp, subErr := indexObj.GetDoc().GetBatchDocsByIds("id", docIds...)
	log.Printf("resp:%v\n", resp)
	return subErr
}

//query doc
func queryDoc() ([]interface{}, interface{}, error) {
	//get index obj
	indexObj, err := getIndexObj(IndexName)
	if err != nil || indexObj == nil {
		return nil, nil, err
	}

	//filter
	//filter := "property.age >= 0 AND property.age < 10"
	//facets := []string{"tags"}

	//filter := "_semanticSimilarity('go', tags) > 0.1"
	//sorter := []string{"_semanticSimilarity('go', tags):desc"}

	//distinct field
	//distinctField := "poster"

	//setup query para
	para := &define.QueryPara{
		//Filter: filter,
		//Distinct: distinctField,
		//Sort: sorter,
		Page: 1,
		PageSize: 10,
	}
	_, resp, facetObj, subErr := indexObj.GetDoc().QueryIndexDocs(para)
	return resp, facetObj, subErr
}

//add batch doc
func addBatchDoc()  {
	now := time.Now().UnixNano()
	for i := int64(0); i < 5; i++ {
		//add new doc
		beginId := now + i
		err := addDoc(beginId)
		if err != nil {
			log.Printf("add doc failed, err:%v\n", err.Error())
			return
		}
	}
}

//recreate index
func recreateIndex() error {
	client, err := mc.GetClient(HostTag)
	if err != nil {
		return err
	}
	err = client.ReCreateIndex(IndexName)
	return err
}

//create index
func createIndex(cfg *conf.IndexConf) error {
	return nil
	if cfg == nil {
		return errors.New("invalid parameter")
	}
	client, err := mc.GetClient(HostTag)
	if err != nil {
		return err
	}
	err = client.CreateIndex(cfg)
	return err
}

func main() {
	var (
		wg sync.WaitGroup
		err error
	)
	sf := func() {
		//mc.Quit()
		log.Printf("tiny meiLi quit..\n")
		time.Sleep(time.Second)
		//wg.Done()
	}
	time.AfterFunc(time.Second * 2, sf)
	wg.Add(1)

	//add batch doc
	//addBatchDoc()

	////del doc
	//err = delDoc()
	//if err != nil {
	//	log.Printf("del doc failed, err:%v\n", err.Error())
	//	return
	//}

	//get multi docs
	//getMultiDoc()

	//query doc
	resp, facets, err := queryDoc()
	log.Printf("query doc, resp:%v, facets:%v, err:%v\n", resp, facets, err)

	//create index
	indexCfg := &conf.IndexConf{
		IndexName: IndexName,
		PrimaryKey: "id",
		FilterableFields: []string{
			"poster",
			"property",
			"tags",
			"vectors",
		},
		CreateIndex: true,
		UpdateFields: true,
	}
	err = createIndex(indexCfg)
	log.Printf("recreate index, resp:%v\n", err)

	wg.Wait()
	log.Printf("doc opt succeed\n")
}