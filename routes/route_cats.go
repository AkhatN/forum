package routes

import (
	"forum/model"
	"forum/pkg"
	"net/http"
	"strconv"
)

//Managecats manages categories
func Managecats(w http.ResponseWriter, r *http.Request) {
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

	if err := Forum.ReadCountNotification(); err != nil {
		pkg.Danger("Cannot get notifications")
		Forum.ErrStatus = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	if r.Method != http.MethodGet {
		pkg.Danger(invMeth + "manage cats")
		w.WriteHeader(http.StatusMethodNotAllowed)
		Forum.ErrStatus = http.StatusMethodNotAllowed
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	if Forum.User.Role != 1 {
		pkg.Danger("User cannot manage categories")
		Forum.ErrStatus = http.StatusNotFound
		w.WriteHeader(http.StatusNotFound)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}
	if err := Forum.Post.ReadAllCats(); err != nil {
		pkg.Danger("Cannot get categories")
		Forum.ErrStatus = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	managecats, err := pkg.NewView("views/categories.html")
	if err != nil {
		pkg.Danger("Cannot execute categories template")
		Forum.ErrStatus = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	managecats.ExecuteTemplate(w, "bootstrap", Forum)

}

//SaveCats creates categories
func SaveCats(w http.ResponseWriter, r *http.Request) {
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
		pkg.Danger(invMeth + "manage cats")
		w.WriteHeader(http.StatusMethodNotAllowed)
		Forum.ErrStatus = http.StatusMethodNotAllowed
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	if Forum.User.Role != 1 {
		pkg.Danger("User cannot manage categories")
		Forum.ErrStatus = http.StatusNotFound
		w.WriteHeader(http.StatusNotFound)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}
	cat := model.Category{
		Name: r.FormValue("savecat"),
	}
	if err := pkg.EmptySpaceCheck(cat.Name); err != nil {
		pkg.Danger(err, "Cannot create empty category")
		Forum.IsErr = true
		Forum.ErrMsg = err.Error()
		postView, err := pkg.NewView("views/categories.html")
		if err != nil {
			pkg.Danger("Cannot execute categories template")
			Forum.ErrStatus = http.StatusInternalServerError
			w.WriteHeader(http.StatusInternalServerError)
			Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
			return
		}

		w.WriteHeader(400)
		postView.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	if err := cat.SaveNewCat(); err != nil {
		pkg.Danger(err, "Cannot save new category in database")
		Forum.IsErr = true
		Forum.ErrMsg = "Internal server error. Please try to create category later."
		postView, err := pkg.NewView("views/categories.html")
		if err != nil {
			pkg.Danger("Cannot execute categories template")
			Forum.ErrStatus = http.StatusInternalServerError
			w.WriteHeader(http.StatusInternalServerError)
			Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
			return
		}

		w.WriteHeader(500)
		postView.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	http.Redirect(w, r, "/managecats", http.StatusSeeOther)
}

//DeleteCats deletes categories
func DeleteCats(w http.ResponseWriter, r *http.Request) {
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
		pkg.Danger(invMeth + "manage cats")
		w.WriteHeader(http.StatusMethodNotAllowed)
		Forum.ErrStatus = http.StatusMethodNotAllowed
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	if Forum.User.Role != 1 {
		pkg.Danger("User cannot manage categories")
		Forum.ErrStatus = http.StatusNotFound
		w.WriteHeader(http.StatusNotFound)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	cat, err := strconv.Atoi(r.FormValue("deletecat"))
	if err != nil {
		pkg.Danger("Cannot delete category - invalid catid")
		Forum.ErrStatus = http.StatusBadRequest
		w.WriteHeader(http.StatusBadRequest)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	if err := Forum.DeleteCatsPosts(cat); err != nil {
		pkg.Danger("Cannot delete categories - manage categories")
		Forum.ErrStatus = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	http.Redirect(w, r, "/managecats", http.StatusSeeOther)
}
