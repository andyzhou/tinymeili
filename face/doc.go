package face

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/andyzhou/tinymeili/define"
	"github.com/andyzhou/tinymeili/lib"
	"github.com/meilisearch/meilisearch-go"
	"log"
	"strconv"
)

/*
 * doc opt face
 */

//inter opt
type (
	syncDocReq struct {
		obj interface{}
		isUpdate bool
	}
	removeDocReq struct {
		docIds []string
		filter []string
	}
)

//face info
type Doc struct {
	client *meilisearch.Client //reference
	index *meilisearch.Index //reference
	worker *lib.Worker
	workers int
}

//construct
func NewDoc(
	client *meilisearch.Client,
	index *meilisearch.Index,
	workers int) *Doc {
	this := &Doc{
		client: client,
		workers: workers,
		index: index,
		worker: lib.NewWorker(),
	}
	this.interInit()
	return this
}

//quit
func (f *Doc) Quit() {
	if f.worker != nil {
		f.worker.Quit()
	}
}

//query batch doc one index
//sync opt
//return total, []docObj, facetMap, error
func (f *Doc) QueryIndexDocs(
		para *define.QueryPara,
	) (int64, []interface{}, map[string]map[string]int64, error) {
	//check
	if para == nil {
		return 0, nil, nil, errors.New("invalid parameter")
	}
	if f.index == nil {
		return 0, nil, nil, errors.New("inter index not init")
	}

	//setup offset
	if para.Page <= 0 {
		para.Page = define.DefaultPage
	}
	if para.PageSize <= 0 {
		para.PageSize = define.DefaultPageSize
	}
	offset := (para.Page - 1) * para.PageSize
	limit := para.PageSize

	//setup search request
	sq := &meilisearch.SearchRequest{
		Query: para.Key,
		Filter: para.Filter,
		Facets: para.Facets,
		Sort: para.Sort,
		Offset: int64(offset),
		Limit: int64(limit),
		HitsPerPage:int64(para.PageSize),
	}

	//query origin doc
	resp, subErr := f.index.Search(para.Key, sq)
	if subErr != nil || resp == nil {
		return 0, nil, nil, subErr
	}

	//gather facet objs
	facetObjs := make(map[string]map[string]int64)
	if resp.FacetDistribution != nil {
		facetMap, ok := resp.FacetDistribution.(map[string]interface{})
		if ok && facetMap != nil {
			for k, v := range facetMap {
				if k == "" || v == nil {
					continue
				}
				//sub facet objs
				facetObj, subOk := v.(map[string]interface{})
				if !subOk || facetObj == nil {
					continue
				}
				subFacetObj := make(map[string]int64)
				for k1, v1 := range facetObj {
					countVal, _ := strconv.ParseInt(fmt.Sprintf("%v", v1), 10, 64)
					subFacetObj[k1] = countVal
				}
				//gather one key and sub facet objs
				facetObjs[k] = subFacetObj
			}
		}
	}
	return resp.TotalHits, resp.Hits, facetObjs, nil
}

//get bach doc by ids
//field need set as filterable
//return []interface{}, error
func (f *Doc) GetBatchDocByIds(
	field string,
	docIds ...string) ([]interface{}, error) {
	//check
	if field == "" || docIds == nil {
		return nil, errors.New("invalid parameter")
	}
	if f.index == nil {
		return nil, errors.New("inter index not init")
	}

	filterBuff := bytes.NewBuffer(nil)
	i := 0
	limit := int64(0)
	for _, docId := range docIds {
		if docId == "" {
			continue
		}
		if i > 0 {
			filterBuff.WriteString(" OR ")
		}
		filterBuff.WriteString(fmt.Sprintf("%v = %v", field, docId))
		i++
		limit++
	}

	//setup search request
	sq := &meilisearch.SearchRequest{
		AttributesToSearchOn:[]string{field},
		Filter: []string{filterBuff.String()},
		Limit: limit,
		HitsPerPage:limit,
	}

	//query origin doc
	resp, err := f.index.Search("", sq)
	if err != nil || resp == nil {
		return nil, err
	}
	return resp.Hits, nil
}

//get one doc by field condition
//sync opt
func (f *Doc) GetOneDocByFieldCond(
	matchField, matchVal string,
	out interface{}) error {
	//check
	if matchVal == "" ||
		matchField == "" || out == nil {
		return errors.New("invalid parameter")
	}
	if f.index == nil {
		return errors.New("inter index not init")
	}

	//setup search request
	sq := &meilisearch.SearchRequest{
		Offset: 0,
		Limit: 1,
		HitsPerPage:1,
		AttributesToSearchOn:[]string{matchField},
	}

	//get origin doc
	resp, subErr := f.index.Search(matchVal, sq)
	if subErr != nil || resp == nil || resp.Hits == nil {
		return subErr
	}

	//get first hit doc
	hitDoc := resp.Hits[0]
	recMap, ok := hitDoc.(map[string]interface{})
	if !ok || recMap == nil {
		return errors.New("invalid hit doc format")
	}

	//decode to out obj
	recBytes, _ := json.Marshal(recMap)
	err := json.Unmarshal(recBytes, out)
	return err
}

//get one doc by id
func (f *Doc) GetOneDocById(
	docId string,
	out interface{}) error {
	//check
	if docId == "" || out == nil {
		return errors.New("invalid parameter")
	}
	if f.index == nil {
		return errors.New("inter index not init")
	}

	//get real doc
	err := f.index.GetDocument(
		docId,
		nil,
		&out)
	return err
}

//del one doc
func (f *Doc) DelDoc(
	docIds ...string) error {
	//check
	if docIds == nil || len(docIds) <= 0 {
		return errors.New("invalid parameter")
	}
	if f.index == nil {
		return errors.New("inter index not init")
	}

	//init request
	req := removeDocReq{
		docIds: docIds,
	}

	//send worker queue
	_, err := f.worker.SendData(req, "")
	return err
}

//del docs by filter
//filter like: 'a = 6 and b < 10'
func (f *Doc) DelDocsByFilter(
	filter []string) error {
	//check
	if filter == nil {
		return errors.New("invalid parameter")
	}
	if f.index == nil {
		return errors.New("inter index not init")
	}

	//init request
	req := removeDocReq{
		filter: filter,
	}

	//send worker queue
	_, err := f.worker.SendData(req, "")
	return err
}

//update one doc
func (f *Doc) UpdateDoc(
	docObj interface{},
	docIds ...string) error {
	var (
		docId string
	)
	//check
	if docObj == nil {
		return errors.New("invalid parameter")
	}
	if f.index == nil {
		return errors.New("inter index not init")
	}
	if docIds != nil && len(docIds) > 0 {
		docId = docIds[0]
	}

	//init request
	req := syncDocReq{
		obj: docObj,
		isUpdate: true,
	}

	//send worker queue
	_, err := f.worker.SendData(req, docId)
	return err
}

//add one or batch doc
func (f *Doc) AddDoc(
	docObj interface{},
	docIds ...string) error {
	var (
		docId string
	)
	//check
	if docObj == nil {
		return errors.New("invalid parameter")
	}
	if f.index == nil {
		return errors.New("inter index not init")
	}
	if docIds != nil && len(docIds) > 0 {
		docId = docIds[0]
	}
	//init request
	req := syncDocReq{
		obj: docObj,
	}

	//send worker queue
	_, err := f.worker.SendData(req, docId)
	return err
}

/////////////////
//private func
/////////////////

//remove doc
func (f *Doc) removeDocObj(req *removeDocReq) error {
	var (
		resp *meilisearch.TaskInfo
		err error
	)
	//check
	if req == nil {
		return errors.New("invalid parameter")
	}
	if f.index == nil {
		return errors.New("inter index not init")
	}

	//remove real doc
	if req.filter != nil {
		//remove by filter
		resp, err = f.index.DeleteDocumentsByFilter(req.filter)
	}else{
		//remove by ids
		resp, err = f.index.DeleteDocuments(req.docIds)
	}
	if err != nil {
		return err
	}

	//wait for task status
	finalTask, _ := f.client.WaitForTask(resp.TaskUID)
	if finalTask.Status != "succeeded" {
		log.Printf("doc.removeDocObj failed, err:%v\n", finalTask.Error.Code)
		return fmt.Errorf(finalTask.Error.Code)
	}
	return nil
}

//add or update doc
func (f *Doc) syncDocObj(req *syncDocReq) (*meilisearch.TaskInfo, error) {
	var (
		resp *meilisearch.TaskInfo
		err error
	)
	//check
	if req == nil || req.obj == nil {
		return nil, errors.New("invalid parameter")
	}
	if f.index == nil {
		return nil, errors.New("inter index not init")
	}

	//add real doc
	if req.isUpdate {
		resp, err = f.index.UpdateDocuments(req.obj, f.index.PrimaryKey)
	}else{
		resp, err = f.index.AddDocuments(req.obj, f.index.PrimaryKey)
	}

	//wait for task status
	finalTask, _ := f.client.WaitForTask(resp.TaskUID)
	if finalTask.Status != "succeeded" {
		log.Printf("doc.syncDocObj failed, err:%v\n", finalTask.Error.Code)
		return nil, fmt.Errorf(finalTask.Error.Code)
	}
	return resp, err
}

//cb for worker opt
func (f *Doc) cbForWorkerOpt(input interface{}) (interface{}, error) {
	var (
		err error
	)
	//check
	if input == nil {
		return nil, errors.New("invalid parameter")
	}

	//do diff opt by data type
	switch dataType := input.(type) {
	case syncDocReq:
		{
			//add or update doc opt
			req, ok := input.(syncDocReq)
			if !ok || &req == nil {
				return nil, errors.New("invalid data type")
			}
			_, err = f.syncDocObj(&req)
			return nil, err
		}
	case removeDocReq:
		{
			//remove doc opt
			req, ok := input.(removeDocReq)
			if !ok || &req == nil {
				return nil, errors.New("invalid data type")
			}
			err = f.removeDocObj(&req)
			return nil, err
		}
	default:
		{
			return nil, fmt.Errorf("invalid data type `%v`", dataType)
		}
	}
}

//inter init
func (f *Doc) interInit() {
	//check worker num
	if f.workers <= 0 {
		f.workers = lib.DefaultWorkers
	}

	//init workers
	f.worker.SetCBForQueueOpt(f.cbForWorkerOpt)
	f.worker.CreateWorkers(f.workers)
}