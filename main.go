package main

import (
	_ "fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"html/template"
	"log"
	"net/http"
	"github.com/shirou/gopsutil/mem"
	"math/rand"
)


var store = sessions.NewCookieStore([]byte("something-very-secret2"))
var router = mux.NewRouter()
var templates = make(map[string]*template.Template)
var authkeys = make(map[string]int32)

func handleIndex(response http.ResponseWriter, request *http.Request) {

	session, err := store.Get(request, "gosession")
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
	
	if session.Values["name"] == "" || session.Values["name"] == nil {
		http.Redirect(response, request, "/login", 302)
		return
	}
	
	// Add authentication key
	log.Println(session.Values["name"])
	authkeys[session.Values["name"].(string)] = rand.Int31()
	
	v, _ := mem.VirtualMemory()

	log.Println("Request for index:", session.Values["name"])

	templates["home"].Execute(response, v)
	
}

func handleLogin(response http.ResponseWriter, request *http.Request) {

	if request.Method == "POST" {
		name := request.FormValue("name")
		password := request.FormValue("password")

		// Some stupid check :)
		if password == name {
			session, _ := store.Get(request, "gosession")
			session.Values["name"] = name
			session.Save(request, response)
			
			// redirect to homepage
			http.Redirect(response, request, "/", 302)
		}

	}

	templates["login"].Execute(response, nil)
}

func handleLogout(response http.ResponseWriter, request *http.Request) {

	session, _ := store.Get(request, "gosession")
	session.Values = nil
	session.Save(request, response)
	
	// redirect to homepage
	http.Redirect(response, request, "/", 302)

}

func main() {
	
	// Initialize all template files
	for _,templateName := range []string{"home", "login"} {
		t, err := template.ParseFiles("theme/" + templateName + ".html")
		if err != nil {
			log.Fatal("Couldn't load template file:", err)
		}
		templates[templateName] = t
	}
	

	router.HandleFunc("/", handleIndex)

	router.HandleFunc("/login", handleLogin)
	router.HandleFunc("/logout", handleLogout)

	http.Handle("/", router)
	http.ListenAndServe(":8000", nil)
}
