package domain

import (
	"errors"
	"github.com/couchbase/vellum"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type UrlIndex struct {
	lastUpdate time.Time
	fst        *vellum.FST
	remotePath string
	localPath  string
	tenant     string
	locale     string
	defaultResponse bool
}

func NewUrlIndex(tenant string, locale string, path string, m int, d bool) (*UrlIndex, error) {

	tempFilename, err := ioutil.TempFile("", "routes")

	if err != nil {
		return nil, err
	}

	u := UrlIndex{
		lastUpdate: time.Now().Add(-time.Hour * 2),
		tenant:     tenant,
		locale:     locale,
		remotePath: path,
		localPath:  tempFilename.Name(),
		defaultResponse: d,
	}

	u.Refresh()
	go refresh(&u, m)

	return &u, nil
}

func refresh(u *UrlIndex, m int) {
	for range time.Tick(time.Minute * time.Duration(m)){
		u.Refresh()
	}
}

func (u *UrlIndex) Refresh() bool {

	if u.lastUpdate.Add(time.Minute*2).Before(time.Now()) {

		u.DownloadRemoteFile()

		fst, err := vellum.Open(u.remotePath)
		if err == nil {
			log.Info("Updated fst")
			u.fst = fst
			return true
		}
		log.Error(err)
	}

	return false
}

func (u *UrlIndex) Exists(s string) bool {

	if u.fst == nil {
		log.Warn("FST is not loaded yet")
		return u.defaultResponse
	}

	_, exists, err := u.fst.Get([]byte(s))
	if err != nil {
		log.Error(err)
	}

	return exists
}

func isFileRemote(remotePath string) bool {
	return strings.HasPrefix(remotePath, "http")
}

func (u *UrlIndex) DownloadRemoteFile() error {

	data, err := u.GetRemoteData()
	if err != nil {
		return err
	}
	defer data.Close()
	out, err := os.Create(u.localPath)
	if err != nil {
		return err
	}

	if _, err = io.Copy(out, data); err != nil {
		out.Close()
		return err
	}

	return out.Close()
}

func (u *UrlIndex) GetRemoteData() (io.ReadCloser, error) {

	var data io.ReadCloser
	if isFileRemote(u.remotePath) {
		resp, err := http.Get(u.remotePath)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != 200 {
			return nil, errors.New("download file failed")
		}
		data = resp.Body
	} else {
		file, err := os.Open(u.remotePath)
		if err != nil {
			return nil, err
		}
		data = file
	}

	return data, nil
}

