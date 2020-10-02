package mangadex

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"mangapdf/common"
	"mime/multipart"
	"net/http"
	"strings"
)

var exmaple = "/search?title=aposimz"
var cookieFile = "cookie.json"

//CheckLogin ...
func (m *Mangadex) CheckLogin() bool {
	//checkFile ...
	if !m.checkFile() {
		return false
	}
	return m.searchCheck()
}

func (m *Mangadex) searchCheck() bool {
	resp, err := m.fetch(req{
		method: "GET",
		link:   exmaple,
	})
	if err != nil {
		log.Println(err)
		return false
	}
	str, err := common.GetBodyString(resp)
	if err != nil {
		log.Println(err)
		return false
	}
	if strings.Contains(str, `<h1 class="text-center">Login</h1>`) {
		m.cookie = m.cookie[:0]
		return false
	}
	return true
}
func (m *Mangadex) checkFile() bool {
	buf, err := ioutil.ReadFile(cookieFile)
	if err != nil {
		log.Println("error found", err)
		return false
	}
	if len(buf) == 0 {
		return false
	}
	var cookie []*http.Cookie
	err = json.Unmarshal(buf, &cookie)
	if err != nil {
		log.Println("error found:", err)
		return false
	}
	m.cookie = cookie
	return true
}

//Login ...
func (m *Mangadex) Login(userName, password string) error {
	// data := url.Values{"login_username": {"0"}, "login_password": {"xxxx"}, "two_factor": {""}, "remember_me": {"1"}}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("login_username", userName)
	writer.WriteField("login_password", password)
	writer.WriteField("two_factor", "")
	writer.WriteField("remember_me", "1")
	resp, err := m.fetch(req{
		link: "/ajax/actions.ajax.php?function=login",
		body: body.Bytes(),
		header: map[string]string{
			"Content-Type": writer.FormDataContentType(),
		},
		method: "POST",
	})
	if err != nil {
		return err
	}
	str, err := common.GetBodyString(resp)
	log.Println(str)
	if err != nil && err.Error() != "zero length" {
		return err
	}
	if strings.Contains(str, "Incorrect username or password.") {
		return errors.New("incorrect username or password")
	}
	m.cookie = common.GetCookie(resp)
	buf, err := json.Marshal(m.cookie)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(cookieFile, buf, 0664)
}
