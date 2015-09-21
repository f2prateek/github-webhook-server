package gws

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/google/go-github/github"
)

func TestInvalidMethod(t *testing.T) {
	_, ts := newTestServer("")
	defer ts.Close()

	resp, err := http.Get(ts.URL)
	check(err)

	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
}

func TestPushEvent(t *testing.T) {
	s, ts := newTestServer("")
	defer ts.Close()

	var event *github.PushEvent
	go func() {
		event = <-s.PushEvents
	}()

	resp := post(ts.URL, "push")

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotEqual(t, nil, event)
}

func TestPullRequestEvent(t *testing.T) {
	s, ts := newTestServer("")
	defer ts.Close()

	var event *github.PullRequestEvent
	go func() {
		event = <-s.PullRequestEvents
	}()

	resp := post(ts.URL, "pull_request")

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotEqual(t, nil, event)
}

func newTestServer(secret string) (*Server, *httptest.Server) {
	s := New(secret)
	ts := httptest.NewServer(s)
	return s, ts
}

func post(url, file string) *http.Response {
	data, err := ioutil.ReadFile("fixtures/" + file + ".json")
	check(err)

	var fixture struct {
		Event     *string         `json:"event"`
		Signature *string         `json:"signature"`
		Payload   json.RawMessage `json:"payload"`
	}
	err = json.Unmarshal(data, &fixture)
	check(err)

	requestBody := bytes.NewReader(fixture.Payload)
	req, err := http.NewRequest("POST", url, requestBody)
	check(err)

	if fixture.Event != nil {
		req.Header.Add("X-Github-Event", *fixture.Event)
	}
	if fixture.Signature != nil {
		req.Header.Add("X-Hub-Signature", *fixture.Signature)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	check(err)

	return resp
}

func check(err error) {
	if err != nil {
		log.Fatalf("error: %s", err)
	}
}
