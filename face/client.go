package face

import (
	"errors"
	"github.com/andyzhou/tinymeili/conf"
	"github.com/andyzhou/tinymeili/define"
	"github.com/meilisearch/meilisearch-go"
	"sync"
	"time"
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
func (f *Client) GetIndex(indexName string) (*Index, error) {
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
func (f *Client) CreateIndex(indexConf *conf.IndexConf) error {
	//check
	if indexConf == nil || indexConf.IndexName == "" {
		return errors.New("invalid parameter")
	}

	//init new index obj
	indexObj := NewIndex(f.client, indexConf, f.cfg.Workers)

	//sync into map
	f.Lock()
	defer f.Unlock()
	f.indexMap[indexConf.IndexName] = indexObj
	return nil
}

////////////////
//private func
////////////////

//inter init
func (f *Client) interInit() {
	//check config
	if f.cfg.TimeOut <= 0 {
		f.cfg.TimeOut = time.Duration(define.DefaultTimeOut) * time.Second
	}

	//setup client config
	clientCfg := meilisearch.ClientConfig{
		Host: f.cfg.Host,
		APIKey: f.cfg.ApiKey,
		Timeout: f.cfg.TimeOut,
	}

	//init search client
	f.client = meilisearch.NewClient(clientCfg)

	//init indexes
	if f.cfg.IndexesConf != nil {
		f.Lock()
		defer f.Unlock()
		for _, indexConf := range f.cfg.IndexesConf {
			if indexConf == nil {
				continue
			}
			//init index obj
			indexObj := NewIndex(f.client, indexConf, f.cfg.Workers)
			f.indexMap[indexConf.IndexName] = indexObj
		}
	}
}