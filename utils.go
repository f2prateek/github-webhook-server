package gws

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
)

// Print an `err` message with the given `code`.
func httpError(w http.ResponseWriter, err string, code int) {
	msg := http.StatusText(code)
	if err != "" {
		msg = msg + ": " + err
	}
	http.Error(w, msg, code)
}

// CheckMAC reports whether messageMAC is a valid HMAC tag for message.
func checkMAC(message, messageMAC, key []byte) bool {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	expectedSig := "sha1=" + hex.EncodeToString(expectedMAC)
	return hmac.Equal(messageMAC, []byte(expectedSig))
}
