package mangadex

import "net/http"

//Mangadex ...
type Mangadex struct {
	cookie []*http.Cookie
	clt    http.Client
}

