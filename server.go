package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

type Mock struct {
	Method      string `json:"method"`
	Path        string `json:"path"`
	StatusCode  string `json:"statusCode"`
	ContentType string `json:"contentType"`
	Response    string `json:"response"`
}

type MockHandler struct {
	mutex sync.Mutex
	mock  Mock
	f     string
	p     string
}

func (h MockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(h.mock)

	if r.Method == h.mock.Method {
		w.Header().Set("Content-Type", h.mock.ContentType)
		var statusCode int
		statusCode, _ = strconv.Atoi(h.mock.StatusCode)
		w.WriteHeader(statusCode)

		w.Write([]byte(h.mock.Response))
	}
}

func main() {
	var (
		f = flag.String("f", "./mock.json", "JSON file path")
		p = flag.String("p", "8888", "Port number")
	)
	flag.Parse()

	mockJsonFile, err := os.Open(*f)
	if err != nil {
		log.Fatal(err)
	}

	defer mockJsonFile.Close()

	mockByteValue, _ := ioutil.ReadAll(mockJsonFile)
	var mocks []Mock
	json.Unmarshal(mockByteValue, &mocks)

	for _, mock := range mocks {
		var handler MockHandler
		handler.mock = mock
		handler.f = *f
		handler.p = *p

		http.Handle(mock.Path, handler)
	}

	log.Fatal(http.ListenAndServe(":"+*p, nil))
}
