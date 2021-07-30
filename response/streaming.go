package response

import "time"

type StartStreaming struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	StartTime time.Time `json:"start_time"`
	RuleID    string    `json:"rule_id"`
	Rule      string    `json:"rule"`
}
