package main

import (
	"fmt"
	"github.com/andyzhou/tinymeili"
	"github.com/andyzhou/tinymeili/face"
	"log"
	"math/rand"
	"time"
)

/*
 * example demo
 */

const (
	NodeKind = "test"
	HostTag = "test"
	Host = "http://127.0.0.1:7700"
	ApiKey = ""
	IndexName = "test3"
)

var (
	mc *tinymeili.MeiLi
)

//init
func init()  {
	//init client
	mc = tinymeili.GetMeiLi()

	//gen and fill node config
	nodeCfg := mc.GenNodeConfig()
	nodeCfg.Kind = NodeKind
	nodeCfg.Hosts = map[string]string{
		HostTag : Host,
	}
	nodeCfg.ApiKey = ApiKey
	nodeCfg.Indexes = []string{IndexName}

	//add node
	err := mc.AddNode(nodeCfg)
	if err != nil {
		panic(any(err))
	}
}

//get index obj
func getIndexObj(indexName string) (*face.Index, error) {
	node, err := mc.GetNode(NodeKind)
	if err != nil || node == nil {
		return nil, err
	}
	client, subErr := node.GetClient(HostTag)
	if subErr != nil || client == nil {
		return nil, subErr
	}
	index, subErrTwo := client.GetIndex(indexName)
	return index, subErrTwo
}

//add new doc
func addDoc() error {
	//init obj
	obj := NewTestDoc()
	obj.Id = time.Now().UnixNano()
	obj.Title = fmt.Sprintf("test-%v", rand.Int63n(time.Now().UnixNano()))
	obj.Tags = []string{
		"beijing",
		"china",
	}

	//get index obj
	indexObj, err := getIndexObj(IndexName)
	if err != nil || indexObj == nil {
		return err
	}

	//add doc
	err = indexObj.GetDoc().AddDoc(obj)
	return err
}

func main() {
	//add new doc
	err := addDoc()
	if err != nil {
		log.Printf("add doc failed, err:%v\n", err.Error())
		return
	}
	log.Printf("add doc succeed\n")
}