package gws

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/google/go-github/github"
)

type Fixture struct {
	Event     string          `json:"event"`
	Signature string          `json:"signature"`
	Payload   json.RawMessage `json:"payload"`
}

func TestServer(t *testing.T) {
	s := New("")
	ts := httptest.NewServer(s)
	client := &http.Client{}
	defer ts.Close()

	data, err := ioutil.ReadFile("fixtures/push.json")
	assert.Equal(t, nil, err)

	var fixture Fixture
	json.Unmarshal(data, &fixture)

	requestBody := bytes.NewReader(fixture.Payload)
	req, err := http.NewRequest("POST", ts.URL, requestBody)
	req.Header.Add("X-Github-Event", fixture.Event)
	assert.Equal(t, nil, err)

	var event *github.PushEvent
	go func() {
		event = <-s.PushEvents
	}()

	resp, err := client.Do(req)
	assert.Equal(t, nil, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotEqual(t, nil, event)
}
