package common

import "net/http"

func GetParam(req *http.Request, param string) (string, bool) {
	keys, ok := req.URL.Query()[param]
	if !ok || len(keys[0]) < 1 {
		return "", ok
	}
	return keys[0], true
}
