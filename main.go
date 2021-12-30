package main

import (
	"crypto/tls"
	"database/sql"
	"forum/model"
	"forum/pkg"
	"forum/routes"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func init() {
	// Open DataBase ...
	db, err := sql.Open("sqlite3", "sqlite/forum.db?_foreign_keys=on")
	if err != nil {
		log.Println(err)
		return
	}

	model.Db = db

	// Creating tables in sqlite database
	if err = model.InitSQL(); err != nil {
		log.Println(err)
		return
	}

	//creates admin
	err = model.WriteAdmin()
	if err != nil {
		log.Println(err)
		return
	}

	// Creating a folder
	err = os.MkdirAll("./static/img_posts", os.ModePerm)
	if err != nil {
		log.Println(err)
		return
	}

	// Loading config file
	pkg.LoadConfig()
	//opening file for logging
	file, err := os.OpenFile("forum.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file", err)
	}
	pkg.Logger = log.New(file, "INFO ", log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {
	pkg.CreateSSL()

	mux := http.NewServeMux()
	files := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", files))

	mux.HandleFunc("/", routes.Home)

	mux.HandleFunc("/login", routes.Login)
	mux.HandleFunc("/logedin", routes.Logedin)
	mux.HandleFunc("/signup", routes.Signup)
	mux.HandleFunc("/signedup", routes.Signedup)
	mux.HandleFunc("/logout", routes.Logout)

	mux.HandleFunc("/title/", routes.Title)
	mux.Handle("/createpost", pkg.Middleware(http.HandlerFunc(routes.Createpost)))
	mux.Handle("/savepost", pkg.Middleware(http.HandlerFunc(routes.Savepost)))

	mux.Handle("/savecomment", pkg.Middleware(http.HandlerFunc(routes.Savecomment)))
	mux.Handle("/like", pkg.Middleware(http.HandlerFunc(routes.Like)))
	mux.Handle("/comlike", pkg.Middleware(http.HandlerFunc(routes.ComLike)))

	mux.HandleFunc("/cats/", routes.Cats)
	mux.HandleFunc("/liked", routes.Liked)
	mux.HandleFunc("/mine", routes.Mine)

	mux.HandleFunc("/deletePost", routes.DeletePost)
	mux.HandleFunc("/deleteComm", routes.DeleteComm)
	mux.HandleFunc("/editPost/", routes.EditPost)
	mux.HandleFunc("/editComm/", routes.EditComm)
	mux.HandleFunc("/editedPost/", routes.EditedPost)
	mux.HandleFunc("/editedComm/", routes.EditedComm)

	mux.HandleFunc("/notification", routes.Notification)
	mux.HandleFunc("/activity", routes.Activity)

	mux.HandleFunc("/request", routes.Request)
	mux.HandleFunc("/handlerequest/", routes.HandleRequest)
	mux.HandleFunc("/saverequest", routes.Saverequest)

	mux.HandleFunc("/deleterequest", routes.DeleteRequest)
	mux.HandleFunc("/report", routes.Report)
	mux.HandleFunc("/makerequest", routes.Makerequest)

	mux.HandleFunc("/managecats", routes.Managecats)
	mux.HandleFunc("/savecats", routes.SaveCats)
	mux.HandleFunc("/deletecats", routes.DeleteCats)

	mycert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		log.Fatal(err)
	}

	server := &http.Server{
		Addr:         pkg.Config.Address,
		Handler:      mux,
		ReadTimeout:  10 * time.Minute,
		WriteTimeout: 10 * time.Minute,
		IdleTimeout:  1 * time.Hour,
		TLSConfig: &tls.Config{
			MinVersion:               tls.VersionTLS13,
			PreferServerCipherSuites: true,
			Certificates:             []tls.Certificate{mycert},
		},
		ErrorLog: log.New(ioutil.Discard, "Cluster-Http-Server", 0),
	}

	go func() {
		//запусакаем http сервер в горутине и перенаправляем все http запросы на https сервер
		if err := http.ListenAndServe(":9090", http.HandlerFunc(redirectTLS)); err != nil {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	log.Printf("Listening on %s port ...\n", ":9090")

	// //запускаем https сервер
	if err := server.ListenAndServeTLS("", ""); err != nil {
		log.Fatal(err)

	}
}

func redirectTLS(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://localhost:8080"+r.RequestURI, http.StatusMovedPermanently)
}
