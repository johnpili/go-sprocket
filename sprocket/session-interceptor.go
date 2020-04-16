package sprocket

import (
	"net/http"

	"github.com/gorilla/sessions"
)

//CheckAuthenticatedSession ...
func CheckAuthenticatedSession(next http.HandlerFunc, cookieStore *sessions.CookieStore, cookieName string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := cookieStore.Get(r, cookieName)
		if auth, ok := session.Values["authenticated"].(bool); !auth || !ok {
			http.Redirect(w, r, "/login", 302)
			return
		}
		next.ServeHTTP(w, r)
	})
}
