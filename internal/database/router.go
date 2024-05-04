package database

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Router struct {
	R    mux.Router
	logg logrus.Logger
	rp   *repo
}

func NewRouter(db *DataBase) *Router {
	return &Router{
		R:    *mux.NewRouter(),
		logg: *logrus.New(),
		rp:   NewRepo(db),
	}
}

func (r *Router) Pref(path string) *Router {
	r.R.PathPrefix(path + "/").Handler(http.StripPrefix(path, &r.R))
	return r
}

func (r *Router) SayHello() {
	r.R.HandleFunc("/", hello).Methods("GET")
}

func (r *Router) AuthEndPoints() {
	r.R.HandleFunc("/sign-up", r.CreateAccount).Methods("POST")
	r.R.HandleFunc("/sign-in", r.signIn).Methods("POST")
}

func (r *Router) UserEndPoints() {
	r.R.HandleFunc("/", r.getList).Methods("GET")
	r.R.HandleFunc("/{id}", r.getUserByID).Methods("GET")
	r.R.HandleFunc("/{id}", r.updateUser).Methods("PUT")
	r.R.HandleFunc("/{id}", r.partUpdateUser).Methods("PATCH")
	r.R.HandleFunc("/{id}", r.deleteUser).Methods("DELETE")
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK) // http_test.go
	io.WriteString(w, "Welcome to our Web-site!")
}

func (rout *Router) signIn(w http.ResponseWriter, r *http.Request) { // Entry
	w.WriteHeader(http.StatusOK) // http_test.go
	io.WriteString(w, "SignIn was successfully")
}

func (rout *Router) getList(w http.ResponseWriter, r *http.Request) { // GET
	w.WriteHeader(http.StatusOK) // http_test.go
	io.WriteString(w, "This is our users")
}

func (rout *Router) getUserByID(w http.ResponseWriter, r *http.Request) { // GET
	w.WriteHeader(http.StatusOK) // http_test.go
	io.WriteString(w, "This is our user by id")
}

func (rout *Router) updateUser(w http.ResponseWriter, r *http.Request) { // PUT
	w.WriteHeader(204) // http_test.go
	io.WriteString(w, "This is our updated user")
}

func (rout *Router) partUpdateUser(w http.ResponseWriter, r *http.Request) { // PATCH
	w.WriteHeader(204) // http_test.go
	io.WriteString(w, "This is our part update user")
}

func (rout *Router) deleteUser(w http.ResponseWriter, r *http.Request) { // DELETE
	w.WriteHeader(204) // http_test.go
	io.WriteString(w, "This is our deleted user")
}
