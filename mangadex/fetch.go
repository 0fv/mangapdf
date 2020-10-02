package mangadex

import (
	"bytes"
	"net/http"
	"strings"
)

var baseURL = "https://mangadex.org"

type req struct {
	method    string
	link      string
	linkParam map[string]string
	header    map[string]string
	body      []byte
}

func (m *Mangadex) fetch(r req) (*http.Response, error) {
	reqURL := r.link
	if (!strings.HasPrefix(reqURL, "http://")) && (!strings.HasPrefix(reqURL, "https://")) {
		reqURL = baseURL + reqURL
	}
	request, err := http.NewRequest(r.method, reqURL, bytes.NewReader(r.body))
	if err != nil {
		return nil, err
	}
	qr := request.URL.Query()
	for k, v := range r.linkParam {
		qr.Add(k, v)
	}
	request.URL.RawQuery = qr.Encode()
	//setCookie
	for _, v := range m.cookie {
		request.AddCookie(v)
	}
	request.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.102 Safari/537.36")
	request.Header.Set("authority", "mangadex.org")
	request.Header.Set("x-requested-with", "XMLHttpRequest")
	//setHeader
	for k, v := range r.header {
		request.Header.Set(k, v)
	}
	m.clt = http.Client{
		Transport: &http.Transport{
			Proxy:             http.ProxyFromEnvironment,
			ForceAttemptHTTP2: false,
		},
	}
	return m.clt.Do(request)
}
