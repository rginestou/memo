package main

import (
	"html/template"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type loginData struct {
	Title string
}

func loginGET(w http.ResponseWriter, r *http.Request) {
	// Load templates
	tmpl, _ := template.ParseFiles(
		"view/login.html",
		"view/head.html",
	)

	tmpl.ExecuteTemplate(w, "login", loginData{
		Title: "Memo • Login",
	})
}

func loginPOST(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	// Auth
	r.ParseForm()

	login := r.FormValue("login")
	pwd := r.FormValue("password")

	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if login == config.Login && bcrypt.CompareHashAndPassword(hash, []byte(config.Password)) == nil {
		// Set user as authenticated
		session.Values["authenticated"] = true
		session.Save(r, w)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	// Revoke users authentication
	session.Values["authenticated"] = false
	session.Save(r, w)
}

func sessionAuth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "cookie-name")

		// Check if user is authenticated
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		handler(w, r)
	}
}
