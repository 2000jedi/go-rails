package url

import (
	"net/http"
	"regexp"
)

type URL map[string]func(r http.Request) string

var Error map[int]func(r http.Request) string

func routing(url URL, r http.Request) string{
	for k, v := range url{
		match, _ := regexp.Match(k, []byte(url))
		if match {
			return v(r)
		}
	}
	return Error[404](r)
}
