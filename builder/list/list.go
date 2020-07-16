package list

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"sort"
)

func GetUrls(s string) []string {

	ul := UrlList{ Urls: make([]string, 0)}

	err := json.Unmarshal([]byte(s), &ul)

	if err != nil {
		log.Fatal(err)
	}

	u := ul.Urls

	sort.Strings(u)

	return u
}

type UrlList struct {
	Urls []string `json:"urls"`
}