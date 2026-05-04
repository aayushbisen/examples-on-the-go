package logprocessor

import (
	"fmt"
	"net/http"
	"strings"
)

func myhandler(inp chan string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		msg := r.URL.Query().Get("msg")
		inp <- msg
		fmt.Fprintf(w, "Log Received!")
	}
}

func process(ch chan string, opch chan string) {
	for msg := range ch {
		msgU := strings.ToUpper(msg)
		opch <- msgU
	}
}

func printer(ch chan string) {
	for msg := range ch {
		fmt.Println(msg)
	}
}

func main() {
	inputChan := make(chan string)
	outputChan := make(chan string)

	go process(inputChan, outputChan)

	go printer(outputChan)

	http.HandleFunc("/log", myhandler(inputChan))
	http.ListenAndServe(":8080", nil)
}
