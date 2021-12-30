package routes

import (
	"forum/model"
	"forum/pkg"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//Savecomment saves comment of a post
func Savecomment(w http.ResponseWriter, r *http.Request) {
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

	if r.Method != http.MethodPost {
		pkg.Danger(invMeth + "save comment")
		w.WriteHeader(http.StatusMethodNotAllowed)
		Forum.ErrStatus = http.StatusMethodNotAllowed
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	// Check if empty string is sent
	if err := pkg.EmptySpaceCheck(r.FormValue("comment")); err != nil {
		pkg.Danger(err, "Characters are not validate")
		w.WriteHeader(http.StatusBadRequest)
		Forum.ErrStatus = http.StatusBadRequest
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	c := model.Comment{
		Comm:   r.FormValue("comment"),
		AuthID: Forum.User.ID,
		Date:   time.Now().Format("2006.01.02 15:04:05"),
	}

	postID := strings.Split(r.Referer(), "/title/")

	if len(postID) < 2 {
		pkg.Danger(err, "Invalid referer - saving comment")
		w.WriteHeader(http.StatusBadRequest)
		Forum.ErrStatus = http.StatusBadRequest
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	c.PostID, err = strconv.Atoi(postID[1])
	if err != nil {
		pkg.Danger(err, "Invalid post id - saving comment")
		w.WriteHeader(http.StatusBadRequest)
		Forum.ErrStatus = http.StatusBadRequest
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	if err := c.ReadPostAuthIDByPostID(); err != nil {
		pkg.Danger(err, "Post doesn't exist")
		w.WriteHeader(http.StatusBadRequest)
		Forum.ErrStatus = http.StatusBadRequest
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	if err := c.WriteCommentToDataBase(); err != nil {
		pkg.Danger(err, "Cannot write comment")
		w.WriteHeader(http.StatusInternalServerError)
		Forum.ErrStatus = http.StatusInternalServerError
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}
