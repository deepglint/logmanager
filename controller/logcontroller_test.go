package controller

import (
	//"encoding/json"
	//"fmt"
	//"github.com/deepglint/glog"
	//"bytes"
	"github.com/deepglint/logmanager/models"
	"github.com/deepglint/muses/util/ripple"
	//"io"
	"os"
	//"path/filepath"
	"encoding/json"
	//"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestString(t *testing.T) {
	config := models.Config{}
	logcontroller, _ := NewLogController("log", config)
	logcontroller.String()
}

func TestNewLogController(t *testing.T) {
	config := models.Config{
		URL:             url.URL{},
		Username:        "",
		Password:        "",
		UserAgent:       "",
		Timeout:         time.Duration(5) * time.Second,
		RetentionPolicy: "",
		Precision:       "n",
		Consistency:     "one",
	}
	_, err := NewLogController("log", config)
	if err != nil {
		t.Logf("unexpected error. expected %v, actual %v", nil, err)
	}
}

func TestPostUpload(t *testing.T) {
	server := respTestServer()
	defer server.Close()

	/*bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	fileWriter, _ := bodyWriter.CreateFormFile("uploadlogfile", "testtest")
	fh, _ := os.OpenFile("testtest", os.O_RDWR|os.O_CREATE, 0666)
	_, _ = io.Copy(fileWriter, fh)
	contentType := bodyWriter.FormDataContentType()
	defer fh.Close()*/

	fd, _ := os.OpenFile("test", os.O_RDWR|os.O_CREATE, 0666)
	defer fd.Close()

	url, _ := url.Parse(server.URL)
	req, _ := http.NewRequest("POST", url.String(), fd)
	req.Header.Set("Content-Type", "multipart/form-data; boundary=a8f16fe282c995eee9ad34a1c955e517327dba7383cab1e528286eccffdd")

	ctx := ripple.NewContext()
	ctx.Request = req

	configs := []models.Config{
		{URL: *url, Username: "admin"},
		{URL: *url},
		{URL: *url, Username: "admin", Timeout: time.Duration(1) * time.Nanosecond},
		{Username: "admin"},
	}

	for _, config := range configs {
		logcontroller, _ := NewLogController("log", config)
		logcontroller.PostUpload(ctx)
	}
}

func respTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var resp http.Response
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	}))
}
