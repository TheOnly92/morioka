package infrastructure

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type HerokuHandler struct {
	apiKey  string
	appName string
}

type HerokuAppRelease struct {
	Env       map[string]string `json:"env"`
	Commit    interface{}       `json:"commit"`
	User      string            `json:"user"`
	CreatedAt string            `json:"created_at"`
	Descr     string            `json:"descr"`
	Pstable   struct {
		Web string `json:"web"`
	} `json:"pstable"`
	Name   string   `json:"name"`
	Addons []string `json:"addons"`
}

func NewHerokuHandler(apiKey, appName string) *HerokuHandler {
	return &HerokuHandler{
		apiKey,
		appName,
	}
}

func (h *HerokuHandler) GetReleases() ([]*HerokuAppRelease, error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	req, err := http.NewRequest("GET", "https://:"+h.apiKey+"@api.heroku.com/apps/"+h.appName+"/releases", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var rt []*HerokuAppRelease
	err = json.Unmarshal(body, &rt)
	if err != nil {
		return nil, err
	}
	return rt, nil
}
