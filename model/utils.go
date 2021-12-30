package model

import (
	"net/http"
)

//SetCookie sets session cookie for a registered user
func SetCookie(w http.ResponseWriter) string {
	uuid := UUID()
	cookie := &http.Cookie{
		Name:     "session",
		Value:    uuid,
		MaxAge:   3600,
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)
	return uuid
}

//DropCookie deletes cookie from the client's browser
func DropCookie(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie("session")
	if err != nil {
		return err
	}
	cookie.MaxAge = -1
	http.SetCookie(w, cookie)
	return nil
}
