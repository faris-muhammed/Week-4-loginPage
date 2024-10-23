package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/icza/session"
)

var temp *template.Template
var err string

func init() {
	temp = template.Must(template.ParseGlob("template/*.html"))
}
func clearCache(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, no-transform, must-revalidate, private, max-age=0")
}
func loginHandler(w http.ResponseWriter, r *http.Request) {
	clearCache(w, r)
	sess := session.Get(r)
	if sess == nil {
		temp.ExecuteTemplate(w, "index.html", err)
		err = ""
	} else {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}
}
func loginCheck(w http.ResponseWriter, r *http.Request) {
	clearCache(w, r)
	user := r.FormValue("username")
	pass := r.FormValue("password")
	if user == "faris_muhd" && pass == "faris@123" {
		sess := session.NewSessionOptions(&session.SessOptions{
			CAttrs: map[string]interface{}{"username": user},
		})
		session.Add(sess, w)
		fmt.Println("session is ", sess)
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	} else {
		err = "Invalid Username or Password"
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
func homeHandler(w http.ResponseWriter, r *http.Request) {
	clearCache(w, r)
	sess := session.Get(r)
	if sess == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		user := sess.CAttr("username")
		temp.ExecuteTemplate(w, "home.html", user)
	}
}
func logoutHandle(w http.ResponseWriter, r *http.Request) {
	clearCache(w, r)
	sess := session.Get(r)
	if sess != nil {
		session.Remove(sess, w)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func main() {
	http.HandleFunc("/", loginHandler)
	http.HandleFunc("/login", loginCheck)
	http.HandleFunc("/home", homeHandler)
	http.HandleFunc("/logout", logoutHandle)
	http.ListenAndServe(":9999", nil)
}
