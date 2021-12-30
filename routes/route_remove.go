package routes

import (
	"forum/model"
	"forum/pkg"
	"net/http"
	"strconv"
)

//DeletePost deletes post from database
func DeletePost(w http.ResponseWriter, r *http.Request) {
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
		pkg.Danger(invMeth + "delete post")
		w.WriteHeader(http.StatusMethodNotAllowed)
		Forum.ErrStatus = http.StatusMethodNotAllowed
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	p := model.Post{
		AuthID: Forum.User.ID,
	}

	//getting post id
	postID := r.PostFormValue("PostId")

	p.ID, err = strconv.Atoi(postID)
	if err != nil {
		pkg.Danger(err, "Invalid post id - deleting post")
		Forum.ErrStatus = http.StatusBadRequest
		w.WriteHeader(http.StatusBadRequest)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	if Forum.User.Role == 1 || Forum.User.Role == 2 {
		//deleting post from database
		if err := p.DeletePostAMFromDataBase(); err != nil {
			pkg.Danger(err, "Cannot delete post - moderator or admin")
			Forum.ErrStatus = http.StatusInternalServerError
			w.WriteHeader(http.StatusInternalServerError)
			Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if Forum.User.Role == 3 {
		//deleting post from database
		if err := p.DeletePostFromDataBase(); err != nil {
			pkg.Danger(err, "Cannot delete post - user")
			Forum.ErrStatus = http.StatusInternalServerError
			w.WriteHeader(http.StatusInternalServerError)
			Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
			return
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

//DeleteComm deletes comment from database
func DeleteComm(w http.ResponseWriter, r *http.Request) {
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
		pkg.Danger(invMeth + "delete comment")
		w.WriteHeader(http.StatusMethodNotAllowed)
		Forum.ErrStatus = http.StatusMethodNotAllowed
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	c := model.Comment{
		AuthID: Forum.User.ID,
	}

	//getting comment id
	commID := r.PostFormValue("IdCom")

	c.ID, err = strconv.Atoi(commID)
	if err != nil {
		pkg.Danger(err, "Invalid comment id - deleting comment")
		Forum.ErrStatus = http.StatusBadRequest
		w.WriteHeader(http.StatusBadRequest)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	//getting post id
	postID := r.PostFormValue("PostId")

	c.PostID, err = strconv.Atoi(postID)
	if err != nil {
		pkg.Danger(err, "Invalid post id - deleting comment")
		Forum.ErrStatus = http.StatusBadRequest
		w.WriteHeader(http.StatusBadRequest)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	if Forum.User.Role == 1 {
		if err := c.DeleteCommentAdminFromDataBase(); err != nil {
			pkg.Danger(err, "Cannot delete comment - admin")
			Forum.ErrStatus = http.StatusInternalServerError
			w.WriteHeader(http.StatusInternalServerError)
			Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
			return
		}
		http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
		return
	}

	//deleting comment from database
	if err := c.DeleteCommentFromDataBase(Forum.User.ID); err != nil {
		pkg.Danger(err, "Cannot delete comment - user")
		Forum.ErrStatus = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}
