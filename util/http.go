package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

//GetPostParams Get the parameters sent by the post method in an http request
func GetPostParams(r *http.Request) url.Values {
	switch {
	case strings.Contains(r.Header.Get("Content-Type"), "application/json"):
		params := map[string]interface{}{}
		result := url.Values{}
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &params)
		for k, v := range params {
			if reflect.ValueOf(v).Kind().String() == "string" {
				result.Set(k, v.(string))
			} else {
				result.Set(k, fmt.Sprint(v))
			}
		}
		return result
	case strings.Contains(r.Header.Get("Content-Type"), "application/x-www-form-urlencoded"):
		r.ParseForm()
		return r.Form
	case strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data"):
		r.ParseMultipartForm(int64(10 * 1000))
		return r.Form
	}
	return url.Values{}
}
