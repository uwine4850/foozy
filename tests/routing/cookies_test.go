package routing

import (
	"io"
	"net/http"
	"testing"
)

func TestLoginSession(t *testing.T) {
	createReq, err := http.NewRequest("GET", "http://localhost:8030/session-create", nil)
	if err != nil {
		t.Error(err)
	}
	createResp, err := http.DefaultClient.Do(createReq)
	if err != nil {
		t.Error(err)
	}
	defer createResp.Body.Close()

	readReq, err := http.NewRequest("GET", "http://localhost:8030/session-read", nil)
	if err != nil {
		t.Error(err)
	}

	readReq.AddCookie(createResp.Cookies()[0])

	readResp, err := http.DefaultClient.Do(readReq)
	if err != nil {
		t.Error(err)
	}
	defer readResp.Body.Close()

	body, err := io.ReadAll(readResp.Body)
	if err != nil {
		t.Error(err)
	}
	if string(body) != "111" {
		t.Errorf("The secure session data was not read correctly.")
	}
}
