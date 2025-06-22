package tutils

import (
	"fmt"
	"io"
	"net/http"
)

const (
	PortSimpleHandlers = ":7000"
	PortForm           = ":7001"
	PortObject         = ":7002"
	PortCookies        = ":7003"
)

func MakeUrl(port string, addres string) string {
	return fmt.Sprintf("http://localhost%s/%s", port, addres)
}

func ReadBody(body io.ReadCloser) (string, error) {
	data, err := io.ReadAll(body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func SendRequest(method string, url string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
