package conf

import "time"

type (
	IndexConf struct {
		IndexName        string //must value
		PrimaryKey       string //must value
		FilterableFields []string
	}
	ClientConf struct {
		Tag         string
		Host        string
		ApiKey      string
		TimeOut     time.Duration
		IndexesConf []*IndexConf //indexes config
		Workers     int          //inter concurrency workers
	}
)
