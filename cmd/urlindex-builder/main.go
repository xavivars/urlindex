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

	flag.StringVar(&tlpReportFile, "tlp-report", "", "File with a TLP report with URLs")
	flag.StringVar(&tlpReportFile, "t", "", "File with a TLP report with URLs")
}

func main() {

	flag.Parse()

	if tlpReportFile == "" && listFile == "" {
		fmt.Errorf("either list-file or tlp-report need to be provided")
		return
	}

	if tlpReportFile != "" && listFile != "" {
		fmt.Errorf("only list-file or tlp-report can be provided")
		return
	}

	format, file := getFileAndFormat()

	text := getContent(file)

	log.Info("About to save")
	o := builder.Save(format, "vistaprint", "en-ie", text)

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

func getFileAndFormat() (string, string) {
	var format string
	var file string
	if tlpReportFile != "" {
		format = builder.TLP
		file = tlpReportFile
	} else {
		format = builder.List
		file = listFile
	}
	return format, file
}
