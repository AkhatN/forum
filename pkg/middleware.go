package pkg

import (
	"context"
	"fmt"
	"forum/model"
	"net/http"
	"text/template"
	"time"
)

var templ []string = []string{
	"views/bootstrap.html",
	"views/navbar.html",
	"views/error.html",
}

//Middleware ...
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		cookie, err := r.Cookie("session")
		if err != nil {
			Warning(err, "Failed to get cookie")
			if r.Method == http.MethodPost || r.URL.Path == "/createpost" {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			Forum := model.Forum{}

			errtmpl, err := template.ParseFiles(templ...)
			if err != nil {
				Forum.ErrStatus = http.StatusNotFound
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, "Page not Found")
				return
			}

			Forum.ErrStatus = http.StatusNotFound
			w.WriteHeader(http.StatusNotFound)
			errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "user", cookie.Value)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
		fmt.Println("Second: ", time.Since(start).Seconds())
	})
}
