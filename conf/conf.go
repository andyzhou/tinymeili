package conf

import "time"

type (
	ClientConf struct {
		Tag     string
		Host    string
		ApiKey  string
		TimeOut time.Duration
		Indexes []string 		//index names
		Workers int      		//inter concurrency workers
	}
)
