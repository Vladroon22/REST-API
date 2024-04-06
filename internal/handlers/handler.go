package handlers

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

type Router struct {
	R *mux.Router
	// L *logrus.Logger
}

func NewRouter() Router {
	return Router{
		R: mux.NewRouter(),
		//		L: logrus.New(),
	}
}

func (r *Router) Pref(path string) Router {
	var handler http.Handler = r.R
	r.R.Handle(path+"/", http.StripPrefix(path, handler))
	return *r
}

func (r Router) SayHello() {
	r.R.HandleFunc("/", Hello).Methods("GET")
	// r.L.Infoln("Handler -> Hello")
}

func (r Router) EndPoints() {
	r.R.HandleFunc("/sign-up", SignUp).Methods("GET")
	r.R.HandleFunc("/sign-in", SignIn).Methods("GET")
	r.R.HandleFunc("/logout", LogOut).Methods("GET")
	// r.L.Infoln("Handlers of AuthEndPoints")
}

func (r Router) UserEndPoints() {
	r.R.HandleFunc("/users", GetList).Methods("GET")
	r.R.HandleFunc("/users", CreateUser).Methods("POST")
	r.R.HandleFunc("/users/{id}", GetUserByID).Methods("GET")
	r.R.HandleFunc("/users/{id}", UpdateUser).Methods("PUT")
	r.R.HandleFunc("/users/{id}", PartUpdateUser).Methods("PATCH")
	r.R.HandleFunc("/users/{id}", DeleteUser).Methods("DELETE")
	// r.L.Infoln("Handlers of UserEndPoints")
}

func Hello(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK) // http_test.go
	io.WriteString(w, "Welcome to our Web-site!")
}

func SignUp(w http.ResponseWriter, r *http.Request) { // register
	w.WriteHeader(http.StatusOK) // http_test.go
	io.WriteString(w, "You have registered")
}

func SignIn(w http.ResponseWriter, r *http.Request) { // Entry
	w.WriteHeader(http.StatusOK) // http_test.go
	io.WriteString(w, "SignIn was successfully")
}

func LogOut(w http.ResponseWriter, r *http.Request) { //logOut
	w.WriteHeader(http.StatusOK) // http_test.go
	io.WriteString(w, "You have been logout")
}

func GetList(w http.ResponseWriter, r *http.Request) { // GET
	w.WriteHeader(http.StatusOK) // http_test.go
	io.WriteString(w, "This is our users")
}

func CreateUser(w http.ResponseWriter, r *http.Request) { // POST
	w.WriteHeader(201) // http_test.go
	io.WriteString(w, "This is our New user")
}

func GetUserByID(w http.ResponseWriter, r *http.Request) { // GET
	w.WriteHeader(http.StatusOK) // http_test.go
	io.WriteString(w, "This is our user by id")
}

func UpdateUser(w http.ResponseWriter, r *http.Request) { // PUT
	w.WriteHeader(204) // http_test.go
	io.WriteString(w, "This is our updated user")
}

func PartUpdateUser(w http.ResponseWriter, r *http.Request) { // PATCH
	w.WriteHeader(204) // http_test.go
	io.WriteString(w, "This is our part update user")
}

func DeleteUser(w http.ResponseWriter, r *http.Request) { // DELETE
	w.WriteHeader(204) // http_test.go
	io.WriteString(w, "This is our deleted user")
}
