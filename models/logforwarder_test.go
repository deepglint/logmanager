package models

import (
	//"bufio"
	//"encoding/json"
	//"errors"
	//"fmt"
	//"github.com/deepglint/glog"
	//"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	//"strings"
	//"bytes"
	"encoding/json"
	"testing"
	"time"
)

func TestNewLogForwarder(t *testing.T) {
	config := Config{
		URL:             url.URL{},
		Username:        "",
		Password:        "",
		UserAgent:       "",
		Timeout:         time.Duration(5) * time.Second,
		RetentionPolicy: "",
		Precision:       "n",
		Consistency:     "one",
	}
	_, err := NewLogForwarder(config)
	if err != nil {
		t.Logf("unexpected error. expected %v, actual %v", nil, err)
	}
}

func TestSaveFile(t *testing.T) {
	server := respTestServer()
	defer server.Close()

	fd, _ := os.Open("test")
	var file multipart.File = fd

	url, _ := url.Parse(server.URL)
	configs := []Config{
		{URL: *url, Username: "admin"},
		{URL: *url},
		{URL: *url, Username: "admin", Timeout: time.Duration(1) * time.Nanosecond},
		{Username: "admin"},
	}
	for _, config := range configs {
		Logforwarder, err := NewLogForwarder(config)

		filenames := []string{
			".test.1734",
			"./home/.test.1734",
			"test",
		}
		for _, filename := range filenames {

			handler := multipart.FileHeader{Filename: filename}

			_, err = Logforwarder.SaveFile(file, &handler)
			if err != nil {
				t.Logf("unexpected error. expected %v, actual %v", nil, err)
			}
		}
	}
}

func TestPing(t *testing.T) {
	server := emptyTestServer()
	defer server.Close()

	url, _ := url.Parse(server.URL)
	configs := []Config{
		{URL: *url, Username: "admin"},
		{URL: *url},
		{URL: *url, Username: "admin", Timeout: time.Duration(1) * time.Nanosecond},
		{Username: "admin"},
	}
	for _, config := range configs {
		Logforwarder, err := NewLogForwarder(config)
		if err != nil {
			t.Logf("unexpected error. expected %v, actual %v", nil, err)
		}
		duration, version, err := Logforwarder.Ping()
		if err != nil {
			t.Logf("unexpected error. expected %v, actual %v", nil, err)
		}
		if duration == 0 {
			t.Logf("Timeout set to zero")
		}
		if version != "x.x" {
			t.Logf("unexpected error. expected %s, actual %v", "x.x", err)
		}
	}
}

func TestQuery(t *testing.T) {
	server := respTestServer()
	defer server.Close()

	url, _ := url.Parse(server.URL)
	configs := []Config{
		{URL: *url, Username: "admin"},
		{URL: *url},
		{URL: *url, Username: "admin", Timeout: time.Duration(1) * time.Nanosecond},
		{Username: "admin"},
	}
	for _, config := range configs {
		Logforwarder, err := NewLogForwarder(config)
		if err != nil {
			t.Logf("unexpected error. expected %v, actual %v", nil, err)
		}

		queries := []Query{
			{Database: "mydb"},
			{},
		}
		for _, query := range queries {
			_, err = Logforwarder.Query(&query)
			if err != nil {
				t.Logf("unexpected error. expected %v, actual %v", nil, err)
			}
		}
	}
}

func TestWrite(t *testing.T) {
	server := respTestServer()
	defer server.Close()

	fd, _ := os.Open("test")
	//req, err := http.NewRequest("POST", server.URL, fd)
	//req.Body = ioutil.NopCloser(bytes.NewReader(ar))
	//_, _, err = req.FormFile("testtest")
	//defer file.Close()
	/*slice := []byte{'t', 'e', 's', 't'}
	//r := ReadHelper{body: slice}
	ra := ReadatHelper{body: slice}
	s := SeekHelper{body: slice}
	c := CloseHelper{body: slice}
	//var ir io.Reader = &r
	var ira io.ReaderAt = &ra
	var is io.Seeker = &s
	var ic io.Closer = &c
	ffile := FileHelper{
		fir:  fd,
		fira: ira,
		fis:  is,
		fic:  ic,
	}*/
	// os.File implements interface multipart.File's all methods
	var file multipart.File = fd

	url, _ := url.Parse(server.URL)
	configs := []Config{
		{URL: *url, Username: "admin"},
		{URL: *url},
		{URL: *url, Username: "admin", Timeout: time.Duration(1) * time.Nanosecond},
		{Username: "admin"},
	}
	for _, config := range configs {
		logforwarder, err := NewLogForwarder(config)
		if err != nil {
			t.Logf("unexpected error. expected %v, actual %v", nil, err)
		}
		cases := []string{
			"test",
			".test.1734",
		}

		for _, filename := range cases {
			//logforwarder.Ping()
			//t.Logf(filename)
			_, err = logforwarder.Write(file, filename)
			//	if err != nil {
			//		t.Logf("unexpected error. expected %v, actual %v", nil, err)
			//	}
		}
	}
}

// helper functions

func emptyTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Influxdb-Version", "x.x")
		return
	}))
}

func respTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var resp http.Response
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	}))
}

/*
type ReadHelper struct {
	body []byte
}

func (this *ReadHelper) Read(p []byte) (n int, err error) {
	//this.body = p
	return len(p), nil
}

func (this *ReadHelper) Close() error {
	return nil
}

type ReadatHelper struct {
	body []byte
}

func (this *ReadatHelper) ReadAt(p []byte, off int64) (n int, err error) {
	//this.body = p[off:]
	return 0, nil
}

type SeekHelper struct {
	body []byte
}

func (this *SeekHelper) Seek(offset int64, whence int) (n int64, err error) {
	return 0, nil
}

type CloseHelper struct {
	body []byte
}

func (this *CloseHelper) Close() error {
	//this.body = nil
	return nil
}

type FileHelper struct {
	multipart.File
	fir  io.Reader
	fira io.ReaderAt
	fis  io.Seeker
	fic  io.Closer
}*/

func TestalterRetentionPolicy(t *testing.T) {
	alterRetentionPolicy("test", "7d")
}

func TestcreateDatabase(t *testing.T) {
	createDatabase("test")
}
