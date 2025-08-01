package secure

import (
	"bytes"
	"net/http"
)

// Setting CSP.
// CSP (Content Security Policy) is a security mechanism implemented through HTTP headers
// or meta tags that helps prevent XSS (Cross-Site Scripting) attacks, injection,
// and other vulnerabilities associated with the introduction of malicious content.
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

// Setting AntiClickjacking.
// Anti-clickjacking refers to protective measures that prevent your website from
// being embedded in an <iframe> so that attackers cannot use clickjacking attacks.
func SetAntiClickjacking(w http.ResponseWriter, acOption string) {
	w.Header().Set("X-Frame-Options", acOption)
}
