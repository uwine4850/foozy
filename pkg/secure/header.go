package secure

import (
	"bytes"
	"net/http"
)

func SetCSP(w http.ResponseWriter, directives map[string][]string) error {
	var directivesString bytes.Buffer
	for name, value := range directives {
		_, err := directivesString.WriteString(name + " ")
		if err != nil {
			return err
		}
		for i := 0; i < len(value); i++ {
			if i == len(value)-1 {
				_, err := directivesString.WriteString(value[i] + "; ")
				if err != nil {
					return err
				}
			} else {
				_, err := directivesString.WriteString(value[i] + " ")
				if err != nil {
					return err
				}
			}
		}
	}
	w.Header().Set("Content-Security-Policy", directivesString.String())
	return nil
}

const (
	ACSameorigin = "SAMEORIGIN"
	ACDeny       = "DENY"
)

func SetAntiClickjacking(w http.ResponseWriter, acOption string) {
	w.Header().Set("X-Frame-Options", acOption)
}
