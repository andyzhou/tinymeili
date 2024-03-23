package face

import (
	"errors"
	"github.com/andyzhou/tinymeili/conf"
	"sync"
	"time"
)

/*
 * node face
 */

//face info
type Node struct {
	cfg *conf.NodeConf //reference
	clientMap map[string]*Client //tag -> *Client
	sync.RWMutex
}

//construct
func NewNode(cfg *conf.NodeConf) *Node {
	this := &Node{
		cfg: cfg,
		clientMap: map[string]*Client{},
	}
	this.interInit()
	return this
}

//get client by host tag
func (f *Node) GetClient(hostTag string) (*Client, error) {
	//check
	if hostTag == "" {
		return nil, errors.New("invalid parameter")
	}
	//get by tag with locker
	f.Lock()
	defer f.Unlock()
	v, ok := f.clientMap[hostTag]
	if ok && v != nil {
		return v, nil
	}
	return nil, errors.New("no such client by tag")
}

//inter init
func (f *Node) interInit() {
	//init batch clients
	f.Lock()
	defer f.Unlock()
	for tag, host := range f.cfg.Hosts {
		if tag == "" || host == "" {
			continue
		}
		//init client conf
		cfg := &clientConf{
			host: host,
			apiKey: f.cfg.ApiKey,
			indexes: f.cfg.Indexes,
		}
		if f.cfg.TimeOut > 0 {
			cfg.timeout = f.cfg.TimeOut * time.Second
		}

		//init new client
		client := NewClient(cfg)
		f.clientMap[tag] = client
	}
}