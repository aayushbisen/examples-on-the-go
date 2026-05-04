package main

import (
	"fmt"
	"my-go-learning/sentinel/internal/dashboard"
	"my-go-learning/sentinel/internal/monitor"
	"my-go-learning/sentinel/internal/storage"
	"my-go-learning/sentinel/internal/types"
	"net/http"
)

func main() {

	s, e := storage.NewStorage()
	if e != nil {
		fmt.Printf("error %s", e)
		return
	}
	websites := []types.Website{
		{Id: "1", Name: "Google", URL: "https://google.com"},
		{Id: "2", Name: "GitHub", URL: "https://github.com"},
		{Id: "3", Name: "My Portfolio", URL: "https://aayushbisen.dev"},
	}

	go func() {
		http.HandleFunc("/dashboard", dashboard.Handler(s))
		http.ListenAndServe(":8080", nil)
	}()

	monit := monitor.NewMonitor(websites, s)
	monit.Start()
}
