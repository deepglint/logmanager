package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/deepglint/glog"
	"github.com/deepglint/logmanager/models"
	"github.com/deepglint/muses/util/ripple"
	//"os"
	//"path/filepath"
	//"mime/multipart"
	//"time"
)

type LogController struct {
	name               string
	log_forwarder      *models.LogForwarder
	logForwarderConfig models.Config
}

func (this *LogController) String() string {
	return fmt.Sprintf("LogControlloer: %s", this.name)
}

func (this *LogController) Version() string {
	return "LogServer v0.2.0"
}

func NewLogController(name string, c models.Config) (*LogController, error) {
	//fmt.Println("New LogController")
	ctl := &LogController{
		name: name,
		log_forwarder: &models.LogForwarder{
			Url:             c.URL,
			Username:        c.Username,
			Password:        c.Password,
			UserAgent:       c.UserAgent,
			HttpClient:      &http.Client{Timeout: c.Timeout},
			RetentionPolicy: c.RetentionPolicy,
			Precision:       c.Precision,
			Consistency:     c.Consistency,
		},
	}
	return ctl, nil
}

func (this *LogController) PostUpload(ctx *ripple.Context) *ripple.Response {
	ctx.Request.ParseMultipartForm(32 << 20)
	file, handler, err := ctx.Request.FormFile("uploadlogfile")
	if err != nil {
		return &ripple.Response{Status: 400, Body: err.Error(), ContentType: "text/plain"}
	}
	defer file.Close()

	status_code, err := this.log_forwarder.SaveFile(file, handler) // decode file and forward to database
	if err != nil {
		return &ripple.Response{Status: status_code, Body: err.Error(), ContentType: "text/plain"}
	}
	return &ripple.Response{Status: status_code}
}

func (this *LogController) GetVersion(ctx *ripple.Context) {
	ctx.Response.Status = 200
	ctx.Response.Body = this.Version()
}

type WhereClause struct {
	Module    string `json:"module"`
	Level     string `json:"level"`
	Begintime int64  `json:"begintime"`
	Endtime   int64  `json:"endtime"`
	Keyword   string `json:"keyword"`
}

type Log struct {
	Database  string       `json:"database"`
	Precision string       `json:"precision"`
	Columns   []string     `json:"columns"`
	Table     string       `json:"table"`
	Where     *WhereClause `json:"where"`
	Pagesize  int64        `json:"pagesize"`
	Pagenum   int64        `json:"pagenum"`
}

// DEPRECATED TIME FORMAT

// func formatInfluxTime(t, now time.Time) string {
// 	d := now.Sub(t)
// 	min := int64(d.Minutes())
// 	min_str := strconv.FormatInt(min, 10)
// 	return min_str + `m`
// }

// func formatTimeClause(begin, end, now time.Time) string {
// 	b := now.Sub(begin)
// 	e := now.Sub(end)
// 	b_min := int64(b.Minutes())
// 	e_min := int64(e.Minutes())
// 	if b_min == e_min {
// 		b_min += 1
// 	}
// 	b_str := strconv.FormatInt(b_min, 10)
// 	e_str := strconv.FormatInt(e_min, 10)
// 	b_str += "m"
// 	e_str += "m"

// 	time_clause := `time > now() - ` + b_str + ` and time < now() - ` + e_str
// 	return time_clause
// }

// DEPRECATED TIME FORMAT END

type InfluxdbResults struct {
	Results []InfluxdbResult `json:"results"`
}

type InfluxdbResult struct {
	Statement_id int64           `json:"statement_id"`
	Series       []InfluxdbSerie `json:"series"`
}

type InfluxdbSerie struct {
	Name    string        `json:"name"`
	Columns []string      `json:"columns"`
	Values  []interface{} `json:"values"`
}

func (this *LogController) PostGetLog(ctx *ripple.Context) *ripple.Response {
	body, _ := ioutil.ReadAll(ctx.Request.Body)
	var req Log
	err := json.Unmarshal(body, &req)
	glog.Errorln("Get log request: ", req, req.Where)
	if err != nil {
		err := fmt.Errorf("Post body format error")
		glog.Errorln(err)
		return &ripple.Response{Status: 400, Body: err.Error(), ContentType: "text/plain"}
	}

	if req.Database == "" || req.Table == "" {
		err := fmt.Errorf("Database name or table name not found")
		glog.Errorln(err)
		return &ripple.Response{Status: 400, Body: err.Error(), ContentType: "text/plain"}
	}
	if req.Precision == "" || (req.Precision != "ns" && req.Precision != "u" && req.Precision != "s" &&
		req.Precision != "m" && req.Precision != "h" && req.Precision != "ms") {

		req.Precision = "ms"
	}
	if len(req.Columns) == 0 {
		req.Columns = append(req.Columns, "*")
	}
	req.Where.Level = strings.ToUpper(req.Where.Level)
	if req.Where.Level != "ERROR" && req.Where.Level != "INFO" && req.Where.Level != "WARNING" && req.Where.Level != "FATAL" {
		req.Where.Level = ""
	}
	if req.Pagesize <= 0 {
		req.Pagesize = 100
	} else if req.Pagesize > 1000 {
		req.Pagesize = 1000
	}
	if req.Pagenum > 500 {
		req.Pagenum = 500
	} else if req.Pagenum < 1 {
		req.Pagenum = 1
	}
	offset := (req.Pagenum - 1) * req.Pagesize
	page_size_str := strconv.FormatInt(req.Pagesize, 10)
	offset_str := strconv.FormatInt(offset, 10)

	var begin, end, now time.Time
	now = time.Now().UTC()
	if req.Where.Begintime != 0 { //get local time coresponding to timestamp
		begin = time.Unix(req.Where.Begintime, 1000000000).UTC()
	} else {
		begin = now.Add(time.Hour * -24)
	}
	if req.Where.Endtime != 0 { //get local time coresponding to timestamp
		end = time.Unix(req.Where.Endtime, 1000000000).UTC()
	} else {
		end = now
	}
	if begin.After(end) {
		end = begin.Add(time.Hour * 24)
	}
	begin_str := begin.Format("2006-01-02 15:04:05")
	end_str := end.Format("2006-01-02 15:04:05")

	time_clause := `time > '` + begin_str + `' and time < '` + end_str + `'`
	// time_clause := formatTimeClause(begin, end, now)

	q := new(models.Query)
	q.Database = req.Database
	q.Epoch = req.Precision

	where_clause := ` where `
	if req.Where.Module != "" {
		where_clause += `program = '` + req.Where.Module + `' and `
	}
	if req.Where.Level != "" {
		where_clause += `type = '` + req.Where.Level + `' and `
	}
	if req.Where.Keyword != "" {
		where_clause += `message =~ /` + req.Where.Keyword + `/`
	}

	if where_clause[len(where_clause)-5:len(where_clause)] == " and " {
		where_clause = where_clause[:len(where_clause)-5]
	}

	if where_clause == ` where ` {
		q.Command = `select * from "` + req.Table + `"` + where_clause + time_clause + ` limit ` + page_size_str + ` offset ` + offset_str
	} else {
		q.Command = `select * from "` + req.Table + `"` + where_clause + ` and ` + time_clause + ` limit ` + page_size_str + ` offset ` + offset_str
	}
	glog.Errorln(q)

	body, err = this.log_forwarder.Query(q)
	if err != nil {
		glog.Errorln(err)
		return &ripple.Response{Status: 500, Body: err.Error(), ContentType: "text/plain"}
	}

	if log_rst, err := parseInfluxdbResult(body); err != nil {
		glog.Errorln(err)
		return &ripple.Response{Status: 500, Body: err.Error(), ContentType: "text/plain"}
	} else {
		return &ripple.Response{Status: 200, Body: string(log_rst)}
	}
}

func (this *LogController) PostLogCount(ctx *ripple.Context) *ripple.Response {
	body, _ := ioutil.ReadAll(ctx.Request.Body)
	var req Log
	err := json.Unmarshal(body, &req)
	glog.Errorln("Get log request: ", req, req.Where)
	if err != nil {
		err := fmt.Errorf("Post body format error")
		glog.Errorln(err)
		return &ripple.Response{Status: 400, Body: err.Error(), ContentType: "text/plain"}
	}

	if req.Database == "" || req.Table == "" {
		err := fmt.Errorf("Database name or table name not found")
		glog.Errorln(err)
		return &ripple.Response{Status: 400, Body: err.Error(), ContentType: "text/plain"}
	}
	if req.Precision == "" || (req.Precision != "ns" && req.Precision != "u" && req.Precision != "s" &&
		req.Precision != "m" && req.Precision != "h" && req.Precision != "ms") {

		req.Precision = "ms"
	}
	if len(req.Columns) == 0 {
		req.Columns = append(req.Columns, "*")
	}
	req.Where.Level = strings.ToUpper(req.Where.Level)
	if req.Where.Level != "ERROR" && req.Where.Level != "INFO" && req.Where.Level != "WARNING" && req.Where.Level != "FATAL" {
		req.Where.Level = ""
	}

	var begin, end, now time.Time
	now = time.Now().UTC()
	if req.Where.Begintime != 0 { //get local time coresponding to timestamp
		begin = time.Unix(req.Where.Begintime, 1000000000).UTC()
	} else {
		begin = now.Add(time.Hour * -24)
	}
	if req.Where.Endtime != 0 { //get local time coresponding to timestamp
		end = time.Unix(req.Where.Endtime, 1000000000).UTC()
	} else {
		end = now
	}
	if begin.After(end) {
		end = begin.Add(time.Hour * 24)
	}
	begin_str := begin.Format("2006-01-02 15:04:05")
	end_str := end.Format("2006-01-02 15:04:05")

	time_clause := `time > '` + begin_str + `' and time < '` + end_str + `'`
	// time_clause := formatTimeClause(begin, end, now)

	q := new(models.Query)
	q.Database = req.Database
	q.Epoch = req.Precision

	where_clause := ` where `
	if req.Where.Module != "" {
		where_clause += `program = '` + req.Where.Module + `' and `
	}
	if req.Where.Level != "" {
		where_clause += `type = '` + req.Where.Level + `' and `
	}
	if req.Where.Keyword != "" {
		where_clause += `message =~ /` + req.Where.Keyword + `/`
	}

	if where_clause[len(where_clause)-5:len(where_clause)] == " and " {
		where_clause = where_clause[:len(where_clause)-5]
	}

	if where_clause == ` where ` {
		q.Command = `select count(*) from "` + req.Table + `"` + where_clause + time_clause
	} else {
		q.Command = `select count(*) from "` + req.Table + `"` + where_clause + ` and ` + time_clause
	}
	glog.Errorln(q)

	body, err = this.log_forwarder.Query(q)
	if err != nil {
		glog.Errorln(err)
		return &ripple.Response{Status: 500, Body: err.Error(), ContentType: "text/plain"}
	}

	if log_rst, err := parseInfluxdbResult(body); err != nil {
		glog.Errorln(err)
		return &ripple.Response{Status: 500, Body: err.Error(), ContentType: "text/plain"}
	} else {
		return &ripple.Response{Status: 200, Body: string(log_rst)}
	}
}

func (this *LogController) GetDatabases(ctx *ripple.Context) *ripple.Response {
	q := new(models.Query)
	q.Command = `show databases`
	body, err := this.log_forwarder.Query(q)
	if err != nil {
		glog.Errorln(err)
		return &ripple.Response{Status: 500, Body: err.Error(), ContentType: "text/plain"}
	}

	if log_rst, err := parseInfluxdbResult(body); err != nil {
		glog.Errorln(err)
		return &ripple.Response{Status: 500, Body: err.Error(), ContentType: "text/plain"}
	} else {
		return &ripple.Response{Status: 200, Body: string(log_rst)}
	}
}

func (this *LogController) GetTables(ctx *ripple.Context) *ripple.Response {
	db := ctx.Params["db"]
	if db == "" {
		err := fmt.Errorf("Database name not found")
		glog.Errorln(err)
		return &ripple.Response{Status: 400, Body: err.Error(), ContentType: "text/plain"}
	}

	q := new(models.Query)
	q.Database = db
	q.Command = `show measurements`

	body, err := this.log_forwarder.Query(q)
	if err != nil {
		glog.Errorln(err)
		return &ripple.Response{Status: 500, Body: err.Error(), ContentType: "text/plain"}
	}

	if log_rst, err := parseInfluxdbResult(body); err != nil {
		glog.Errorln(err)
		return &ripple.Response{Status: 500, Body: err.Error(), ContentType: "text/plain"}
	} else {
		return &ripple.Response{Status: 200, Body: string(log_rst)}
	}
}

func parseInfluxdbResult(data []byte) ([]byte, error) {
	var rst InfluxdbResults
	err := json.Unmarshal(data, &rst)
	if err != nil {
		return nil, err
	}

	var log_rst []byte
	if len(rst.Results) > 0 && len(rst.Results[0].Series) > 0 {
		log_rst, err = json.Marshal(rst.Results[0].Series[0])
	} else {
		var tmp InfluxdbSerie
		log_rst, err = json.Marshal(tmp)
	}
	if err != nil {
		return nil, err
	}
	return log_rst, nil
}
