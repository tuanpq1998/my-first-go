package auth

import (
	"errors"
	"net/http"
	"strings"
)

/*
Extract apikey from headers request
Example of headers:
`Authorization: ApiKey {XXXX}`
return {XXXX}
*/
func ParseApiKey(headers http.Header) (string, error) {
	value := headers.Get("Authorization")
	if value == "" {
		return "", errors.New("not found authentication")
	}
	splited := strings.Split(value, " ")
	if len(splited) != 2 || splited[0] != "ApiKey" {
		return "", errors.New("malformed auth header")
	}
	return splited[1], nil
}
