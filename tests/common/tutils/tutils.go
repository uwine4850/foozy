package tutils

import "fmt"

func MakeUrl(port string, addres string) string {
	return fmt.Sprintf("http://localhost%s/%s", port, addres)
}
