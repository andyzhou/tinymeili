package face

import (
	"errors"
	"github.com/andyzhou/tinymeili/conf"
	"github.com/meilisearch/meilisearch-go"
	"sync"
)

/*
 * one node client face
 * - one host, one client
 * - multi indexes for one host
 */

//face info
type Client struct {
	cfg *conf.ClientConf //reference
	client *meilisearch.Client
	indexMap map[string]*Index //tag -> *Index
	sync.RWMutex
}

//construct
func NewClient(cfg *conf.ClientConf) *Client {
	this := &Client{
		cfg: cfg,
		indexMap: map[string]*Index{},
	}
	this.interInit()
	return this
}

//quit
func (f *Client) Quit() {
	for _, v := range f.indexMap {
		v.Quit()
	}
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
	indexObj := NewIndex(indexName, index, f.cfg.Workers)

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
		Host: f.cfg.Host,
		APIKey: f.cfg.ApiKey,
		Timeout: f.cfg.TimeOut,
	}

	//init search client
	f.client = meilisearch.NewClient(clientCfg)

	//init indexes
	if f.cfg.Indexes != nil {
		f.Lock()
		defer f.Unlock()
		for _, indexName := range f.cfg.Indexes {
			if indexName == "" {
				continue
			}
			//init index
			index := f.client.Index(indexName)
			indexObj := NewIndex(indexName, index, f.cfg.Workers)
			f.indexMap[indexName] = indexObj
		}
	}
}