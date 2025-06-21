package session

import (
	"net/http"
	"sync"
)

func NewSession(fname string) Session {
	return &sessionImpl{
		Filename: fname,
		Ua:       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36",
	}
}

type Session interface {
	GetCookies() []*http.Cookie
	UserAgent() string
	AddCookiesToRequest(req *http.Request)
	Update([]*http.Cookie) error
	FindCookie(string) string
}

type sessionImpl struct {
	Filename string `json:"filename"`
	Ua       string
	Cookies  []*http.Cookie
	sync.Mutex
}

// AddCookiesToRequest implements Session.
func (b *sessionImpl) AddCookiesToRequest(req *http.Request) {
	for _, cookie := range b.Cookies {
		req.AddCookie(cookie)
	}
}

// FindCookie implements Session.
func (b *sessionImpl) FindCookie(key string) string {
	for _, cookie := range b.Cookies {
		if cookie.Name == key {
			return cookie.Value
		}
	}

	return ""
}

// GetCookies implements Session.
func (b *sessionImpl) GetCookies() []*http.Cookie {
	return b.Cookies
}

// Update implements Session.
func (b *sessionImpl) Update(cookies []*http.Cookie) error {
	b.Lock()
	defer b.Unlock()

	for _, cookie := range cookies {
		err := b.updateCookie(cookie)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *sessionImpl) updateCookie(cookie *http.Cookie) error {

	fixCookies := []*http.Cookie{}

	isExist := false
	for _, oldCookie := range b.Cookies {
		if oldCookie.Name != cookie.Name {
			fixCookies = append(fixCookies, oldCookie)
			continue
		}

		fixCookies = append(fixCookies, cookie)
		isExist = true
	}

	if !isExist {
		fixCookies = append(fixCookies, cookie)
	}

	b.Cookies = fixCookies

	return nil
}

// UserAgent implements Session.
func (b *sessionImpl) UserAgent() string {
	return b.Ua
}
