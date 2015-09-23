package main

import (
	//"bytes"
	// "bufio"
	"fmt"
	// "os"
	// "io/ioutil"
	"net/http"
	// "net/url"
	"strconv"
	"strings"
	"time"
	// "runtime"
)

func main() {
	i := 0
	for {
		// data := url.Values{}
		str := "cpc,type=ERROR host=2,region=1 " + strconv.FormatInt(time.Now().UnixNano(), 10)
		// s := strconv.FormatInt(time.Now().UnixNano(), 10)

		// fd, _ := os.OpenFile("LOG.x."+time.Now().Round(time.Millisecond*5).Format("MST2006-01-02T15:04:05Z"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		resp, _ := http.Post("http://localhost:8086/write?db=x", "", strings.NewReader(str))
		// str := "cpc,type=INFO,sensor=\"aa bb cc\",kaixin=12345,zero=\"ads;nf qefnqnfqewfn oqnef anfqwef\" host=2,region=1,message=\"haha\" " + strconv.FormatInt(time.Now().UnixNano(), 10)
		// writer := bufio.NewWriter(fd)
		// writer.WriteString(str)
		// writer.WriteByte('\n')
		// writer.Flush()
		fmt.Println(resp.StatusCode, i)
		i++
		// fd.Close()
	}
}
