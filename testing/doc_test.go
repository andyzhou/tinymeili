package testing

import (
	"fmt"
	"github.com/andyzhou/tinymeili/define"
	"math/rand"
	"testing"
	"time"
)

//add new doc
func addDoc() error {
	//init obj
	obj := NewTestDoc()
	obj.Id = time.Now().UnixNano()
	obj.Title = fmt.Sprintf("test-%v", rand.Int63n(time.Now().UnixNano()))
	obj.Tags = []string{
		"beijing",
		"china",
	}

	//get index obj
	indexObj, err := getIndexObj(IndexName)
	if err != nil || indexObj == nil {
		return err
	}

	//add doc
	err = indexObj.GetDoc().AddDoc(obj)
	return err
}

//get one doc
func getDoc() (*TestDoc, error) {
	//get index obj
	indexObj, err := getIndexObj(IndexName)
	if err != nil || indexObj == nil {
		return nil, err
	}

	//get doc
	docId := "1711277520797064000"
	out := NewTestDoc()
	err = indexObj.GetDoc().GetOneDocById(docId, out)
	return out, err
}

//query doc
func queryDoc() error {
	//get index obj
	indexObj, err := getIndexObj(IndexName)
	if err != nil || indexObj == nil {
		return err
	}

	para := &define.QueryPara{
		Key: "test",
		Page: 1,
		PageSize: 10,
	}
	_, _, _, subErr := indexObj.GetDoc().QueryIndexDocs(para)
	return subErr
}

//test add doc
func TestAddDoc(t *testing.T) {
	err := addDoc()
	if err != nil {
		t.Errorf("add doc failed, err:%v\n", err.Error())
	}else{
		t.Logf("add doc succeed")
	}
}

func TestGetDoc(t *testing.T) {
	doc, err := getDoc()
	t.Logf("doc:%v, err:%v\n", doc, err)
}

//benchmark get doc
func BenchmarkGetDoc(b *testing.B) {
	var (
		err error
	)
	succeed := 0
	failed := 0
	for i := 0; i < b.N; i++ {
		_, err = getDoc()
		if err != nil {
			failed++
		}else{
			succeed++
		}
	}
	b.Logf("benchmark get doc, succeed:%v, failed:%v\n", succeed, failed)
}

//benchmark add doc
func BenchmarkAddDoc(b *testing.B) {
	var (
		err error
	)
	succeed := 0
	failed := 0
	for i := 0; i < b.N; i++ {
		err = addDoc()
		if err != nil {
			failed++
		}else{
			succeed++
		}
	}
	b.Logf("benchmark add doc, succeed:%v, failed:%v\n", succeed, failed)
}

//benchmark query doc
func BenchmarkQueryDoc(b *testing.B) {
	var (
		err error
	)
	succeed := 0
	failed := 0
	for i := 0; i < b.N; i++ {
		err = queryDoc()
		if err != nil {
			failed++
		}else{
			succeed++
		}
	}
	b.Logf("benchmark query doc, succeed:%v, failed:%v\n", succeed, failed)
	time.Sleep(time.Second)
}