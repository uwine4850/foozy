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

func TestSetStandartCookie(t *testing.T) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "http://localhost:8030/cookie", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %v, but got %v", http.StatusOK, resp.StatusCode)
	}

	cookies := resp.Cookies()
	if len(cookies) == 0 {
		t.Fatal("Expected cookies to be set, but none were found")
	}

	var mycookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "cookie" {
			mycookie = cookie
			break
		}
	}

	if mycookie == nil {
		t.Fatal("Expected mycookie to be set, but it was not found")
	}

	if mycookie.Value != "value" {
		t.Errorf("Expected mycookie to have value %v, but got %v", "value", mycookie.Value)
	}
}
