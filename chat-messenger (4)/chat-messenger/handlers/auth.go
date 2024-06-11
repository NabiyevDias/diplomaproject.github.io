package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"sync"

	"github.com/gorilla/sessions"
)

var (
	users     = make(map[string]string)
	userMutex = &sync.Mutex{}
	store     = sessions.NewCookieStore([]byte("secret-key"))

	// Use absolute or relative paths correctly
	loginTmpl = template.Must(template.ParseFiles("static/login/login.html"))
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		userMutex.Lock()
		users[email] = password
		userMutex.Unlock()

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	loginTmpl.Execute(w, nil)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		userMutex.Lock()
		storedPassword, ok := users[email]
		userMutex.Unlock()

		if ok && storedPassword == password {
			// Set up session
			session, _ := store.Get(r, "chat-session")
			session.Values["email"] = email
			err := session.Save(r, w)
			if err != nil {
				fmt.Println("Error saving session:", err)
			}
			http.Redirect(w, r, "/chat", http.StatusSeeOther)
			return
		}

		// Incorrect login attempt
		fmt.Println("Incorrect email or password")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	loginTmpl.Execute(w, nil)
}
