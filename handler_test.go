package ginlogger_test

import (
	"testing"
	"net/http"
	"io/ioutil"
	"github.com/ilackarms/ginlogger"
	"os"
	"net/http/httptest"
	"bytes"
)

func TestLogger(t *testing.T) {
	echo := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request){
		p, err := ioutil.ReadAll(req.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed reading request"))
		}
		defer req.Body.Close()
		w.Write(p)
	})
	loggedEcho := ginlogger.Handler(echo, os.Stdout)
	server := httptest.NewServer(loggedEcho)
	body := []byte("hi there bub")
	resp, err := http.Post(server.URL, "text/plain", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("err doing request: %v", err)
	}
	p, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed reading response: %v", err)
	}
	defer resp.Body.Close()
	if !bytes.Equal(p, body) {
		t.Fatalf("actual: %s, expected: %s", p, body)
	}
}