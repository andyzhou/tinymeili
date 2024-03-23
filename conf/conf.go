package conf

import "time"

type (
	NodeConf struct {
		Kind    string
		Hosts   map[string]string //tag -> host
		ApiKey  string
		TimeOut time.Duration
		Indexes []string
	}
)
