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
	nodeMap map[string]*Node //kind -> *Node
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
		nodeMap: map[string]*Node{},
	}
	return this
}

//remove node by kind
func (f *InterFace) RemoveNode(nodeKind string) error {
	//check
	if nodeKind == "" {
		return errors.New("invalid parameter")
	}

	//remove with locker
	f.Lock()
	defer f.Unlock()
	delete(f.nodeMap, nodeKind)
	return nil
}

//get node by kind
func (f *InterFace) GetNode(nodeKind string) (*Node, error) {
	//check
	if nodeKind == "" {
		return nil, errors.New("invalid parameter")
	}

	//get by tag with locker
	f.Lock()
	defer f.Unlock()
	v, ok := f.nodeMap[nodeKind]
	if ok && v != nil {
		return v, nil
	}
	return nil, errors.New("no such node by kind")
}

//add new node
func (f *InterFace) AddNode(cfg *conf.NodeConf) error {
	//check
	if cfg == nil || cfg.Kind == "" {
		return errors.New("invalid parameter")
	}

	//check and init node
	f.Lock()
	defer f.Unlock()
	_, ok := f.nodeMap[cfg.Kind]
	if ok {
		return errors.New("node has exists")
	}

	//init new node
	node := NewNode(cfg)
	f.nodeMap[cfg.Kind] = node
	return nil
}

//gen node conf
func (f *InterFace) GenNodeConf() *conf.NodeConf {
	return &conf.NodeConf{
		Hosts: map[string]string{},
		Indexes: []string{},
	}
}