package routes

import (
	"forum/model"
	"forum/pkg"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

//Signup страница регистрации
func Signup(w http.ResponseWriter, r *http.Request) {
	Forum := model.Forum{}

	cookie, err := r.Cookie("session")
	if err == nil {
		if err := Forum.User.ReadByUiid(cookie.Value); err == nil {
			pkg.Danger("User already have cookie")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}

	// Check Method ...
	if r.Method != http.MethodGet {
		pkg.Danger(invMeth + "signup page")
		Forum.ErrStatus = http.StatusMethodNotAllowed
		w.WriteHeader(http.StatusMethodNotAllowed)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	signupView, err := pkg.NewView("views/signup.html")
	if err != nil {
		pkg.Danger("Cannot execute signup template")
		Forum.ErrStatus = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	signupView.ExecuteTemplate(w, "bootstrap", nil)
}

//Login страница логина
func Login(w http.ResponseWriter, r *http.Request) {
	Forum := model.Forum{}

	cookie, err := r.Cookie("session")
	if err == nil {
		if err := Forum.User.ReadByUiid(cookie.Value); err == nil {
			pkg.Danger("User already have cookie")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}

	// Check Method ...
	if r.Method != http.MethodGet {
		pkg.Danger(invMeth + "login page")
		Forum.ErrStatus = http.StatusMethodNotAllowed
		w.WriteHeader(http.StatusMethodNotAllowed)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	loginView, err := pkg.NewView("views/login.html")
	if err != nil {
		pkg.Danger("Cannot execute login template")
		Forum.ErrStatus = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	loginView.ExecuteTemplate(w, "bootstrap", nil)
}

//Signedup регистрирует нового юзера
func Signedup(w http.ResponseWriter, r *http.Request) {
	Forum := model.Forum{
		User: model.User{
			Username: r.FormValue("username"),
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
		},
	}

	// Check Method ...
	if r.Method != http.MethodPost {
		pkg.Danger(invMeth + "posts dedicated to a specific topic page")
		Forum.ErrStatus = http.StatusMethodNotAllowed
		w.WriteHeader(http.StatusMethodNotAllowed)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	// Check Valid Characters ...
	if err := pkg.CharactersValidate(Forum.User.Username); err != nil {
		pkg.Danger(err, "Characters are not validate")
		Forum.IsErr = true
		Forum.ErrMsg = "Username is consists of " + err.Error()
		loginView, err := pkg.NewView("views/signup.html")
		if err != nil {
			pkg.Danger("Cannot execute login template")
			Forum.ErrStatus = http.StatusInternalServerError
			w.WriteHeader(http.StatusInternalServerError)
			Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		loginView.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	if err := pkg.EmailValidate(Forum.User.Email); err != nil {
		pkg.Danger(err, "Characters are not validate")
		Forum.IsErr = true
		Forum.ErrMsg = "Email is consists of " + err.Error()
		loginView, err := pkg.NewView("views/signup.html")
		if err != nil {
			pkg.Danger("Cannot execute signup template")
			Forum.ErrStatus = http.StatusInternalServerError
			w.WriteHeader(http.StatusInternalServerError)
			Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		loginView.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	if err := Forum.User.WriteNewUser(); err != nil {
		pkg.Danger(err, "Cannot create a user")
		Forum.IsErr = true
		Forum.ErrMsg = err.Error()
		Forum.ErrMsg = Forum.ErrMsg[31:] + " already exists"
		loginView, err := pkg.NewView("views/signup.html")
		if err != nil {
			pkg.Danger("Cannot execute login template")
			Forum.ErrStatus = http.StatusInternalServerError
			w.WriteHeader(http.StatusInternalServerError)
			Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		loginView.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	if err := Forum.User.WriteSession(w); err != nil {
		pkg.Danger(err, "Cannot create session")
		Forum.IsErr = true
		Forum.ErrMsg = err.Error()
		loginView, err := pkg.NewView("views/signup.html")
		if err != nil {
			pkg.Danger("Cannot execute login template")
			Forum.ErrStatus = http.StatusInternalServerError
			w.WriteHeader(http.StatusInternalServerError)
			Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		loginView.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	http.Redirect(w, r, "/", 301)
}

//Logedin проверят на валидность логина
func Logedin(w http.ResponseWriter, r *http.Request) {
	Forum := model.Forum{
		User: model.User{
			Username: r.FormValue("username"),
			Password: r.FormValue("password"),
		},
	}

	//Check Method ...
	if r.Method != http.MethodPost {
		pkg.Danger(invMeth + "posts dedicated to a specific topic page")
		Forum.ErrStatus = http.StatusMethodNotAllowed
		w.WriteHeader(http.StatusMethodNotAllowed)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	// Check Valid Characters ...
	if err := pkg.CharactersValidate(Forum.User.Username); err != nil {
		pkg.Danger(err, "Characters are not validate")
		Forum.IsErr = true
		Forum.ErrMsg = "You typed invalid character"
		loginView, err := pkg.NewView("views/login.html")
		if err != nil {
			pkg.Danger("Cannot execute login template")
			Forum.ErrStatus = http.StatusInternalServerError
			w.WriteHeader(http.StatusInternalServerError)
			Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		loginView.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	// Check if there is such user in database ...
	if err := Forum.User.ReadByUsername(); err != nil {
		pkg.Danger(err, "Cannot find user")
		Forum.IsErr = true
		Forum.ErrMsg = "You typed incorrect username or password"
		loginView, err := pkg.NewView("views/login.html")
		if err != nil {
			pkg.Danger("Cannot execute login template")
			Forum.ErrStatus = http.StatusInternalServerError
			w.WriteHeader(http.StatusInternalServerError)
			Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		loginView.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	// Check if the password is correct
	if err := bcrypt.CompareHashAndPassword(Forum.User.HashPassword, []byte(Forum.User.Password)); err != nil {
		pkg.Danger(err, "Incorrect password")
		Forum.IsErr = true
		Forum.ErrMsg = "You typed incorrect username or password"
		loginView, err := pkg.NewView("views/login.html")
		if err != nil {
			pkg.Danger("Cannot execute login template")
			Forum.ErrStatus = http.StatusInternalServerError
			w.WriteHeader(http.StatusInternalServerError)
			Errtmpl.ExecuteTemplate(w, "bootstrap", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		loginView.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	session := model.Session{AuthID: Forum.User.ID}
	if err := session.DeleteSessionFromDB(); err != nil {
		pkg.Info(err, "No other session in db")
	}

	err := Forum.User.WriteSession(w)
	if err != nil {
		pkg.Danger(err, "Cannot create session")
		http.Redirect(w, r, "/login", 302)
		return
	}

	http.Redirect(w, r, "/", 301)
}

//Logout для выхода из сессии
func Logout(w http.ResponseWriter, r *http.Request) {
	Forum := model.Forum{}

	cookie, err := r.Cookie("session")
	if err != nil {
		pkg.Danger(err, "No cookie")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	session := model.Session{
		UUID: cookie.Value,
	}

	if err = session.ReadAuthIdbyUUID(); err != nil {
		model.DropCookie(w, r)
		pkg.Danger(err, "Cannot get authID from db")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if err = session.DeleteSessionFromDB(); err != nil {
		model.DropCookie(w, r)
		pkg.Danger(err, "Cannot delete cookie from db")
		Forum.ErrStatus = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	if err := model.DropCookie(w, r); err != nil {
		pkg.Danger(err, "Cannot drop cookie")
		Forum.ErrStatus = http.StatusBadRequest
		w.WriteHeader(http.StatusBadRequest)
		Errtmpl.ExecuteTemplate(w, "bootstrap", Forum)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
