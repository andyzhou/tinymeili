package testing

import (
	"github.com/andyzhou/tinymeili"
	"github.com/andyzhou/tinymeili/conf"
	"github.com/andyzhou/tinymeili/face"
)

const (
	HostTag = "test"
	Host = "http://127.0.0.1:7700"
	ApiKey = "test"
	IndexName = "test2"
	PrimaryKey = "id"
)

var (
	mc *tinymeili.MeiLi
	clientCfg *conf.ClientConf
)

//init
func init()  {
	//init client
	mc = tinymeili.GetMeiLi()

	//gen and fill client config
	clientCfg = mc.GenClientConfig()
	clientCfg.Tag = HostTag
	clientCfg.Host = Host
	clientCfg.ApiKey = ApiKey
	clientCfg.IndexesConf = []*conf.IndexConf{
		{
			IndexName: IndexName,
			PrimaryKey: PrimaryKey,
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

