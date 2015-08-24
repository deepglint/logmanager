package client

import (
	"bytes"
	"fmt"
	"io"
	//"io/ioutil"
	"errors"
	"github.com/deepglint/glog"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

func PostFile(filename, targetUrl string) error {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	fileWriter, err := bodyWriter.CreateFormFile("uploadlogfile", filename)
	if err != nil {
		glog.Errorf("Error writing to buffer")
		return err
	}
	//var mu sync.Mutex
	//mu.Lock()
	//glog.Flush()
	fh, err := os.Open(filename)
	if err != nil {
		glog.Errorf("Error opening file")
		return err
	}
	defer fh.Close()

	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		glog.Errorf("%v", err)
		return err
	}
	//mu.Unlock()
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(targetUrl, contentType, bodyBuf)
	if err != nil {
		glog.Errorf("%v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 204 && resp.StatusCode != 200 {
		if resp.StatusCode == 400 {
			code_err := errors.New("Bad filename format")
			glog.Errorf("%s: %v %d", filename, code_err, resp.StatusCode)
			return code_err
		} else if resp.StatusCode == 500 {
			code_err := errors.New("Log server internal error")
			glog.Errorf("%v %d", code_err, resp.StatusCode)
			return code_err
		} else {
			code_err := errors.New("Http status code error")
			glog.Errorf("%v %d", code_err, resp.StatusCode)
			return code_err
		}
	}
	return nil
}

func SendLog(url, dir string, upload, keep time.Duration) error {
	err := filepath.Walk(dir, func(filename string, f os.FileInfo, err error) error {
		if f == nil {
			info_err := errors.New("Filepath.Walk() returned no fileinfo")
			glog.Errorf("%v", info_err)
			return nil
		}
		if f.IsDir() {
			return nil
		}
		if b, _ := path.Match("LOG.*.*????-??-??-??:??:??", f.Name()); b {
			fields := strings.Split(f.Name(), ".")
			t, _ := time.Parse("MST2006-01-02-15:04:05", fields[2])
			time_int := t.Unix()
			// glog.Infoln("Log created at: ", t, "\tTime now: ", time.Now())
			// fmt.Println("Log created at: ", t, "\tTime now: ", time.Now())
			if time.Now().Unix()-int64(upload.Seconds()) > time_int {
				post_err := PostFile(filename, url)
				if post_err != nil {
					// glog.Errorf("Upload file failed, %v", post_err)
				} else {
					os.Rename(filename, path.Dir(filename)+"/OLD."+f.Name())
					glog.Infof("%s upload succeed \n", filename)
					fmt.Printf("%s upload succeed \n", filename)
				}
			}
			return nil
		}
		if b, _ := path.Match("OLD.LOG.*.ULAT????-??-??-??:??:??", f.Name()); b {
			fields := strings.Split(f.Name(), ".")
			t, _ := time.Parse("MST2006-01-02-15:04:05", fields[3])
			time_int := t.Unix()
			if time.Now().Unix()-int64(keep.Seconds()) > time_int {
				err := os.Remove(filename)
				if err == nil {
					glog.Infof("%s removed \n", filename)
					fmt.Printf("%s removed \n", filename)
				}
			}
			return nil
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func SendLogNow(url, dir string) error {
	err := filepath.Walk(dir, func(filename string, f os.FileInfo, err error) error {
		if f == nil {
			info_err := errors.New("Filepath.Walk() returned no fileinfo")
			glog.Errorf("%v", info_err)
			return nil
		}
		if f.IsDir() {
			return nil
		}
		if b, _ := path.Match("LOG.*.*????-??-??-??:??:??", f.Name()); b {
			post_err := PostFile(filename, url)
			if post_err != nil {
				glog.Errorf("Upload file failed, %v", post_err)
			} else {
				glog.Infof("%s upload succeed \n", filename)
				fmt.Printf("%s upload succeed \n", filename)
			}
			return nil
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func Flush(url, dir string) error {
	err := SendLogNow(url, dir)
	return err
}
