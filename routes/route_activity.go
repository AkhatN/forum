package routes

import (
	"forum/model"
	"forum/pkg"
	"net/http"
)

//Activity shows activity page
func Activity(w http.ResponseWriter, r *http.Request) {
	Forum := model.Forum{}

	cookie, err := r.Cookie("session")
	if err != nil {
		pkg.Danger(err, "Cannot get liked posts of registered user")
		Forum.ErrStatus = http.StatusNotFound
		w.WriteHeader(http.StatusNotFound)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	if err := Forum.User.ReadByUiid(cookie.Value); err != nil {
		model.DropCookie(w, r)
		pkg.Danger(err, "Cannot get liked posts of registered user")
		Forum.ErrStatus = http.StatusNotFound
		w.WriteHeader(http.StatusNotFound)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	if err := Forum.ReadCountNotification(); err != nil {
		pkg.Danger("Cannot get notifications")
		Forum.ErrStatus = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	if r.Method != http.MethodGet {
		pkg.Danger("Invalid method title")
		Forum.ErrStatus = http.StatusMethodNotAllowed
		w.WriteHeader(http.StatusMethodNotAllowed)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	if err := Forum.ReadLikedDislikedmyPostsFromDataBase(w, r); err != nil {
		pkg.Danger(err, "Cannot get liked posts of registered user")
		Forum.ErrStatus = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	activity, err := pkg.NewView("views/activity.html")
	if err != nil {
		pkg.Danger(err, "Cannot execute home template of created posts")
		Forum.ErrStatus = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	activity.ExecuteTemplate(w, "bootstrap", Forum)
}
