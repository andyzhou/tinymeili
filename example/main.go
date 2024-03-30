package main

import (
	"github.com/andyzhou/tinymeili"
	"github.com/andyzhou/tinymeili/conf"
	"github.com/andyzhou/tinymeili/define"
	"github.com/andyzhou/tinymeili/face"
	"log"
	"sync"
	"time"
)

/*
 * example demo
 */

const (
	HostTag = "test"
	Host = "http://127.0.0.1:7700"
	ApiKey = "test"
	IndexName = "test_3"
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
			FilterableFields: []string{"dataId"},
			CreateIndex: false,
			UpdateFields: false,
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
	now := time.Now().Unix()

	//init obj
	obj := NewReviewJson()
	obj.Id = beginId
	obj.DataId = 1
	obj.CreateAt = now

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
	attributeFields := []string{"id", "tags"}
	err = indexObj.UpdateFilterableAttributes(attributeFields)

	//del doc by ids
	docIds := []string{"1711273936556697000", "1711273936562740000"}
	_, err = indexObj.GetDoc().GetBatchDocByIds("id", docIds...)
	return nil
}

//query doc
func queryDoc() ([]interface{}, interface{}, error) {
	//get index obj
	indexObj, err := getIndexObj(IndexName)
	if err != nil || indexObj == nil {
		return nil, nil, err
	}

	//filter
	filter := "dataId = 1"
	//facets := []string{"tags"}

	//setup query para
	para := &define.QueryPara{
		Filter: filter,
		Page: 1,
		PageSize: 2,
	}
	_, resp, facetObj, subErr := indexObj.GetDoc().QueryIndexDocs(para)
	return resp, facetObj, subErr
}

func main() {
	var (
		wg sync.WaitGroup
	)
	sf := func() {
		//mc.Quit()
		log.Printf("tiny meiLi quit..\n")
		time.Sleep(time.Second)
		//wg.Done()
	}
	time.AfterFunc(time.Second * 2, sf)
	wg.Add(1)

	//now := time.Now().UnixNano()
	//for i := int64(0); i < 1; i++ {
	//	//add new doc
	//	beginId := now + i
	//	err := addDoc(beginId)
	//	if err != nil {
	//		log.Printf("add doc failed, err:%v\n", err.Error())
	//		return
	//	}
	//}

	////del doc
	//err = delDoc()
	//if err != nil {
	//	log.Printf("del doc failed, err:%v\n", err.Error())
	//	return
	//}

	//get multi docs
	//getMultiDoc()

	////query doc
	//resp, facets, err := queryDoc()
	//log.Printf("query doc, resp:%v, facets:%v, err:%v\n", resp, facets, err)

	wg.Wait()
	log.Printf("doc opt succeed\n")
}