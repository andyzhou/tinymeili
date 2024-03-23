package testing

import (
	"fmt"
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
func getDoc() error {
	//get index obj
	indexObj, err := getIndexObj(IndexName)
	if err != nil || indexObj == nil {
		return err
	}

	//get doc
	docId := "1711194742684362000"
	out := NewTestDoc()
	err = indexObj.GetDoc().GetOneDocById(docId, out)
	return err
}

//benchmark get doc
func BenchmarkGetDoc(b *testing.B) {
	var (
		err error
	)
	succeed := 0
	failed := 0
	for i := 0; i < b.N; i++ {
		err = getDoc()
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
