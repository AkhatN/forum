package routes

import (
	"forum/model"
	"forum/pkg"
	"net/http"
	"strconv"
)

//Request shows request page for a moderator or an admin
func Request(w http.ResponseWriter, r *http.Request) {
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
		pkg.Danger("Invalid method title")
		Forum.ErrStatus = http.StatusMethodNotAllowed
		w.WriteHeader(http.StatusMethodNotAllowed)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	switch Forum.User.Role {
	case 1:
		if err := Forum.ReadRequestsForAdminFromDb(); err != nil {
			pkg.Danger("Cannot get requests for admin from database")
			Forum.ErrStatus = http.StatusInternalServerError
			w.WriteHeader(http.StatusInternalServerError)
			Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
			return
		}
	case 2:
		if err := Forum.ReadRequestsForModeratorFromDB(); err != nil {
			pkg.Danger("Cannot get requests for moderator from database")
			Forum.ErrStatus = http.StatusInternalServerError
			w.WriteHeader(http.StatusInternalServerError)
			Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
			return
		}
	case 3:
		if err := Forum.ReadRequestsForUserFromDB(); err != nil {
			pkg.Danger("Cannot get requests for user from database")
			Forum.ErrStatus = http.StatusInternalServerError
			w.WriteHeader(http.StatusInternalServerError)
			Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
			return
		}
	}

	requests, err := pkg.NewView("views/request.html")
	if err != nil {
		pkg.Danger(err, "Cannot execute request template of request posts")
		Forum.ErrStatus = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	requests.ExecuteTemplate(w, "bootstrap", Forum)
}

// HandleRequest shows a page of specific request
func HandleRequest(w http.ResponseWriter, r *http.Request) {
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
		pkg.Danger("Invalid method title")
		Forum.ErrStatus = http.StatusMethodNotAllowed
		w.WriteHeader(http.StatusMethodNotAllowed)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	Forum.Post.ID, err = strconv.Atoi(r.RequestURI[15:])
	if err != nil {
		pkg.Danger(err, "Invalid postid - handlerequest")
		Forum.ErrStatus = http.StatusBadRequest
		w.WriteHeader(http.StatusBadRequest)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	if Forum.User.Role != 1 && Forum.User.Role != 2 {
		pkg.Danger(err, "Cannot get specific request for user - handlerequest")
		Forum.ErrStatus = http.StatusNotFound
		w.WriteHeader(http.StatusNotFound)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	switch Forum.User.Role {
	case 1:
		if err := Forum.ReadRequestForAdmin(w, r); err != nil {
			pkg.Danger("Cannot get specific request for admin - handlerequest")
			Forum.ErrStatus = http.StatusInternalServerError
			w.WriteHeader(http.StatusInternalServerError)
			Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
			return
		}

	case 2:
		if err := Forum.Post.ReadRequestFromDb(); err != nil {
			pkg.Danger(err, "Cannot get specific request for moderator - handlerequest")
			Forum.ErrStatus = http.StatusInternalServerError
			w.WriteHeader(http.StatusInternalServerError)
			Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
			return
		}
	}

	if Forum.Post.Title == "" {
		pkg.Danger(err, "Cannot get handlerequest with invalid isvalid status")
		Forum.ErrStatus = http.StatusNotFound
		w.WriteHeader(http.StatusNotFound)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	if Forum.Post.Valid == true {
		pkg.Danger(err, "Cannot get specific request for moderator because it is handled already - handlerequest")
		Forum.ErrStatus = http.StatusNotFound
		w.WriteHeader(http.StatusNotFound)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	Forum.Post.Image, err = Forum.Post.ReadImagefromDb()
	if err != nil {
		Forum.Post.IsImage = false
	}

	requestView, err := pkg.NewView("views/handlerequest.html")
	if err != nil {
		pkg.Danger(err, "Cannot execute handlerequest template")
		Forum.ErrStatus = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	requestView.ExecuteTemplate(w, "bootstrap", Forum)
}

//DeleteRequest deletes post from database
func DeleteRequest(w http.ResponseWriter, r *http.Request) {
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

	if Forum.User.Role != 2 && Forum.User.Role != 1 {
		pkg.Danger(err, "The user is not moderator or admin for deleting a request")
		Forum.ErrStatus = http.StatusBadRequest
		w.WriteHeader(http.StatusBadRequest)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	p := model.Post{}

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

	if err := p.ReadAuthorID(); err != nil {
		pkg.Danger("Cannot find userid by postid from db")
		w.WriteHeader(http.StatusBadRequest)
		Forum.ErrStatus = http.StatusBadRequest
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	//deleting post from database
	if err := p.DeleteRequestFromDataBase(Forum.User.ID); err != nil {
		pkg.Danger(err, "Cannot delete post")
		Forum.ErrStatus = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

//Saverequest saves a new post
func Saverequest(w http.ResponseWriter, r *http.Request) {
	Forum := model.Forum{}

	// Getting user's ID to hand it to article auth_id
	cookie, err := r.Cookie("session")
	if err != nil {
		pkg.Warning(err, "Failed to get cookie")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if err := Forum.User.ReadByUiid(cookie.Value); err != nil {
		pkg.Danger(err, usUUID)
		model.DropCookie(w, r)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method != http.MethodPost {
		pkg.Danger(invMeth + "save post")
		w.WriteHeader(http.StatusMethodNotAllowed)
		Forum.ErrStatus = http.StatusMethodNotAllowed
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	if Forum.User.Role != 2 && Forum.User.Role != 1 {
		pkg.Danger(err, "The user is not moderator or admin for saving a request")
		Forum.ErrStatus = http.StatusBadRequest
		w.WriteHeader(http.StatusBadRequest)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	Forum.Post.ID, err = strconv.Atoi(r.FormValue("PostId"))
	if err != nil {
		pkg.Danger("Cannot update post_valid to 1")
		w.WriteHeader(http.StatusBadRequest)
		Forum.ErrStatus = http.StatusBadRequest
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	if err := Forum.Post.ReadAuthorID(); err != nil {
		pkg.Danger("Cannot find userid by postid from db")
		w.WriteHeader(http.StatusBadRequest)
		Forum.ErrStatus = http.StatusBadRequest
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	// Writing post to database
	if err := Forum.UpdatePostValid(); err != nil {
		pkg.Danger("Cannot update post to valid template")
		Forum.ErrStatus = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

//Report writes report about post for admin into database
func Report(w http.ResponseWriter, r *http.Request) {
	Forum := model.Forum{}
	var err error

	cookie, err := r.Cookie("session")
	if err != nil {
		pkg.Warning(err, "Failed to get cookie")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if err := Forum.User.ReadByUiid(cookie.Value); err != nil {
		pkg.Danger(err, usUUID)
		model.DropCookie(w, r)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method != http.MethodPost {
		pkg.Danger("Invalid method report")
		Forum.ErrStatus = http.StatusMethodNotAllowed
		w.WriteHeader(http.StatusMethodNotAllowed)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	Forum.Post.ID, err = strconv.Atoi(r.PostFormValue("PostId"))
	if err != nil {
		pkg.Danger(err, "Not found postid - report")
		Forum.ErrStatus = http.StatusNotFound
		w.WriteHeader(http.StatusNotFound)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	if Forum.User.Role != 2 {
		pkg.Danger(err, "The user is not moderator for reporting a post")
		Forum.ErrStatus = http.StatusBadRequest
		w.WriteHeader(http.StatusBadRequest)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	Forum.Post.FullText = r.PostFormValue("report")
	err = Forum.WriteReportToDb()
	if err != nil {
		pkg.Danger(err, "Cannot write report to db")
		Forum.ErrStatus = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	http.Redirect(w, r, "/request", http.StatusSeeOther)
}

//Makerequest ...
func Makerequest(w http.ResponseWriter, r *http.Request) {
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

	if r.Method == http.MethodGet {
		if Forum.User.Role == 1 {
			if err := Forum.ReadUsersFromDB(); err != nil {
				pkg.Danger(err)
				Forum.ErrStatus = http.StatusInternalServerError
				w.WriteHeader(http.StatusInternalServerError)
				Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
				return
			}
		}

		requestView, err := pkg.NewView("views/makerequest.html")
		if err != nil {
			pkg.Danger(err, "Cannot execute handlerequest template")
			Forum.ErrStatus = http.StatusInternalServerError
			w.WriteHeader(http.StatusInternalServerError)
			Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
			return
		}

		requestView.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	//Method POST
	if r.Method == http.MethodPost {
		switch Forum.User.Role {
		case 1:
			newuser := model.User{
				Username: r.FormValue("makereq"),
			}

			role := r.FormValue("role")
			if role != "user" && role != "moderator" {
				pkg.Danger("Invalid role - makerequest - admin")
				Forum.ErrStatus = http.StatusBadRequest
				w.WriteHeader(http.StatusBadRequest)
				Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
				return
			}

			if role == "user" {
				newuser.NewRole = 3
			}

			if role == "moderator" {
				newuser.NewRole = 2
			}

			if err = newuser.ReadsByUsername(); err != nil {
				pkg.Danger("No such user in db - makerequest - admin")
				Forum.ErrStatus = http.StatusBadRequest
				w.WriteHeader(http.StatusBadRequest)
				Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
				return
			}

			if err = newuser.UpdateRoleAdmin(); err != nil {
				pkg.Danger("Cannot update user role with admin")
				Forum.ErrStatus = http.StatusBadRequest
				w.WriteHeader(http.StatusBadRequest)
				Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
				return
			}
		case 2:
			role := r.FormValue("makereq")
			if role == "cancel" {
				if err := Forum.CancelModerator(); err != nil {
					pkg.Danger(err, "Cannot change role moderator")
					Forum.ErrStatus = http.StatusInternalServerError
					w.WriteHeader(http.StatusInternalServerError)
					Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
					return
				}
			}
		case 3:
			if err := Forum.WriteRequestRegisteredUserToAdmin(w, r); err != nil {
				pkg.Danger(err, "User can not write request to database")
				Forum.ErrStatus = http.StatusInternalServerError
				w.WriteHeader(http.StatusInternalServerError)
				Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
				return
			}
		}

		http.Redirect(w, r, "/makerequest", http.StatusSeeOther)
		return
	}
}
