package builder

import (
	"github.com/couchbase/vellum"
	log "github.com/sirupsen/logrus"
	"github.com/xavivars/urlindex/builder/list"
	"github.com/xavivars/urlindex/builder/tlp"
	"io/ioutil"
)

const (
	TLP = "TLP"
	List = "List"
)

func Save(ft string, t string, l string, s string) string {

	var urls []string

	switch ft {
		case TLP: urls = tlp.GetUrls(l, s)
		case List: urls = list.GetUrls(l, s)
	default:
		log.Fatal("Invalid format: %s", ft)
	}

	filename := getFst(urls)

	return filename
}

func getFst(urls []string) string {
	f, err := ioutil.TempFile("", "example")
	if err != nil {
		log.Fatal(err)
	}

	b, err := vellum.New(f, nil)
	if err != nil {
		log.Fatal(err)
	}

	for _, u := range urls {
		err = b.Insert([]byte(u), 1)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = b.Close()
	if err != nil {
		log.Fatal(err)
	}

	return f.Name()
}
