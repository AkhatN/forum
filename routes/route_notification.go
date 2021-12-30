package routes

import (
	"forum/model"
	"forum/pkg"
	"net/http"
)

//Notification shows notification page for a registered user
func Notification(w http.ResponseWriter, r *http.Request) {
	Forum := model.Forum{}

	cookie, err := r.Cookie("session")
	if err != nil {
		pkg.Warning(err, "Failed to get cookie")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if err := Forum.User.ReadByUiid(cookie.Value); err != nil {
		model.DropCookie(w, r)
		pkg.Danger(err, usUUID)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method != http.MethodGet {
		pkg.Danger("Invalid method title")
		Forum.ErrStatus = http.StatusMethodNotAllowed
		w.WriteHeader(http.StatusMethodNotAllowed)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	if err := Forum.ReadLikedNotification(); err != nil {
		pkg.Danger(err, "Cannot get liked and disliked posts notifications")
		Forum.ErrStatus = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	if err := Forum.ReadCommentNotification(); err != nil {
		pkg.Danger(err, "Cannot get liked posts of registered user")
		Forum.ErrStatus = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	if err := Forum.UpdateNotifications(); err != nil {
		pkg.Danger(err, "Cannot delete notifications")
		Forum.ErrStatus = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	activity, err := pkg.NewView("views/notification.html")
	if err != nil {
		pkg.Danger(err, "Cannot execute home template of created posts")
		Forum.ErrStatus = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	activity.ExecuteTemplate(w, "bootstrap", Forum)
}
