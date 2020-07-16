package main

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/xavivars/urlindex/builder"
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
		fmt.Errorf("list-file is required")
		return
	}

	text := getContent(file)

	log.Info("About to save")
	o := builder.Save(text)

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