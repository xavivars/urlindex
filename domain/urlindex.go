package domain

import (
	"errors"
	"fmt"
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
	lastUpdate      time.Time
	fst             *vellum.FST
	remotePath      string
	localPath       string
	tenant          string
	locale          string
	defaultResponse bool
	refreshRate     time.Duration
}

func NewUrlIndex(path string, t time.Duration, d bool) (*UrlIndex, error) {

	tempFilename, err := ioutil.TempFile("", "routes")

	if err != nil {
		return nil, err
	}

	u := UrlIndex{
		lastUpdate: time.Now().Add(-time.Hour * 2),
		remotePath: path,
		localPath:  tempFilename.Name(),
		defaultResponse: d,
		refreshRate: t,
	}

	u.Refresh()
	go refresh(&u)

	return &u, nil
}

func refresh(u *UrlIndex) {
	for range time.Tick(time.Minute){
		u.Refresh()
	}
}

func (u *UrlIndex) Refresh() bool {

	if u.lastUpdate.Add(time.Minute * u.refreshRate).Before(time.Now()) {

		u.downloadRemoteFile()

		fst, err := vellum.Open(u.localPath)
		if err == nil {
			log.Info(fmt.Sprintf("Updated fst: %s", u.remotePath))
			u.fst = fst
			u.lastUpdate = time.Now()
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

func (u *UrlIndex) downloadRemoteFile() error {

	data, err := u.getRemoteData()
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

func (u *UrlIndex) getRemoteData() (io.ReadCloser, error) {

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

