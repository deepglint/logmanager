package client

import (
	"bytes"
	// "encoding/json"
	"errors"
	"fmt"
	"io"
	// "io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/deepglint/glog"
	"github.com/hpcloud/tail"
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
		glog.Errorf("Error opening file", err)
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
	glog.Errorln(resp, err)
	defer resp.Body.Close()

	if resp.StatusCode != 204 && resp.StatusCode != 200 {
		if resp.StatusCode == 411 {
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
		if b, _ := path.Match("LOG.*.*????-??-??T??:??:??Z", f.Name()); b {
			fields := strings.Split(f.Name(), ".")
			t, _ := time.Parse("MST2006-01-02T15:04:05Z", fields[2])
			time_int := t.Unix()
			// glog.Infoln("Log created at: ", t, "\tTime now: ", time.Now())
			// fmt.Println("Log created at: ", t, "\tTime now: ", time.Now())
			if time.Now().Unix()-int64(upload.Seconds()) > time_int {
				post_err := PostFile(filename, url)
				if post_err != nil {
					glog.Errorf("Upload file failed, %v", post_err)
				} else {
					os.Rename(filename, path.Dir(filename)+"/OLD."+f.Name())
					glog.Infof("%s upload succeed \n", filename)
					// fmt.Printf("%s upload succeed \n", filename)
				}
			}
			return nil
		}
		if b, _ := path.Match("OLD.LOG.*.*????-??-??T??:??:??Z", f.Name()); b {
			fields := strings.Split(f.Name(), ".")
			t, _ := time.Parse("MST2006-01-02T15:04:05Z", fields[3])
			time_int := t.Unix()
			if time.Now().Unix()-int64(keep.Seconds()) > time_int {
				err := os.Remove(filename)
				if err == nil {
					glog.Infof("%s removed \n", filename)
					// fmt.Printf("%s removed \n", filename)
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
		if b, _ := path.Match("LOG.*.*????-??-??T??:??:??Z", f.Name()); b {
			post_err := PostFile(filename, url)
			if post_err != nil {
				glog.Errorf("Upload file failed, %v", post_err)
			} else {
				glog.Infof("%s upload succeed \n", filename)
				// fmt.Printf("%s upload succeed \n", filename)
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

func shortHostname(hostname string) string {
	if i := strings.Index(hostname, "."); i >= 0 {
		return hostname[:i]
	}
	return hostname
}

type SensorId struct {
	Key           string
	Value         string
	ModifiedIndex int
	CreatedIndex  int
}

type Sensor struct {
	Action string
	Node   SensorId
}

func GetHost() (string, error) {
	// var host string
	// resp, err := http.Get("http://localhost:4001/v2/keys/config/global/sensor_uid")
	// if err != nil || resp.StatusCode != 200 {
	// 	fmt.Printf("%v", err)
	// 	// fmt.Println(os.Hostname())
	//
	// 	if err == nil {
	// 		host = shortHostname(h)
	// 	}
	// } else {
	// 	defer resp.Body.Close()
	// 	body, err := ioutil.ReadAll(resp.Body)
	// 	// fmt.Printf("%s", string(body))
	// 	var sen Sensor
	// 	err = json.Unmarshal(body, &sen)
	// 	if err == nil {
	// 		host = sen.Node.Value
	// 		return host, nil
	// 	}
	// }
	h, _ := os.Hostname()
	return h, nil
}

func writeFile(filename, text string) {
	fd, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		glog.Errorf("%v", err)
		return
	}
	defer fd.Close()
	fd.WriteString(text)
	fd.WriteString("\n")
}

func TailLog(host string, interval time.Duration, dir string, limit int) {
	fmt.Println("START TAILING LOG.......")
	count := 0
	for {
	NEW_TAIL_FILE:
		tail_file := fmt.Sprintf("%sLOG.%4d-%02d-%02dT00:00:00Z", dir, time.Now().Year(), time.Now().Month(), time.Now().Day())
		fmt.Printf("Tailing %v\n", tail_file)
		t, err := tail.TailFile(tail_file, tail.Config{Follow: true, MustExist: true})
		if err != nil {
			glog.Errorf("Tail error: %v", err)
			// fmt.Printf("Tail error: %v", err)
			time.Sleep(time.Second * 30)
			continue
		}
		glog.Infoln(tail_file)
		var fname_now, fname_old string = "", ""
	NEW_ROTATE_FILE:
		count = 0
		for {
			select {
			case line := <-t.Lines:
				if count < limit {
					tmp := dir + "LOG." + host + "." + time.Now().Round(interval).Format("MST2006-01-02T15:04:05Z")
					if fname_old != tmp {
						count = 0
						fname_old = tmp
					}
					writeFile(fname_old, line.Text)
					count++
				} else {
					new_tail_file := fmt.Sprintf("%sLOG.%4d-%02d-%02dT00:00:00Z", dir, time.Now().Year(), time.Now().Month(), time.Now().Day())
					if new_tail_file != tail_file {
						glog.Infoln("New Tail File, go to NEW_TAIL_FILE")
						t.Stop()
						goto NEW_TAIL_FILE
					}
					glog.Warningln("Reach line number limit, no more write to current file.")
					for _ = range time.NewTicker(30 * time.Second).C {
						fname_now = dir + "LOG." + host + "." + time.Now().Round(interval).Format("MST2006-01-02T15:04:05Z")
						if fname_old != fname_now {
							glog.Infoln("New rotate log file, go to NEW_ROTATE_FILE")
							goto NEW_ROTATE_FILE
						} else {
							glog.Infoln("Reach line number limit, ignore")
						}
					}
				}
			case <-time.After(time.Minute * 5):
				new_tail_file := fmt.Sprintf("%sLOG.%4d-%02d-%02dT00:00:00Z", dir, time.Now().Year(), time.Now().Month(), time.Now().Day())
				if new_tail_file != tail_file {
					glog.Infoln("New Tail File, go to NEW_TAIL_FILE")
					t.Stop()
					goto NEW_TAIL_FILE
				} else {
					glog.Infoln("Go to LOOP")
				}
			}
		}
	}
}

func CleanLog(dir string, tt time.Duration) error {
	err := filepath.Walk(dir, func(filename string, f os.FileInfo, err error) error {
		if f == nil {
			info_err := errors.New("Filepath.Walk() returned no fileinfo")
			glog.Errorf("%v", info_err)
			return nil
		}
		if f.IsDir() {
			return nil
		}
		if b, _ := path.Match("LOG.????-??-??T??:??:??Z", f.Name()); b {
			glog.Infof("Processing %v", f.Name())
			fields := strings.Split(f.Name(), ".")
			if len(fields) < 2 {
				return nil
			}
			t, _ := time.Parse("2006-01-02T15:04:05Z", fields[1])
			time_int := t.Unix()
			if time.Now().Unix()-int64(tt.Seconds()) > time_int {
				err := os.Remove(filename)
				if err == nil {
					glog.Infof("%s removed", filename)
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
