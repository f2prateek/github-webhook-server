package gws

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/google/go-github/github"
)

type Server struct {
	secret     string
	PushEvents chan *github.PushEvent
}

func New(secret string) *Server {
	return &Server{
		secret:     secret,
		PushEvents: make(chan *github.PushEvent, 1),
	}
}

// Satisfy the http.Handler interface
func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		httpError(w, req.Method, http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		httpError(w, "", http.StatusInternalServerError)
		return
	}
	defer req.Body.Close()

	if s.secret != "" {
		sig := req.Header.Get("X-Hub-Signature")
		if sig == "" {
			httpError(w, "Missing X-Hub-Signature", http.StatusForbidden)
			return
		}

		mac := hmac.New(sha1.New, []byte(s.secret))
		mac.Write(body)
		expectedMAC := mac.Sum(nil)
		expectedSig := "sha1=" + hex.EncodeToString(expectedMAC)
		if !hmac.Equal([]byte(expectedSig), []byte(sig)) {
			httpError(w, "HMAC verification failed", http.StatusForbidden)
			return
		}
	}

	eventType := req.Header.Get("X-GitHub-Event")
	if eventType == "" {
		httpError(w, "Missing X-GitHub-Event Header", http.StatusBadRequest)
		return
	}

	if eventType == "push" {
		var pushEvent *github.PushEvent
		err := json.Unmarshal(body, &pushEvent)
		if err != nil {
			httpError(w, "Could not decode body", http.StatusInternalServerError)
			return
		}

		select {
		case s.PushEvents <- pushEvent:
		default:
		}
	}
}
