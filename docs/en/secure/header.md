## package secure

__SetCSP__
```go
SetCSP(w http.ResponseWriter, directives map[string][]string) error
```
Sets the “Content-Security-Policy” parmeter in the header. Passes directives with the directives map.
* Key — the name of the directive.
* Value — a slice of directive values.

__SetAntiClickjacking__
```go
SetAntiClickjacking(w http.ResponseWriter, acOption string)
```
Sets the “X-Frame-Options” parmeter in the header. Two values are available for this parameter:
* ACSameorigin(SAMEORIGIN)
* ACDeny(DENY)