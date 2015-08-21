package client

import (
	"testing"
)

func TestPostFile(t *testing.T) {
	urls := []string{
		"http://localhost:1734/upload",
		"http://localhost:1735/upload",
	}
	for _, url := range urls {
		filenames := []string{
			".test.1734",
			".test.1738",
			".test.1737",
			"logmanager_client.html",
		}
		for _, filename := range filenames {
			PostFile(filename, url)
		}
	}
}

func TestSendLog(t *testing.T) {
	urls := []string{
		"http://localhost:1734/upload",
		"http://localhost:1735/upload",
	}
	for _, url := range urls {
		SendLog(url, "./")
	}
}
