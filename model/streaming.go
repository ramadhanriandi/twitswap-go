package model

import "time"

type Streaming struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	RuleID    string    `json:"rule_id"`
}
