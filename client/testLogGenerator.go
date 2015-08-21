package main

import (
	//"bytes"
	//"fmt"
	"bufio"
	"os"
	"strconv"
	"time"
)

func main() {
	for {
		//fd, _ := os.OpenFile("test."+strconv.FormatInt(time.Now().Round(time.Minute*1).UnixNano(), 10), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		fd, _ := os.OpenFile("LOG.x."+time.Now().Round(time.Minute*1).Format("MST2006-01-02-15:04:05"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		str := "cpc,type=ERROR,kaixin=767845,me=\"ha\\ ha\" host=2,region=1,message=\"haha\" " + strconv.FormatInt(time.Now().UnixNano(), 10)
		// str := "cpc,type=INFO,sensor=\"aa bb cc\",kaixin=12345,zero=\"ads;nf qefnqnfqewfn oqnef anfqwef\" host=2,region=1,message=\"haha\" " + strconv.FormatInt(time.Now().UnixNano(), 10)
		writer := bufio.NewWriter(fd)
		writer.WriteString(str)
		writer.WriteByte('\n')
		writer.Flush()
		fd.Close()
		time.Sleep(time.Millisecond * 100)
	}
}
