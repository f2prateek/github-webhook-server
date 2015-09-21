package gws

import (
	"bytes"
	"encoding/json"
	"io"
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
	defer resp.Body.Close()

	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	body, err := ioutil.ReadAll(resp.Body)
	check(err)
	assert.Equal(t, "Method Not Allowed: GET\n", string(body))
}

func TestMissingEventType(t *testing.T) {
	_, ts := newTestServer("")
	defer ts.Close()

	resp, err := http.Post(ts.URL, "", nil)
	check(err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	body, err := ioutil.ReadAll(resp.Body)
	check(err)
	assert.Equal(t, "Bad Request: Missing X-GitHub-Event Header\n", string(body))
}

func TestNoBody(t *testing.T) {
	_, ts := newTestServer("")
	defer ts.Close()

	resp := post(ts.URL, "no_body")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	body, err := ioutil.ReadAll(resp.Body)
	check(err)
	assert.Equal(t, "Internal Server Error: Could not decode body\n", string(body))
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
		Event     *string          `json:"event"`
		Signature *string          `json:"signature"`
		Payload   *json.RawMessage `json:"payload"`
	}
	err = json.Unmarshal(data, &fixture)
	check(err)

	var requestBody io.Reader
	if fixture.Payload != nil {
		requestBody = bytes.NewReader(*fixture.Payload)
	}
	req, err := http.NewRequest("POST", url, requestBody)
	check(err)

	if fixture.Event != nil {
		req.Header.Add("X-Github-Event", *fixture.Event)
	}
	if fixture.Signature != nil {
		req.Header.Add("X-Hub-Signature", *fixture.Signature)
	}

	resp, err := http.DefaultClient.Do(req)
	check(err)

	return resp
}

func check(err error) {
	if err != nil {
		log.Fatalf("error: %s", err)
	}
}
