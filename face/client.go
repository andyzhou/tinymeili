package face

import (
	"errors"
	"github.com/meilisearch/meilisearch-go"
	"sync"
	"time"
)

/*
 * one node client face
 * - one host, one client
 * - multi indexes for one host
 */

//client conf
type clientConf struct {
	host    string
	apiKey  string
	timeout time.Duration
	indexes []string
}

//face info
type Client struct {
	cfg *clientConf //reference
	client *meilisearch.Client
	indexMap map[string]*Index //tag -> *Index
	sync.RWMutex
}

//construct
func NewClient(cfg *clientConf) *Client {
	this := &Client{
		cfg: cfg,
		indexMap: map[string]*Index{},
	}
	this.interInit()
	return this
}

//get index by name
func (f *Client) GetIndex(
	indexName string) (*Index, error) {
	//check
	if indexName == "" {
		return nil, errors.New("invalid parameter")
	}
	//get with locker
	f.Lock()
	defer f.Unlock()
	v, ok := f.indexMap[indexName]
	if ok && v != nil {
		return v, nil
	}
	return nil, errors.New("no such index by tag")
}

//create and init index
func (f *Client) CreateIndex(
	indexName string,
	filterableFields ...string) error {
	//check
	if indexName == "" {
		return errors.New("invalid parameter")
	}

	//init index
	index := f.client.Index(indexName)

	//update filterable fields
	if filterableFields != nil && len(filterableFields) > 0 {
		_, err := index.UpdateFilterableAttributes(&filterableFields)
		if err != nil {
			return err
		}
	}

	//init new index
	indexObj := NewIndex(indexName, index)

	//sync into map
	f.Lock()
	defer f.Unlock()
	f.indexMap[indexName] = indexObj
	return nil
}

////////////////
//private func
////////////////

//inter init
func (f *Client) interInit() {
	//setup client config
	clientCfg := meilisearch.ClientConfig{
		Host: f.cfg.host,
		APIKey: f.cfg.apiKey,
		Timeout: f.cfg.timeout,
	}

	//init search client
	f.client = meilisearch.NewClient(clientCfg)

	//init indexes
	if f.cfg.indexes != nil {
		f.Lock()
		defer f.Unlock()
		for _, indexName := range f.cfg.indexes {
			if indexName == "" {
				continue
			}
			//init index
			index := f.client.Index(indexName)
			indexObj := NewIndex(indexName, index)
			f.indexMap[indexName] = indexObj
		}
	}
}