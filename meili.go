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

//remove client
func (f *MeiLi) RemoveClient(tag string) error {
	return f.interFace.RemoveClient(tag)
}

//get client
func (f *MeiLi) GetClient(tag string) (*face.Client, error) {
	return f.interFace.GetClient(tag)
}

//add client
func (f *MeiLi) AddClient(cfg *conf.ClientConf) error {
	return f.interFace.AddClient(cfg)
}

//gen client config
func (f *MeiLi) GenClientConfig() *conf.ClientConf {
	return f.interFace.GenClientConf()
}
