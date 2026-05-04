package types

import "time"

type Website struct {
	Id   string
	Name string
	URL  string
}

type CheckResult struct {
	WebsiteID string
	Time      time.Time
	Status    bool
}
