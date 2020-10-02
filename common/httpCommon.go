package common

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

//GetCookie ...
func GetCookie(resp *http.Response) []*http.Cookie {
	return resp.Cookies()
}

//GetBody ...
func GetBody(resp *http.Response, i interface{}) error {
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(buf, i)
}

//GetBodyString ...
func GetBodyString(resp *http.Response) (string, error) {
	buf, err := ioutil.ReadAll(resp.Body)
	if len(buf) == 0 {
		return "", errors.New("zero length")
	}
	return string(buf), err
}

//GetDoc ...
func GetDoc(resp *http.Response) (*goquery.Document, error) {
	return goquery.NewDocumentFromResponse(resp)
}
