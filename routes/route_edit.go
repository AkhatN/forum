package routes

import (
	"forum/model"
	"forum/pkg"
	"net/http"
	"os"
	"strconv"
	"time"
)

//EditPost edit post page
func EditPost(w http.ResponseWriter, r *http.Request) {
	Forum := model.Forum{}
	var err error

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

	if err := Forum.ReadCountNotification(); err != nil {
		pkg.Danger("Cannot get notifications")
		Forum.ErrStatus = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	if r.Method != http.MethodGet {
		pkg.Danger(invMeth + "edit post")
		Forum.ErrStatus = http.StatusMethodNotAllowed
		w.WriteHeader(http.StatusMethodNotAllowed)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	//берем id поста
	Forum.Post.ID, err = strconv.Atoi(r.FormValue("PostId"))
	if err != nil {
		pkg.Danger(err, "Not found: title - edit post")
		Forum.ErrStatus = http.StatusBadRequest
		w.WriteHeader(http.StatusBadRequest)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	userID := 0
	//находим post по id, есть ли вообще такой
	userID, err = Forum.Post.ReadPostForEdit(Forum.User.ID)
	if err != nil {
		pkg.Danger(err, "Not found in db: post")
		Forum.ErrStatus = http.StatusNotFound
		w.WriteHeader(http.StatusNotFound)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	//затем сравниваем post_id с предоставленным id поста
	if err := Forum.CheckPost(userID); err != nil {
		pkg.Danger(err, "Different UserID")
		Forum.ErrStatus = http.StatusBadRequest
		w.WriteHeader(http.StatusBadRequest)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	if _, err := Forum.Post.ReadImagefromDb(); err != nil {
		Forum.Post.IsImage = false
	}

	editPostView, err := pkg.NewView("views/editpost.html")
	if err != nil {
		pkg.Danger("Cannot execute edit post` template")
		Forum.ErrStatus = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	editPostView.ExecuteTemplate(w, "bootstrap", Forum)
}

//EditedPost ...
func EditedPost(w http.ResponseWriter, r *http.Request) {
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
		pkg.Danger(invMeth + "edit post")
		w.WriteHeader(http.StatusMethodNotAllowed)
		Forum.ErrStatus = http.StatusMethodNotAllowed
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	//parses url
	url := r.Referer()
	u, err := r.URL.Parse(url)
	if err != nil {
		pkg.Danger("Cannot parse url" + "edit post")
		w.WriteHeader(http.StatusInternalServerError)
		Forum.ErrStatus = http.StatusInternalServerError
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	//parsquery from url
	p := u.Query()

	//берем id поста
	postID := p.Get("PostId")
	if postID == "" {
		postID = r.FormValue("PostId")
	}

	postid, err := strconv.Atoi(postID)
	if err != nil {
		pkg.Danger("Invalid post id" + " edit post")
		w.WriteHeader(http.StatusBadRequest)
		Forum.ErrStatus = http.StatusBadRequest
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	oldForum := model.Forum{
		User: model.User{
			ID:       Forum.User.ID,
			Username: Forum.User.Username,
			IsStatus: Forum.User.IsStatus,
		},
		Post: model.Post{
			ID: postid,
		},
	}

	_, err = oldForum.Post.ReadPostForEdit(Forum.User.ID)
	if err != nil {
		pkg.Danger(err, "Invalid post for edit")
		oldForum.ErrStatus = http.StatusBadRequest
		w.WriteHeader(http.StatusBadRequest)
		Errtmpl.ExecuteTemplate(w, "bootstrap", oldForum)
		return
	}

	//checks if the post have an image
	oldForum.Post.Image, err = oldForum.Post.ReadImagefromDb()
	if err != nil {
		oldForum.Post.IsImage = false
	}

	// Check if empty string is sent
	if err := pkg.EmptySpaceCheck(r.FormValue("title"), r.FormValue("text")); err != nil {
		pkg.Danger(err, "Characters are not validate")
		oldForum.IsErr = true
		oldForum.ErrMsg = err.Error()
		postView, err := pkg.NewView("views/editpost.html")
		if err != nil {
			pkg.Danger("Cannot execute editpost template")
			oldForum.ErrStatus = http.StatusInternalServerError
			w.WriteHeader(http.StatusInternalServerError)
			Errtmpl.ExecuteTemplate(w, "bootstrap", oldForum)
			return
		}

		w.WriteHeader(400)
		postView.ExecuteTemplate(w, "bootstrap", oldForum)
		return
	}

	// Read data categories from client
	cs := r.MultipartForm.Value["addmore[]"]
	// Check if empty string is sent
	if err := pkg.EmptySpaceCheck(cs...); err != nil {
		pkg.Danger(err, "Characters are not validate")
		oldForum.IsErr = true
		oldForum.ErrMsg = err.Error()
		postView, err := pkg.NewView("views/editpost.html")
		if err != nil {
			pkg.Danger("Cannot execute editpost template")
			oldForum.ErrStatus = http.StatusInternalServerError
			w.WriteHeader(http.StatusInternalServerError)
			Errtmpl.ExecuteTemplate(w, "bootstrap", oldForum)
			return
		}

		w.WriteHeader(400)
		postView.ExecuteTemplate(w, "bootstrap", oldForum)
		return
	}

	// Writing data to struct
	Forum.Post = model.Post{
		ID:       postid,
		Title:    r.FormValue("title"),
		FullText: r.FormValue("text"),
		AuthID:   Forum.User.ID,
		Date:     time.Now().Format("2006.01.02 15:04:05"),
	}

	//saving an image to the local machine before saving anything to database
	if Forum.Post.Image, err = pkg.ImageUpload(w, r); err != nil {
		pkg.Danger(err, "Cannot save image in the local machine")
		oldForum.IsErr = true
		oldForum.ErrMsg = err.Error()
		postView, err := pkg.NewView("views/editpost.html")
		if err != nil {
			pkg.Danger("Cannot execute post template")
			oldForum.ErrStatus = http.StatusInternalServerError
			w.WriteHeader(http.StatusInternalServerError)
			Errtmpl.ExecuteTemplate(w, "bootstrap", oldForum)
			return
		}

		w.WriteHeader(400)
		postView.ExecuteTemplate(w, "bootstrap", oldForum)
		return
	}

	// Writing post to database
	if err := Forum.EditPostInDataBase(); err != nil {
		pkg.Danger(err, "Limited amount of characters")
		oldForum.IsErr = true
		oldForum.ErrMsg = "Limited amount of characters in title(40)"
		postView, err := pkg.NewView("views/editpost.html")
		if err != nil {
			pkg.Danger("Cannot execute post template")
			oldForum.ErrStatus = http.StatusInternalServerError
			w.WriteHeader(http.StatusInternalServerError)
			Errtmpl.ExecuteTemplate(w, "bootstrap", oldForum)
			return
		}
		w.WriteHeader(400)
		postView.ExecuteTemplate(w, "bootstrap", oldForum)
		return
	}

	cats := []model.Category{}

	//Adding uniq in categories ...
	pkg.CheckUniq(&cats, cs)

	if err := Forum.EditCats(cats); err != nil {
		pkg.Danger(err, "Cannot save categories in database")
		oldForum.IsErr = true
		oldForum.ErrMsg = "Internal server error. Please try to create the post later."
		postView, err := pkg.NewView("views/editpost.html")
		if err != nil {
			pkg.Danger("Cannot execute post template")
			oldForum.ErrStatus = http.StatusInternalServerError
			w.WriteHeader(http.StatusInternalServerError)
			Errtmpl.ExecuteTemplate(w, "bootstrap", oldForum)
			return
		}

		w.WriteHeader(500)
		postView.ExecuteTemplate(w, "bootstrap", oldForum)
		return
	}

	//delete image path from database
	toDeleteImage := r.FormValue("deleteimage")
	if toDeleteImage == "on" {
		toDelete := true
		image := ""
		image, err = Forum.Post.ReadImagefromDb()
		if err != nil {
			toDelete = false
		}

		if toDelete {
			err := os.Remove("." + image)
			if err != nil {
				pkg.Danger(err, "Cannot delete image path in database")
				oldForum.IsErr = true
				oldForum.ErrMsg = err.Error()
				postView, err := pkg.NewView("views/editpost.html")
				if err != nil {
					pkg.Danger("Cannot execute editpost template")
					oldForum.ErrStatus = http.StatusInternalServerError
					w.WriteHeader(http.StatusInternalServerError)
					Errtmpl.ExecuteTemplate(w, "bootstrap", oldForum)
					return
				}

				w.WriteHeader(400)
				postView.ExecuteTemplate(w, "bootstrap", oldForum)
				return
			}

			if err := Forum.Post.DeleteImageFromDb(); err != nil {
				pkg.Danger(err, "Cannot delete image path in database")
				oldForum.IsErr = true
				oldForum.ErrMsg = err.Error()
				postView, err := pkg.NewView("views/editpost.html")
				if err != nil {
					pkg.Danger("Cannot execute editpost template")
					oldForum.ErrStatus = http.StatusInternalServerError
					w.WriteHeader(http.StatusInternalServerError)
					Errtmpl.ExecuteTemplate(w, "bootstrap", oldForum)
					return
				}

				w.WriteHeader(400)
				postView.ExecuteTemplate(w, "bootstrap", oldForum)
				return
			}
		}
	}

	if Forum.Post.Image != "" {
		toDelete := true
		image := ""
		image, err = Forum.Post.ReadImagefromDb()
		if err != nil {
			toDelete = false
		}

		if toDelete {
			err := os.Remove("." + image)
			if err != nil {
				pkg.Danger(err, "Cannot delete image path in database")
				oldForum.IsErr = true
				oldForum.ErrMsg = err.Error()
				postView, err := pkg.NewView("views/editpost.html")
				if err != nil {
					pkg.Danger("Cannot execute editpost template")
					oldForum.ErrStatus = http.StatusInternalServerError
					w.WriteHeader(http.StatusInternalServerError)
					Errtmpl.ExecuteTemplate(w, "bootstrap", oldForum)
					return
				}

				w.WriteHeader(400)
				postView.ExecuteTemplate(w, "bootstrap", oldForum)
				return
			}
		}
		// Writing an image path to database
		if err := Forum.Post.WriteImagetoDb(); err != nil {
			pkg.Danger(err, "Cannot save image path in database")
			oldForum.IsErr = true
			oldForum.ErrMsg = "Internal server error. Please try to create the post later."
			postView, err := pkg.NewView("views/editpost.html")
			if err != nil {
				pkg.Danger("Cannot execute post template")
				oldForum.ErrStatus = http.StatusInternalServerError
				w.WriteHeader(http.StatusInternalServerError)
				Errtmpl.ExecuteTemplate(w, "bootstrap", oldForum)
				return
			}

			w.WriteHeader(500)
			postView.ExecuteTemplate(w, "bootstrap", oldForum)
			return
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

//EditComm ...
func EditComm(w http.ResponseWriter, r *http.Request) {
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

	if err := Forum.ReadCountNotification(); err != nil {
		pkg.Danger("Cannot get notifications")
		Forum.ErrStatus = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	if r.Method != http.MethodGet {
		pkg.Danger(invMeth + "create post")
		Forum.ErrStatus = http.StatusMethodNotAllowed
		w.WriteHeader(http.StatusMethodNotAllowed)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	//берем id поста
	PostID, err := strconv.Atoi(r.FormValue("PostId"))
	if err != nil {
		pkg.Danger(err, "Not found: title - edit comment")
		Forum.ErrStatus = http.StatusBadRequest
		w.WriteHeader(http.StatusBadRequest)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	//берем id комментария
	idcomm := r.FormValue("IdCom")
	Forum.Post.Com.ID, err = strconv.Atoi(idcomm)
	if err != nil {
		pkg.Danger(err, "Not found: comment - edit comment")
		Forum.ErrStatus = http.StatusBadRequest
		w.WriteHeader(http.StatusBadRequest)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	//находим комментарий по id, есть ли вообще такой
	if err := Forum.ReadCommentByID(); err != nil {
		pkg.Danger(err, "Not found in db: comment - edit comment")
		Forum.ErrStatus = http.StatusNotFound
		w.WriteHeader(http.StatusNotFound)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	//затем сравниваем post_id с предоставленным id поста
	if err := Forum.CheckComment(PostID); err != nil {
		pkg.Danger(err, "Different PostID")
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

	editcommView, err := pkg.NewView("views/editcomm.html")
	if err != nil {
		pkg.Danger("Cannot execute post template")
		Forum.ErrStatus = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	editcommView.ExecuteTemplate(w, "bootstrap", Forum)
}

//EditedComm edites a comment
func EditedComm(w http.ResponseWriter, r *http.Request) {
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
		pkg.Danger(invMeth + "edit comment")
		w.WriteHeader(http.StatusMethodNotAllowed)
		Forum.ErrStatus = http.StatusMethodNotAllowed
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	//parses url
	url := r.Referer()
	u, err := r.URL.Parse(url)
	if err != nil {
		pkg.Danger("Cannot parse url" + "edit comment")
		w.WriteHeader(http.StatusInternalServerError)
		Forum.ErrStatus = http.StatusInternalServerError
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	c := model.Comment{
		Comm:   r.FormValue("comment"),
		AuthID: Forum.User.ID,
	}

	//parsquery from url
	p := u.Query()

	//берем id поста
	postID := p.Get("PostId")
	if postID == "" {
		postID = r.FormValue("PostId")
	}

	c.PostID, err = strconv.Atoi(postID)
	if err != nil {
		pkg.Danger("Invalid post id" + "edit comment")
		w.WriteHeader(http.StatusBadRequest)
		Forum.ErrStatus = http.StatusBadRequest
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}
	Forum.Post.ID = c.PostID

	//берем id комментария
	idcomm := p.Get("IdCom")
	if idcomm == "" {
		idcomm = r.FormValue("IdCom")
	}

	c.ID, err = strconv.Atoi(idcomm)
	if err != nil {
		pkg.Danger("Invalid comment id" + "edit comment")
		w.WriteHeader(http.StatusBadRequest)
		Forum.ErrStatus = http.StatusBadRequest
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}
	Forum.Post.Com.ID = c.ID
	oldForum := model.Forum{
		User: model.User{
			ID:       Forum.User.ID,
			Username: Forum.User.Username,
			IsStatus: Forum.User.IsStatus,
		},
		Post: model.Post{
			ID:  c.PostID,
			Com: model.Comment{ID: c.ID},
		},
	}

	//находим комментарий по id, есть ли вообще такой
	if err := oldForum.ReadCommentByID(); err != nil {
		pkg.Danger(err, "Not found in db: comment - edit comment")
		oldForum.ErrStatus = http.StatusNotFound
		w.WriteHeader(http.StatusNotFound)
		Errtmpl.ExecuteTemplate(w, "bootstrap", oldForum)
		return
	}

	// Check if empty string is sent
	if err := pkg.EmptySpaceCheck(r.FormValue("comment")); err != nil {
		pkg.Danger(err, "Characters are not validate")
		oldForum.IsErr = true
		oldForum.ErrMsg = err.Error()
		postView, err := pkg.NewView("views/editcomm.html")
		if err != nil {
			pkg.Danger("Cannot execute editcomm template")
			oldForum.ErrStatus = http.StatusInternalServerError
			w.WriteHeader(http.StatusInternalServerError)
			Errtmpl.ExecuteTemplate(w, "bootstrap", oldForum)
			return
		}

		w.WriteHeader(400)
		postView.ExecuteTemplate(w, "bootstrap", oldForum)
		return
	}

	if err := c.EditCommentInDataBase(); err != nil {
		pkg.Danger("Cannot edit comment in database")
		w.WriteHeader(http.StatusInternalServerError)
		oldForum.ErrStatus = http.StatusInternalServerError
		Errtmpl.ExecuteTemplate(w, "bootstrap", oldForum)
		return
	}

	http.Redirect(w, r, "/title/"+postID, http.StatusSeeOther)
}
