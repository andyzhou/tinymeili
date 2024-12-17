package face

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/andyzhou/tinymeili/define"
	"github.com/andyzhou/tinymeili/lib"
	"github.com/meilisearch/meilisearch-go"
)

/*
 * doc opt face
 * @author <AndyZhou>
 * @mail <diudiu8848@163.com>
 */

//inter opt
type (
	syncDocReq struct {
		obj      interface{}
		isUpdate bool
	}
	removeDocReq struct {
		docIds []string
		filter []string
	}
)

//face info
type Doc struct {
	client  *meilisearch.Client //reference
	index   *meilisearch.Index  //reference
	worker  *lib.Worker
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

	//setup search request
	sq := &meilisearch.SearchRequest{
		Query: para.Key,
		AttributesToSearchOn: para.AttributesToSearch,
		Filter: para.Filter,
		Facets: para.Facets,
		Sort: para.Sort,
		Page: int64(para.Page),
		HitsPerPage:int64(para.PageSize),
	}
	if para.Distinct != "" {
		sq.Distinct = para.Distinct
	}
	if para.AttributesToSearch != nil && len(para.AttributesToSearch) > 0 {
		sq.AttributesToSearchOn = para.AttributesToSearch
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

//get batch doc by ids
//field need set as filterable
func (f *Doc) GetBatchDocsByIds(
		condField string,
		docIds ...string,
	) ([]map[string]interface{}, error) {
	//check
	if docIds == nil || len(docIds) <= 0 {
		return nil, errors.New("invalid parameter")
	}
	if f.index == nil {
		return nil, errors.New("inter index not init")
	}

	//setup filter
	filterBuff := bytes.NewBuffer(nil)
	i := 0
	for _, docId := range docIds {
		if docId == "" {
			continue
		}
		if i > 0 {
			filterBuff.WriteString(" OR ")
		}
		filterBuff.WriteString(fmt.Sprintf("%v = %v", condField, docId))
		i++
	}

	//setup doc query
	dq := &meilisearch.DocumentsQuery{
		Filter: []string{filterBuff.String()},
	}
	resp := &meilisearch.DocumentsResult{
		Results: []map[string]interface{}{},
	}

	//get real doc
	err := f.index.GetDocuments(dq, resp)
	if err != nil || resp == nil {
		return nil, err
	}
	return resp.Results, err
}

//get one doc by field condition
//sync opt
func (f *Doc) GetOneDocByFieldCond(
	filters interface{},
	out interface{}) error {
	//check
	if filters == nil || out == nil {
		return errors.New("invalid parameter")
	}
	if f.index == nil {
		return errors.New("inter index not init")
	}

	//setup search request
	sq := &meilisearch.SearchRequest{
		Filter: filters,
		Offset: 0,
		Limit: 1,
		HitsPerPage:1,
		//AttributesToSearchOn:[]string{matchField},
	}

	//get origin doc
	resp, subErr := f.index.Search("", sq)
	if subErr != nil || resp == nil ||
		resp.Hits == nil || len(resp.Hits) <= 0 {
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
	err := f.index.GetDocument(docId, nil, &out)
	return err
}

//del one doc
//dataId used for pick hashed son worker
func (f *Doc) DelDoc(
	dataId string,
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
	_, err := f.worker.SendData(req, dataId)
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
//dataIds used for pick hashed son worker
func (f *Doc) UpdateDoc(
	docObj interface{},
	dataIds ...string) error {
	var (
		dataId string
	)
	//check
	if docObj == nil {
		return errors.New("invalid parameter")
	}
	if f.index == nil {
		return errors.New("inter index not init")
	}
	if dataIds != nil && len(dataIds) > 0 {
		dataId = dataIds[0]
	}

	//init request
	req := syncDocReq{
		obj: docObj,
		isUpdate: true,
	}

	//send worker queue
	_, err := f.worker.SendData(req, dataId)
	return err
}

//add one or batch doc
//dataIds used for pick hashed son worker
func (f *Doc) AddDoc(
	docObj interface{},
	dataIds ...string) error {
	var (
		dataId string
	)
	//check
	if docObj == nil {
		return errors.New("invalid parameter")
	}
	if f.index == nil {
		return errors.New("inter index not init")
	}
	if dataIds != nil && len(dataIds) > 0 {
		dataId = dataIds[0]
	}

	//init request
	req := syncDocReq{
		obj: docObj,
	}

	//send worker queue
	_, err := f.worker.SendData(req, dataId)
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
	if resp == nil {
		return errors.New("no any response from meili search")
	}

	//wait for task status
	finalTask, subErr := f.client.WaitForTask(resp.TaskUID)
	if subErr != nil {
		return subErr
	}
	if finalTask.Status != "succeeded" {
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
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, errors.New("no any response from meili search")
	}

	//wait for task status
	finalTask, subErr := f.client.WaitForTask(resp.TaskUID)
	if subErr != nil {
		return nil, subErr
	}
	if finalTask.Status != "succeeded" {
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