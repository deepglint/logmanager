package models

import (
	//"bufio"
	//"bytes"
	//"encoding/json"
	"errors"
	//"fmt"
	"github.com/deepglint/glog"
	//"io"
	"mime/multipart"
	"net/http"
	"net/url"
	//"os"
	"path"
	"strings"
	"time"
)

type Config struct {
	URL             url.URL
	Username        string
	Password        string
	UserAgent       string
	Timeout         time.Duration
	RetentionPolicy string
	Precision       string
	Consistency     string
}

type LogForwarder struct {
	url             url.URL
	username        string
	password        string
	httpClient      *http.Client
	userAgent       string
	retentionPolicy string
	precision       string
	consistency     string
}

type LogModel struct {
	Type    string
	Time    string
	Debug   string
	Process string
	Pid     string
	Message string
}

type Query struct {
	Command  string
	Database string
}

func NewLogForwarder(c Config) (*LogForwarder, error) {
	logforwarder := &LogForwarder{
		url:             c.URL,
		username:        c.Username,
		password:        c.Password,
		userAgent:       c.UserAgent,
		httpClient:      &http.Client{Timeout: c.Timeout},
		retentionPolicy: c.RetentionPolicy,
		precision:       c.Precision,
		consistency:     c.Consistency,
	}
	if logforwarder.userAgent == "" {
		logforwarder.userAgent = "InfluxDBLogForwarder"
	}
	return logforwarder, nil
}

func (this *LogForwarder) SaveFile(file multipart.File, handler *multipart.FileHeader) (int, error) {
	_, _, err := this.Ping()
	if err != nil {
		return 500, err
	}
	filename := path.Base(handler.Filename)
	//fmt.Println(filename)
	var database string
	if strings.ContainsAny(filename, ".") {
		fields := strings.Split(filename, ".")
		database = fields[1]
	} else {
		database = filename
	}

	q := createDatabase(database)
	_, err = this.Query(q)
	if err != nil {
		return 400, err
	}
	q = alterRetentionPolicy(database, this.retentionPolicy)
	_, err = this.Query(q)
	if err != nil {
		return 400, err
	}

	_, write_err := this.Write(file, database)

	if write_err != nil {
		/*fd, err := os.OpenFile("/tmp/logserver/"+database, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			glog.Errorf("Create tmp log file failed: %v", err)
			return 400, err
		}
		defer fd.Close()
		file.Seek(0, 0)
		_, err = io.Copy(fd, file) // file used once, so need to reset to pointer to start
		if err != nil {
			glog.Errorf("Copy to tmp log file failed: %v", err)
			return 400, err
		}*/
		return 400, write_err
	}
	return 204, nil
}

func (this *LogForwarder) Ping() (time.Duration, string, error) {
	now := time.Now()
	url := this.url
	url.Path = "ping"
	// fmt.Println(url.String())

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		glog.Errorf("%v", err)
		return 0, "", err
	}
	req.Header.Set("User-Agent", this.userAgent)
	if this.username != "" {
		req.SetBasicAuth(this.username, this.password)
	}

	resp, err := this.httpClient.Do(req)
	if err != nil {
		glog.Errorf("%v", err)
		return 0, "", err
	}
	defer resp.Body.Close()
	//fmt.Println(resp.StatusCode)
	version := resp.Header.Get("X-Influxdb-Version")
	return time.Since(now), version, nil
}

func (this *LogForwarder) Query(q *Query) (*http.Response, error) {
	//fmt.Println(q.Command)
	//fmt.Println(q.Database)
	url := this.url
	url.Path = "query"
	params := url.Query()
	params.Set("q", q.Command)
	if q.Database != "" {
		params.Set("db", q.Database)
	}
	url.RawQuery = params.Encode()
	// fmt.Println(url.String())

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		glog.Errorf("%v", err)
		return nil, err
	}
	req.Header.Set("User-Agent", this.userAgent)
	if this.username != "" {
		req.SetBasicAuth(this.username, this.password)
	}

	resp, err := this.httpClient.Do(req)
	if err != nil {
		glog.Errorf("%v", err)
		return nil, err
	}
	defer resp.Body.Close()
	//fmt.Println(resp.StatusCode)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		query_err := errors.New("Query database failed")
		glog.Errorf("%v", query_err)
		return nil, query_err
	}
	return resp, nil
}

func (this *LogForwarder) Write(file multipart.File, filename string) (*http.Response, error) {
	url := this.url
	url.Path = "write"

	//fields := strings.Split(filename, ".")

	params := url.Query()
	params.Set("consistency", this.consistency)
	params.Set("db", filename)
	params.Set("precision", this.precision)
	//params.Set("rp", "default")
	url.RawQuery = params.Encode()
	// fmt.Println(url.String())

	//Send the multipart.file to influxdb directly
	req, err := http.NewRequest("POST", url.String(), file)

	//io.copy multipart.file to fd (unnecessary)
	/*
		fd, err := os.OpenFile("/tmp/logserver/"+filename, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			glog.Errorf("Create tmp log file failed: %v", err)
			return nil, err
		}
		defer fd.Close()
		_, err = io.Copy(fd, file)
		if err != nil {
			glog.Errorf("Copy to tmp log file failed: %v", err)
			return nil, err
		}
		fd.Seek(0, 0)
		//	req, err := http.NewRequest("POST", url.String(), fd)

		//Resolve JSON to LINE protocol, and sent to influxdb
		scanner := bufio.NewScanner(fd)
		bodyBuf := &bytes.Buffer{}
		for scanner.Scan() {
			var log LogModel
			//fmt.Println(scanner.Bytes())
			str := scanner.Bytes()
			err = json.Unmarshal(str[:], &log)
			if err != nil {
				continue
			}
			query := log.Process + ",process=" + log.Process + ",type=" + log.Type + " value=\"" + log.Pid + "\",message=\"" + log.Message + "\",debug=\"" + log.Debug + "\" " + log.Time
			//fmt.Println(query)
			bodyBuf.WriteString(query)
			bodyBuf.WriteByte('\n')
		}
		req, err := http.NewRequest("POST", url.String(), bodyBuf)
	*/
	if err != nil {
		glog.Errorf("%v", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "")
	req.Header.Set("User-Agent", this.userAgent)
	if this.username != "" {
		req.SetBasicAuth(this.username, this.password)
	}

	resp, err := this.httpClient.Do(req)
	if err != nil {
		return nil, err
	} else {
		//os.Remove("/tmp/logserver/" + filename)
	}
	defer resp.Body.Close()
	//fmt.Println(resp.StatusCode)
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		write_err := errors.New("Write to influxdb failed, bad request format")
		glog.Errorf("%v", write_err)
		return nil, write_err
	}
	return resp, nil
}

/*
func Save(logs []LogModel) error {
	cnt := 0
	for _, value := range logs {
		query := value.process + ",type=" + value.logtype + ",process=" + value.process + " value=\"" + value.pid + "\" message=\"" + value.message + "\" debug=\"" + value.debug + "\" "
		resp, err := http.Post("http://localhost:8086/write?db=jsonlog", "application/x-www-form-urlencoded", strings.NewReader(query))
		if err != nil || resp.StatusCode != 204 {
			glog.Errorf("Http post to database failed: %v", err)
		} else {
			cnt++
		}
		defer resp.Body.Close()
	}
	if cnt == len(logs) {
		return nil
	} else {
		return errors.New("Http post to database incomplete")
	}
}
*/

func createDatabase(database string) *Query {
	cmd := "CREATE DATABASE " + database
	q := &Query{
		Command:  cmd,
		Database: "",
	}
	return q
}

func alterRetentionPolicy(database, duration string) *Query {
	cmd := "ALTER RETENTION POLICY default ON " + database + " DURATION " + duration
	q := &Query{
		Command:  cmd,
		Database: "",
	}
	return q
}
