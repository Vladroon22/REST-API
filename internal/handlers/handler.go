package handlers

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

type Router struct {
	R *mux.Router
}

func NewRouter() Router {
	return Router{
		R: mux.NewRouter(),
	}
}

func (r *Router) Pref(path string) Router {
	r.R.PathPrefix(path + "/").Handler(http.StripPrefix(path, r.R))
	return *r
}

func (r Router) SayHello() {
	r.R.HandleFunc("/", hello).Methods("GET")
}

func (r Router) EndPoints() {
	r.R.HandleFunc("/sign-up", signUp).Methods("GET")
	r.R.HandleFunc("/sign-in", signIn).Methods("GET")
	r.R.HandleFunc("/logout", logOut).Methods("GET")
}

func (r Router) UserEndPoints() {
	r.R.HandleFunc("/users", getList).Methods("GET")
	r.R.HandleFunc("/users", createUser).Methods("POST")
	r.R.HandleFunc("/users/{id}", getUserByID).Methods("GET")
	r.R.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	r.R.HandleFunc("/users/{id}", partUpdateUser).Methods("PATCH")
	r.R.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK) // http_test.go
	io.WriteString(w, "Welcome to our Web-site!")
}

func signUp(w http.ResponseWriter, r *http.Request) { // register
	w.WriteHeader(http.StatusOK) // http_test.go
	io.WriteString(w, "You have registered")
}

func signIn(w http.ResponseWriter, r *http.Request) { // Entry
	w.WriteHeader(http.StatusOK) // http_test.go
	io.WriteString(w, "SignIn was successfully")
}

func logOut(w http.ResponseWriter, r *http.Request) { //logOut
	w.WriteHeader(http.StatusOK) // http_test.go
	io.WriteString(w, "You have been logout")
}

func getList(w http.ResponseWriter, r *http.Request) { // GET
	w.WriteHeader(http.StatusOK) // http_test.go
	io.WriteString(w, "This is our users")
}

func createUser(w http.ResponseWriter, r *http.Request) { // POST
	w.WriteHeader(201) // http_test.go
	io.WriteString(w, "This is our New user")
}

func getUserByID(w http.ResponseWriter, r *http.Request) { // GET
	w.WriteHeader(http.StatusOK) // http_test.go
	io.WriteString(w, "This is our user by id")
}

func updateUser(w http.ResponseWriter, r *http.Request) { // PUT
	w.WriteHeader(204) // http_test.go
	io.WriteString(w, "This is our updated user")
}

func partUpdateUser(w http.ResponseWriter, r *http.Request) { // PATCH
	w.WriteHeader(204) // http_test.go
	io.WriteString(w, "This is our part update user")
}

func deleteUser(w http.ResponseWriter, r *http.Request) { // DELETE
	w.WriteHeader(204) // http_test.go
	io.WriteString(w, "This is our deleted user")
}
