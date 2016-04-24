package main

import (
	_ "fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"html/template"
	"log"
	"net/http"
	"github.com/shirou/gopsutil/cpu"
	"math/rand"
	"golang.org/x/net/websocket"
	"time"
	"encoding/json"
)


var store = sessions.NewCookieStore([]byte("something-very-secret2"))
var router = mux.NewRouter()
var templates = make(map[string]*template.Template)
var authkeys = make(map[string]int32)


type cpu_json struct {
	UsedPercent float64
}


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
	
	v, _ := cpu.Percent(100 * time.Millisecond, false)

	log.Println("Request for index:", session.Values["name"])

	templates["home"].Execute(response, cpu_json{UsedPercent: v[0]})
	
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


func HandleSocket(ws *websocket.Conn) {

	// Greet the client
	if err := websocket.Message.Send(ws, "Hello!"); err != nil {
		log.Println("Can't send to socket:", err)
		return
	}

	// Wait for incoming websocket messages
	for {
		
		time.Sleep(2500 * time.Millisecond)
		
		v, _ := cpu.Percent(100 * time.Millisecond, false)
		
		msgBytes, err := json.Marshal(cpu_json{UsedPercent: v[0]})
		if err != nil {
			log.Println("Can't marshal:", err)
		}
		
		msg := string(msgBytes)
		
		
		log.Println("Sending to client: ", msg)
		if err := websocket.Message.Send(ws, msg); err != nil {
			log.Println("Can't send to socket:", err)
			break
		}
		
		/*
		var reply string

		if err := websocket.Message.Receive(ws, &reply); err != nil {
			log.Println("Can't receive from socket:", err)
			break
		}

		log.Println("Received back from client: " + reply)

		msg := "Received:  " + reply
		log.Println("Sending to client: " + msg)

		if err := websocket.Message.Send(ws, msg); err != nil {
			log.Println("Can't send to socket:", err)
			break
		}
		*/
	}
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
	

	// Handlers for pages
	router.HandleFunc("/", handleIndex)
	router.HandleFunc("/login", handleLogin)
	router.HandleFunc("/logout", handleLogout)
	
	// Handler for static content
	fs := http.FileServer(http.Dir("./theme/assets"))
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fs))

	// Handler for websockets
	http.Handle("/ws", websocket.Handler(HandleSocket))

	http.Handle("/", router)
	http.ListenAndServe(":8000", nil)
}
