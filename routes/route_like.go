package routes

import (
	"forum/model"
	"forum/pkg"
	"net/http"
	"strconv"
)

//Like ставит лайк на пост для юзера и сохраняет в бд
func Like(w http.ResponseWriter, r *http.Request) {
	Forum := model.Forum{}

	if r.Method != http.MethodPost {
		pkg.Danger(invMeth + "save like on post")
		Forum.ErrStatus = http.StatusMethodNotAllowed
		w.WriteHeader(http.StatusMethodNotAllowed)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	var err error
	cookie := r.Context().Value("user").(string)
	user := model.User{}
	if err := user.ReadByUiid(cookie); err != nil {
		model.DropCookie(w, r)
		pkg.Danger(err, usUUID)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	saveLiked := model.SelectLike{}
	saveLiked.AuthID = user.ID
	saveLiked.PostID, err = strconv.Atoi(r.FormValue("IdPOST"))
	if err != nil {
		pkg.Danger(err, "Invalid IdPost")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if err := saveLiked.ReadPostAuthIDByPostID(); err != nil {
		pkg.Danger(err, "Invalid IdPost")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	switch r.PostFormValue("selector") {
	case "like":
		saveLiked.IsLike = true
		saveLiked.IsDislike = false
	case "clike":
		saveLiked.IsLike = false
		saveLiked.IsDislike = false
	case "dislike":
		saveLiked.IsLike = false
		saveLiked.IsDislike = true
	case "cdislike":
		saveLiked.IsLike = false
		saveLiked.IsDislike = false
	}

	if err = saveLiked.WriteLikeToDataBase(); err != nil {
		pkg.Danger(err, "Cannot write post like to database")
		http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}

//ComLike ставит лайк на коммент для юзера и сохраняет в бд
func ComLike(w http.ResponseWriter, r *http.Request) {
	Forum := model.Forum{}

	if r.Method != http.MethodPost {
		pkg.Danger(invMeth + "save like on comment")
		Forum.ErrStatus = http.StatusMethodNotAllowed
		w.WriteHeader(http.StatusMethodNotAllowed)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	var err error
	cookie := r.Context().Value("user").(string)
	user := model.User{}
	if err = user.ReadByUiid(cookie); err != nil {
		model.DropCookie(w, r)
		pkg.Danger(err, usUUID)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	saveLiked := model.SelectComLike{}
	saveLiked.AuthID = user.ID
	saveLiked.ComID, err = strconv.Atoi(r.FormValue("IdCom"))
	if err != nil {
		pkg.Danger(err, "Invalid IdPost")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	switch r.PostFormValue("selector") {
	case "like":
		saveLiked.IsLike = true
		saveLiked.IsDislike = false
	case "clike":
		saveLiked.IsLike = false
		saveLiked.IsDislike = false
	case "dislike":
		saveLiked.IsLike = false
		saveLiked.IsDislike = true
	case "cdislike":
		saveLiked.IsLike = false
		saveLiked.IsDislike = false
	}

	if err = saveLiked.WriteComLikeToDataBase(); err != nil {
		pkg.Danger(err, "Cannot write comment like to database")
		http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}
