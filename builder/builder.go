package builder

import (
	"github.com/couchbase/vellum"
	log "github.com/sirupsen/logrus"
	"github.com/xavivars/urlindex/builder/list"
	"io/ioutil"
)

func Save(s string) string {

	urls := list.GetUrls(s)
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
