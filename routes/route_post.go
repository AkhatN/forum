package routes

import (
	"forum/model"
	"forum/pkg"
	"net/http"
	"strconv"
	"time"
)

//Createpost creates a new post
func Createpost(w http.ResponseWriter, r *http.Request) {
	Forum := model.Forum{}

	var err error
	cookie := r.Context().Value("user").(string)
	Forum.User.IsStatus = true
	if err = Forum.User.ReadByUiid(cookie); err != nil {
		model.DropCookie(w, r)
		pkg.Danger(err, usUUID)
		Forum.User.IsStatus = false
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

	if Forum.User.ID == 0 {
		pkg.Danger(err, "No user to create post")
		Forum.User.IsStatus = false
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	postView, err := pkg.NewView("views/post.html")
	if err != nil {
		pkg.Danger("Cannot execute post template")
		Forum.ErrStatus = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	postView.ExecuteTemplate(w, "bootstrap", Forum)
}

//Savepost saves a new post
func Savepost(w http.ResponseWriter, r *http.Request) {
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

	// Check if empty string is sent
	if err := pkg.EmptySpaceCheck(r.FormValue("title"), r.FormValue("text")); err != nil {
		pkg.Danger(err, "Characters are not validate")
		Forum.IsErr = true
		Forum.ErrMsg = err.Error()
		postView, err := pkg.NewView("views/post.html")
		if err != nil {
			pkg.Danger("Cannot execute post template")
			Forum.ErrStatus = http.StatusInternalServerError
			w.WriteHeader(http.StatusInternalServerError)
			Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
			return
		}

		w.WriteHeader(400)
		postView.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	// Read data categories from client
	cs := r.MultipartForm.Value["addmore[]"]

	// Check if empty string is sent
	if err := pkg.EmptySpaceCheck(cs...); err != nil {
		pkg.Danger(err, "Characters are not validate")
		Forum.IsErr = true
		Forum.ErrMsg = err.Error()
		postView, err := pkg.NewView("views/post.html")
		if err != nil {
			pkg.Danger("Cannot execute post template")
			Forum.ErrStatus = http.StatusInternalServerError
			w.WriteHeader(http.StatusInternalServerError)
			Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
			return
		}

		w.WriteHeader(400)
		postView.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	// Writing data to struct
	Forum.Post = model.Post{
		Title:    r.FormValue("title"),
		FullText: r.FormValue("text"),
		AuthID:   Forum.User.ID,
		Date:     time.Now().Format("2006.01.02 15:04:05"),
	}

	//saving an image to the local machine before saving anything to database
	if Forum.Post.Image, err = pkg.ImageUpload(w, r); err != nil {
		pkg.Danger(err, "Cannot save image path in database")
		Forum.IsErr = true
		Forum.ErrMsg = err.Error()
		postView, err := pkg.NewView("views/post.html")
		if err != nil {
			pkg.Danger("Cannot execute post template")
			Forum.ErrStatus = http.StatusInternalServerError
			w.WriteHeader(http.StatusInternalServerError)
			Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
			return
		}

		w.WriteHeader(400)
		postView.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	// Writing post to database
	if err := Forum.WritePosttoDataBase(); err != nil {
		pkg.Danger(err, "Limited amount of characters")
		Forum.IsErr = true
		Forum.ErrMsg = "Limited amount of characters in title(40)"
		postView, err := pkg.NewView("views/post.html")
		if err != nil {
			pkg.Danger("Cannot execute post template")
			Forum.ErrStatus = http.StatusInternalServerError
			w.WriteHeader(http.StatusInternalServerError)
			Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
			return
		}

		w.WriteHeader(400)
		postView.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	cats := []model.Category{}

	//Checking and saving only unique categories
	pkg.CheckUniq(&cats, cs)

	if err := Forum.WriteCats(cats); err != nil {
		pkg.Danger(err, "Cannot save categories in database")
		Forum.IsErr = true
		Forum.ErrMsg = "Internal server error. Please try to create the post later."
		postView, err := pkg.NewView("views/post.html")
		if err != nil {
			pkg.Danger("Cannot execute post template")
			Forum.ErrStatus = http.StatusInternalServerError
			w.WriteHeader(http.StatusInternalServerError)
			Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
			return
		}

		w.WriteHeader(500)
		postView.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	if Forum.Post.Image != "" {
		// Writing an image path to database
		if err := Forum.Post.WriteImagetoDb(); err != nil {
			pkg.Danger(err, "Cannot save image path in database")
			Forum.IsErr = true
			Forum.ErrMsg = "Internal server error. Please try to create the post later."
			postView, err := pkg.NewView("views/post.html")
			if err != nil {
				pkg.Danger("Cannot execute post template")
				Forum.ErrStatus = http.StatusInternalServerError
				w.WriteHeader(http.StatusInternalServerError)
				Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
				return
			}

			w.WriteHeader(500)
			postView.ExecuteTemplate(w, "bootstrap", Forum)
			return
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

//Title shows a page of specific post
func Title(w http.ResponseWriter, r *http.Request) {
	Forum := model.Forum{}
	var err error

	cookie, err := r.Cookie("session")
	if err == nil {
		if err := Forum.User.ReadByUiid(cookie.Value); err != nil {
			model.DropCookie(w, r)
			Forum.User.IsStatus = false
		}
	}

	if Forum.User.IsStatus {
		if err := Forum.ReadCountNotification(); err != nil {
			pkg.Danger("Cannot get notifications")
			Forum.ErrStatus = http.StatusInternalServerError
			w.WriteHeader(http.StatusInternalServerError)
			Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
			return
		}
	}

	Forum.Post.ID, err = strconv.Atoi(r.RequestURI[7:])

	if err != nil {
		pkg.Danger(err, "Not found: title")
		Forum.ErrStatus = http.StatusNotFound
		w.WriteHeader(http.StatusNotFound)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	if err := Forum.ReadPostFromDataBase(w, r); err != nil {
		pkg.Danger(err, "Cannot read a post")
		Forum.ErrStatus = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	if Forum.Post.AuthName == "" || Forum.Post.ID == 0 {
		pkg.Danger(err, "Page not found: No such post")
		Forum.ErrStatus = http.StatusNotFound
		w.WriteHeader(http.StatusNotFound)
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

	Forum.Post.Image, err = Forum.Post.ReadImagefromDb()
	if err != nil {
		Forum.Post.IsImage = false
	}

	articleView, err := pkg.NewView("views/article.html")
	if err != nil {
		pkg.Danger(err, "Cannot execute article template")
		Forum.ErrStatus = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	articleView.ExecuteTemplate(w, "bootstrap", Forum)
}
