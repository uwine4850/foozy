package tutils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	PortSimpleHandlers = ":7000"
	PortForm           = ":7001"
	PortObject         = ":7002"
	PortCookies        = ":7003"
	PortSocket         = ":7004"
	PortMiddlewares    = ":7005"
	PortPageRender     = ":7006"
	PortCSRFToken      = ":7007"
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

func FilesAreEqual(path1, path2 string) (bool, error) {
	file1, err := os.Open(path1)
	if err != nil {
		return false, err
	}
	defer file1.Close()

	file2, err := os.Open(path2)
	if err != nil {
		return false, err
	}
	defer file2.Close()

	const chunkSize = 4096
	buf1 := make([]byte, chunkSize)
	buf2 := make([]byte, chunkSize)

	for {
		n1, err1 := file1.Read(buf1)
		n2, err2 := file2.Read(buf2)

		if n1 != n2 || !bytes.Equal(buf1[:n1], buf2[:n2]) {
			return false, nil
		}
		if err1 == io.EOF && err2 == io.EOF {
			break
		}
		if err1 != nil && err1 != io.EOF {
			return false, err1
		}
		if err2 != nil && err2 != io.EOF {
			return false, err2
		}
	}

	return true, nil
}
