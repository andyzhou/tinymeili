package main

type ReviewDoc struct {
	Id       int64  `json:"id"` //auto inc val
	DataId   int64  `json:"dataId"`
	Parent   int64  `json:"parent"`
	Poster   int64  `json:"poster"`
	Receiver int64  `json:"receiver"`
	Content  string `json:"content"`
	File     string `json:"file"`
	Quote    int64  `json:"quote"`  //quote data id
	Topped   int64  `json:"topped"` //topped time stamp, for root review
	Praise   int64  `json:"praise"` //for root review
	Score    int64  `json:"score"`  //for root review
	Upped    bool   `json:"upped"`  //virtual field
	Status   int    `json:"status"` //0:normal 1:removed
	EditAt   int64  `json:"editAt"`
	CreateAt int64  `json:"createAt"`
	BaseJson
}

type MyReviewJson struct {
	ReviewId int64  `json:"reviewId"`
	DataId   int64  `json:"dataId"`
	Parent   int64  `json:"parent"`
	Poster   int64  `json:"poster"`
	Receiver int64  `json:"receiver"`
	Content  string `json:"content"`
	IsSent	 bool	`json:"isSent"`
	CreateAt int64  `json:"createAt"`
	BaseJson
}

type TestDoc struct {
	Id    int64    `json:"id"`
	Title string   `json:"title"`
	Tags  []string `json:"tags"`
	BaseJson
}

//construct
func NewTestDoc() *TestDoc {
	this := &TestDoc{
		Tags: []string{},
	}
	return this
}

func NewReviewDoc() *ReviewDoc {
	this := &ReviewDoc{}
	return this
}

func NewMyReviewJson() *MyReviewJson {
	this := &MyReviewJson{}
	return this
}
