package tlp

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"sort"
	"strings"
)

type report map[string]Result

type Result struct {
	Result MpvDatas `json:"result"`
}

type MpvDatas map[string]MpvData

type MpvData struct {
	Index    Index   `json:"index"`
	Children []Facet `json:"children"`
}

type Facet struct {
	Facet       string  `json:"facet"`
	Name        string  `json:"name"`
	PathSegment string  `json:"pathSegment"`
	Children    []Facet `json:"children,omitempty"`
}

type Index struct {
	PathSegment string `json:"pathSegment"`
}

func GetUrls(l string, s string) []string {
	rep := make(report)
	urls := make(map[string]bool)

	err := json.Unmarshal([]byte(s), &rep)

	if err != nil {
		log.Fatal(err)
	}

	loc, res, err := getResults(rep)

	if err != nil || loc != l {
		log.Fatal("Could not read file")
	}

	for mpvId, mpv := range res.Result {
		b := fmt.Sprintf("/%s/%s", l, strings.ToLower(mpvId))
		u := fmt.Sprintf("%s/index.html", b)
		urls[u] = true

		addInnerUrls(urls, b, mpv.Children)
	}

	keys := make([]string, 0)
	for k := range urls {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, v := range keys {
		fmt.Println(v)
	}

	return keys
}

func getResults(r report) (string, *Result, error) {
	for k, v := range r {
		return k, &v, nil
	}

	return "", nil, errors.New("Can't properly parse")
}

func addInnerUrls(urls map[string]bool, b string, facets []Facet) {

	for _, f := range facets {
		ib:= fmt.Sprintf("%s%s", b, f.PathSegment)
		u := fmt.Sprintf("/%s/index.html", ib)

		urls[u] = true

		addInnerUrls(urls, ib, f.Children)
	}
}
