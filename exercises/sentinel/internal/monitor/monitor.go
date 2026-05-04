package monitor

import (
	"fmt"
	"my-go-learning/sentinel/internal/storage"
	"my-go-learning/sentinel/internal/types"
	"net/http"
	"sync"
	"time"
)

type Monitor struct {
	websiteList []types.Website
	storage     *storage.Storage
}

func (m Monitor) checkWebsite(w types.Website, wg *sync.WaitGroup) {
	defer wg.Done()
	d, e := http.Get(w.URL)
	if e != nil {
		fmt.Printf("error %s", e)
		return
	}

	cr := types.CheckResult{WebsiteID: w.Id, Time: time.Now(), Status: d.StatusCode == 200}

	m.storage.SaveResult(cr)
}

func NewMonitor(list []types.Website, s *storage.Storage) *Monitor {
	return &Monitor{
		websiteList: list,
		storage:     s,
	}
}

func (m Monitor) Start() {
	var wg sync.WaitGroup

	ticker := time.NewTicker(10 * time.Second)

	for range ticker.C {
		for _, website := range m.websiteList {
			wg.Add(1)
			go m.checkWebsite(website, &wg)
		}
		wg.Wait()
		fmt.Println("Tick finished!")
	}
}
