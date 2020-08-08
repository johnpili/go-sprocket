package sprocket

import (
	"net/http"
	"net/url"

	"github.com/gorilla/sessions"
)

// CheckAuthenticatedSession ...
func CheckAuthenticatedSession(next http.HandlerFunc, cookieStore *sessions.CookieStore, cookieName string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		returnURL := r.RequestURI
		session, _ := cookieStore.Get(r, cookieName)
		if auth, ok := session.Values["authenticated"].(bool); !auth || !ok {
			if len(returnURL) <= 0 {
				http.Redirect(w, r, "/login", 302)
				return
			}
			http.Redirect(w, r, "/login?returnURL="+url.QueryEscape(returnURL), 302)
			return
		}
		next.ServeHTTP(w, r)
	})
}
