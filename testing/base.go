package testing

import (
	"github.com/andyzhou/tinymeili"
	"github.com/andyzhou/tinymeili/conf"
	"github.com/andyzhou/tinymeili/face"
)

const (
	NodeKind = "test"
	HostTag = "test"
	Host = "http://127.0.0.1:7700"
	ApiKey = ""
	IndexName = "test3"
)

var (
	mc *tinymeili.MeiLi
	nodeCfg *conf.NodeConf
)

//init
func init()  {
	//init client
	mc = tinymeili.GetMeiLi()

	//gen and fill node config
	nodeCfg = mc.GenNodeConfig()
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

