package conf

import "time"

type (
	IndexConf struct {
		IndexName        string //must value
		PrimaryKey       string //must value
		FilterableFields []string
		SortableFields   []string
		CreateIndex      bool
		UpdateFields 	 bool
	}
	ClientConf struct {
		Tag         string
		Host        string
		ApiKey      string
		TimeOut     time.Duration
		IndexesConf []*IndexConf //indexes config
		CreateIndex bool         //create index or not
		Workers     int          //inter concurrency workers
	}
)
