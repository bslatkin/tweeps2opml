package main

import (
	"os"
)

var (
	Globals = map[string]string{
		"AnalyticsId": os.Getenv("ANALYTICS_ID"),
	}
)

type Params struct {
	Globals interface{}
	Data    interface{}
}
