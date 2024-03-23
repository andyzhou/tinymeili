package face

import (
	"encoding/json"
	"errors"
	"github.com/andyzhou/tinymeili/define"
	"github.com/meilisearch/meilisearch-go"
)

/*
 * doc opt face
 */

//face info
type Doc struct {
	index *meilisearch.Index //reference
}

//construct
func NewDoc(index *meilisearch.Index) *Doc {
	this := &Doc{
		index: index,
	}
	return this
}

//query batch doc one index
//return total, []docObj, error
func (f *Doc) QueryIndexDocs(
		key string,
		filter interface{},
		sort []string,
		page, pageSize int,
	) (int64, []interface{}, error) {
	//check
	if f.index == nil {
		return 0, nil, errors.New("inter index not init")
	}

	//setup offset
	if page <= 0 {
		page = define.DefaultPage
	}
	offset := (page - 1) * pageSize
	limit := pageSize

	//setup search request
	sq := &meilisearch.SearchRequest{
		Query: key,
		Filter: filter,
		Sort: sort,
		Offset: int64(offset),
		Limit: int64(limit),
		HitsPerPage:int64(pageSize),
	}

	//query origin doc
	resp, subErr := f.index.Search(key, sq)
	if subErr != nil || resp == nil {
		return 0, nil, subErr
	}
	return resp.TotalHits, resp.Hits, nil
}

//get one doc by field condition
func (f *Doc) GetOneDocByFieldCond(
	docId, matchField string,
	out interface{}) error {
	//check
	if docId == "" ||
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
	resp, subErr := f.index.Search("", sq)
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
		docIds ...string,
	) error {
	//check
	if docIds == nil || len(docIds) <= 0 {
		return errors.New("invalid parameter")
	}
	if f.index == nil {
		return errors.New("inter index not init")
	}

	//del real doc
	_, err := f.index.DeleteDocuments(docIds)
	return err
}

//del docs by filter
//filter like: 'a = 6 and b < 10'
func (f *Doc) DelDocsByFilter(
		filter interface{},
	) error {
	//check
	if filter == nil {
		return errors.New("invalid parameter")
	}
	if f.index == nil {
		return errors.New("inter index not init")
	}

	//del real docs by filter
	_, err := f.index.DeleteDocumentsByFilter(filter)
	return err
}

//update one doc
func (f *Doc) UpdateDoc(
	docObj interface{}) error {
	//check
	if docObj == nil {
		return errors.New("invalid parameter")
	}
	if f.index == nil {
		return errors.New("inter index not init")
	}

	//update real doc
	_, err := f.index.UpdateDocuments(docObj)
	return err
}

//add one or batch doc
func (f *Doc) AddDoc(
	docObj interface{}) error {
	//check
	if docObj == nil {
		return errors.New("invalid parameter")
	}
	if f.index == nil {
		return errors.New("inter index not init")
	}

	//add real doc
	_, err := f.index.AddDocuments(docObj)
	return err
}