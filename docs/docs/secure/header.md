## Header
Ready-made modules for protection via Header.

#### SetCSP
Setting CSP.<br>
CSP (Content Security Policy) is a security mechanism implemented through HTTP headers or meta tags that helps prevent XSS (Cross-Site Scripting) attacks, injection, and other vulnerabilities associated with the introduction of malicious content.
```golang
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
```

#### SetAntiClickjacking
Setting AntiClickjacking.<br>
Anti-clickjacking refers to protective measures that prevent your website from
being embedded in an `<iframe>` so that attackers cannot use clickjacking attacks.
```golang
func SetAntiClickjacking(w http.ResponseWriter, acOption string) {
	w.Header().Set("X-Frame-Options", acOption)
}
```