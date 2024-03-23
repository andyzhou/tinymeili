package tinymeili

import (
	"github.com/andyzhou/tinymeili/conf"
	"github.com/andyzhou/tinymeili/face"
	"sync"
)

/*
 * lib face
 */

//global variable
var (
	_meiLi *MeiLi
	_meiLiOnce sync.Once
)

//face info
type MeiLi struct {
	interFace *face.InterFace
}

//get single instance
func GetMeiLi() *MeiLi {
	_meiLiOnce.Do(func() {
		_meiLi = NewMeiLi()
	})
	return _meiLi
}

//construct
func NewMeiLi() *MeiLi {
	this := &MeiLi{
		interFace: face.NewInterFace(),
	}
	return this
}

///////////
//api
///////////

//remove node
func (f *MeiLi) RemoveNode(nodeKind string) error {
	return f.interFace.RemoveNode(nodeKind)
}

//get node
func (f *MeiLi) GetNode(nodeKind string) (*face.Node, error) {
	return f.interFace.GetNode(nodeKind)
}

//add node
func (f *MeiLi) AddNode(cfg *conf.NodeConf) error {
	return f.interFace.AddNode(cfg)
}

//gen node config
func (f *MeiLi) GenNodeConfig() *conf.NodeConf {
	return f.interFace.GenNodeConf()
}
