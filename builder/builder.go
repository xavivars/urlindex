package builder

import (
	"github.com/couchbase/vellum"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

func GetFst(urls []string) string {
	f, err := ioutil.TempFile("", "fst-")
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
