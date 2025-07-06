package face

import (
	"errors"
	"runtime"
	"sync"
	"time"

	"github.com/andyzhou/tinymeili/conf"
	"github.com/andyzhou/tinymeili/define"
	"github.com/meilisearch/meilisearch-go"
)

/*
 * one node client face
 * @author <AndyZhou>
 * @mail <diudiu8848@163.com>
 * - one host, one client
 * - multi indexes for one host
 */

//face info
type Client struct {
	cfg      *conf.ClientConf //reference
	client   meilisearch.ServiceManager
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
	//release index map
	f.Lock()
	defer f.Unlock()
	for k, v := range f.indexMap {
		v.Quit()
		delete(f.indexMap, k)
	}

	//gc opt
	f.indexMap = map[string]*Index{}
	runtime.GC()
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
	if indexConf == nil ||
		indexConf.IndexName == "" ||
		indexConf.PrimaryKey == "" {
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

//re-create and init index
func (f *Client) ReCreateIndex(indexName string) error {
	//check
	if indexName == "" {
		return errors.New("invalid parameter")
	}

	//get index by name with locker
	f.Lock()
	defer f.Unlock()
	index, ok := f.indexMap[indexName]
	if !ok || index == nil {
		return errors.New("no such index")
	}

	//begin recreate index
	err := index.ReCreateIndex()
	return err
}

////////////////
//private func
////////////////

//inter init
func (f *Client) interInit() {
	var (
		err error
	)
	//check config
	if f.cfg.TimeOut <= 0 {
		f.cfg.TimeOut = time.Duration(define.DefaultTimeOut) * time.Second
	}

	////setup client config, for old api version
	//clientCfg := meilisearch.ClientConfig{
	//	Host: f.cfg.Host,
	//	APIKey: f.cfg.ApiKey,
	//	Timeout: f.cfg.TimeOut,
	//}
	//
	//client := meilisearch.New(f.cfg.Host, meilisearch.WithAPIKey(f.cfg.ApiKey))

	//init search client
	f.client = meilisearch.New(f.cfg.Host, meilisearch.WithAPIKey(f.cfg.ApiKey))

	//init indexes
	if f.cfg.IndexesConf != nil {
		for _, indexConf := range f.cfg.IndexesConf {
			if indexConf == nil {
				continue
			}
			//init index obj
			err = f.CreateIndex(indexConf)
			if err != nil {
				panic(any(err))
			}
		}
	}
}