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
	tenant          string
	locale          string
	defaultResponse bool
	refreshRate     time.Duration
	localPath		string
}

func NewUrlIndex(path string, t time.Duration, d bool) (*UrlIndex, error) {

	u := UrlIndex{
		lastUpdate: time.Now().Add(-time.Hour * 2),
		remotePath: path,
		defaultResponse: d,
		refreshRate: t,
	}

	_, err := u.Refresh()
	if err != nil {
		return nil, err
	}

	go refresh(&u)

	return &u, nil
}

func refresh(u *UrlIndex) {
	for range time.Tick(time.Minute){
		_, _ = u.Refresh()
	}
}

func (u *UrlIndex) Refresh() (bool, error) {

	if u.lastUpdate.Before(time.Now().Add(-u.refreshRate)) {
		log.Info(fmt.Sprintf("Updating fst: %s", u.remotePath))
		p, err := u.downloadRemoteFile()

		if err != nil {
			return false, err
		}

		fst, err := vellum.Open(p)
		if err == nil {

			oldFst := u.fst
			localPath := u.localPath

			log.Info(fmt.Sprintf("Updated fst: %s", u.remotePath))
			u.fst = fst
			u.localPath = p
			u.lastUpdate = time.Now()

			if oldFst != nil {
				oldFst.Close()
				os.Remove(localPath)
				log.Info(fmt.Sprintf("Temp file %s deleted", localPath))
			}

			return true, nil
		}
		log.Error(err)
	}

	return false, nil
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

func (u *UrlIndex) downloadRemoteFile() (string, error) {
	tempFile, err := ioutil.TempFile("", "routes")
	fileName := tempFile.Name()

	data, err := u.getRemoteData()
	if err != nil {
		return "", err
	}
	defer data.Close()

	if _, err = io.Copy(tempFile, data); err != nil {
		tempFile.Close()
		return "", err
	}

	return fileName, tempFile.Close()
}

func (u *UrlIndex) getRemoteData() (io.ReadCloser, error) {

	var data io.ReadCloser
	if isFileRemote(u.remotePath) {
		resp, err := get(u.remotePath)
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

func get(url string) (*http.Response, error) {
	client := &http.Client{ }

	req, _ := http.NewRequest("GET", url, nil)

	h := fmt.Sprintf("X-%s", time.Now().Format(time.RFC3339))

	// Try to avoid caches by varying headers
	req.Header.Add("Origin", h)
	req.Header.Add("Access-Control-Request-Headers", h)

	return client.Do(req)
}