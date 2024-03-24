package face

import (
	"errors"
	"github.com/andyzhou/tinymeili/conf"
	"sync"
)

/*
 * inter face
 * - base on meili search `https://github.com/meilisearch/meilisearch-go`
 */

//global variable
var (
	_face *InterFace
	_faceOnce sync.Once
)

//face info
type InterFace struct {
	clientMap map[string]*Client //kind -> *Client
	sync.RWMutex
}

//get single instance
func GetInterFace() *InterFace {
	_faceOnce.Do(func() {
		_face = NewInterFace()
	})
	return _face
}

//construct
func NewInterFace() *InterFace {
	this := &InterFace{
		clientMap: map[string]*Client{},
	}
	return this
}

//remove client by tag
func (f *InterFace) RemoveClient(tag string) error {
	//check
	if tag == "" {
		return errors.New("invalid parameter")
	}

	//remove with locker
	f.Lock()
	defer f.Unlock()
	delete(f.clientMap, tag)
	return nil
}

//get all clients
func (f *InterFace) GetAllClient() map[string]*Client {
	f.Lock()
	defer f.Unlock()
	return f.clientMap
}

//get client by tag
func (f *InterFace) GetClient(tag string) (*Client, error) {
	//check
	if tag == "" {
		return nil, errors.New("invalid parameter")
	}

	//get by k with locker
	f.Lock()
	defer f.Unlock()
	v, ok := f.clientMap[tag]
	if ok && v != nil {
		return v, nil
	}
	return nil, errors.New("no such node by kind")
}

//add new client
func (f *InterFace) AddClient(cfg *conf.ClientConf) error {
	//check
	if cfg == nil || cfg.Tag == "" {
		return errors.New("invalid parameter")
	}

	//check and init client
	f.Lock()
	defer f.Unlock()
	_, ok := f.clientMap[cfg.Tag]
	if ok {
		return errors.New("node has exists")
	}

	//init new client
	client := NewClient(cfg)
	f.clientMap[cfg.Tag] = client
	return nil
}

//gen client conf
func (f *InterFace) GenClientConf() *conf.ClientConf {
	return &conf.ClientConf{
		Indexes: []string{},
	}
}