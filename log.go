package main

import (
	"encoding/json"
	"fmt"
)

type Log struct {
	Time        string
	Level       string
	Source      Source
	Diagnostics string
}

type Source struct {
	Function string
	File     string
	Line     int
}

func JSONParser(raw string) (*Log, error) {
	log := &Log{}
	if err := json.Unmarshal([]byte(raw), log); err != nil {
		return nil, err
	}
	return log, nil
}

func (l *Log) String() string {
	return fmt.Sprintf("Log { Timestamp: %s, Level: %s, Diagnostics: %s, Source: %+v}",
		l.Time, l.Level, l.Diagnostics, l.Source,
	)

}
