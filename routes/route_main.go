package routes

import (
	"forum/model"
	"forum/pkg"
	"html/template"
	"log"
	"net/http"
)

const (
	usUUID  = "Cannot get user by UUID"
	invMeth = "Invalid method "
	pnFound = "Page not found"
	erTemp  = "Cannot execute error template"
)

//Errtmpl ...
var Errtmpl *template.Template

func init() {
	var err error
	Errtmpl, err = pkg.NewView("views/error.html")
	if err != nil {
		log.Fatal(err)
	}
}

//Home shows home page with posts
func Home(w http.ResponseWriter, r *http.Request) {
	Forum := model.Forum{}

	cookie, err := r.Cookie("session")
	if err == nil {
		if err := Forum.User.ReadByUiid(cookie.Value); err != nil {
			model.DropCookie(w, r)
			Forum.User.IsStatus = false
		}
	}

	if err := Forum.ReadDataFromDataBase(w, r); err != nil {
		pkg.Danger("Cannot get posts from database")
		Forum.ErrStatus = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	if r.URL.Path != "/" {
		pkg.Danger(pnFound)
		Forum.ErrStatus = http.StatusNotFound
		w.WriteHeader(http.StatusNotFound)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	if r.Method != http.MethodGet {
		pkg.Danger(invMeth + "home page")
		Forum.ErrStatus = http.StatusMethodNotAllowed
		w.WriteHeader(http.StatusMethodNotAllowed)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	homeView, err := pkg.NewView("views/home.html")
	if err != nil {
		pkg.Danger("Cannot execute home template of all posts")
		Forum.ErrStatus = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	homeView.ExecuteTemplate(w, "bootstrap", Forum)
}
