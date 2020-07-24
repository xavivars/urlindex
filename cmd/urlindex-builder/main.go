package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"github.com/xavivars/urlindex/builder"
	"github.com/xavivars/urlindex/builder/list"
	"io/ioutil"
)

var listFile string
var tlpReportFile string

func init() {
	log.SetFormatter(&log.JSONFormatter{})

	flag.StringVar(&listFile, "list-file", "", "File with a list of URLs")
	flag.StringVar(&listFile, "l", "", "File with a list of URLs")
}

func main() {

	flag.Parse()

	if listFile == "" {
		log.Error("list-file is required")
		return
	}

	text := getContent(listFile)

	urls := list.GetUrls(text)

	log.Info("About to save")
	o := builder.GetFst(urls)

	log.Info(o)
}

func getContent(file string) string {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	text := string(content)
	return text
}