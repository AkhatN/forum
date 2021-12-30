package routes

import (
	"forum/model"
	"forum/pkg"
	"net/http"
)

//Cats фильтрирует по категориям
func Cats(w http.ResponseWriter, r *http.Request) {
	Forum := model.Forum{}
	Forum.Post.Title = r.RequestURI[6:]

	cookie, err := r.Cookie("session")
	if err == nil {
		if err := Forum.User.ReadByUiid(cookie.Value); err != nil {
			model.DropCookie(w, r)
			Forum.User.IsStatus = false
		}
	}

	if err := Forum.ReadCatPostsFromDataBase(w, r); err != nil {
		pkg.Danger(err, "Cannot get posts dedicated to a specific topic")
		Forum.ErrStatus = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	if r.URL.RequestURI() == "/cats/" {
		pkg.Danger("Not found: cats")
		Forum.ErrStatus = http.StatusNotFound
		w.WriteHeader(http.StatusNotFound)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	if r.Method != http.MethodGet {
		pkg.Danger(invMeth + "posts dedicated to a specific topic page")
		Forum.ErrStatus = http.StatusMethodNotAllowed
		w.WriteHeader(http.StatusMethodNotAllowed)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	if len(Forum.Posts) < 1 {
		pkg.Info(err, "Page not found: No such posts dedicated to that category")
		Forum.ErrStatus = http.StatusOK
		w.WriteHeader(http.StatusOK)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	homeView, err := pkg.NewView("views/home.html")
	if err != nil {
		pkg.Danger(err, "Cannot execute home template of posts dedicated to a specific topic")
		Forum.ErrStatus = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	homeView.ExecuteTemplate(w, "bootstrap", Forum)
}

//Liked фильтрирует по постам, которые лайкнул юзер
func Liked(w http.ResponseWriter, r *http.Request) {
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
	if err := Forum.ReadLikedPostsFromDataBase(w, r); err != nil {
		pkg.Danger(err, "Cannot get liked posts of registered user")
		Forum.ErrStatus = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	homeView, err := pkg.NewView("views/home.html")
	if err != nil {
		pkg.Danger(err, "Cannot execute home template of liked posts")
		Forum.ErrStatus = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	homeView.ExecuteTemplate(w, "bootstrap", Forum)
}

//Mine фильтрирует по постам, которые создал юзер
func Mine(w http.ResponseWriter, r *http.Request) {
	Forum := model.Forum{}

	cookie, err := r.Cookie("session")
	if err != nil {
		pkg.Danger(err, "Cannot get mine posts of registered user")
		Forum.ErrStatus = http.StatusNotFound
		w.WriteHeader(http.StatusNotFound)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	if err := Forum.User.ReadByUiid(cookie.Value); err != nil {
		model.DropCookie(w, r)
		pkg.Danger(err, "Cannot get mine posts of registered user")
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

	if err := Forum.ReadMyPostsFromDataBase(w, r); err != nil {
		pkg.Danger(err, "Cannot get created posts of registered user")
		Forum.ErrStatus = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	homeView, err := pkg.NewView("views/home.html")
	if err != nil {
		pkg.Danger(err, "Cannot execute home template of created posts")
		Forum.ErrStatus = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	homeView.ExecuteTemplate(w, "bootstrap", Forum)
}
