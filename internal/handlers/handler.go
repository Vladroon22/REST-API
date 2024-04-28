package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Vladroon22/REST-API/internal/database"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Router struct {
	R    *mux.Router
	logg *logrus.Logger
	db   *database.DataBase
}

func NewRouter() *Router {
	return &Router{
		R:    mux.NewRouter(),
		logg: logrus.New(),
		db:   &database.DataBase{},
	}
}

func (r *Router) Pref(path string) *Router {
	r.R.PathPrefix(path + "/").Handler(http.StripPrefix(path, r.R))
	return r
}

func (r *Router) SayHello() {
	r.R.HandleFunc("/", hello).Methods("GET")
}

func (r *Router) EndPoints() {
	r.R.HandleFunc("/sign-up", r.CreateAccount).Methods("POST")
	r.R.HandleFunc("/sign-in", r.signIn).Methods("POST")
	r.R.HandleFunc("/logout", r.logOut).Methods("POST")
}

func (r *Router) UserEndPoints() {
	r.R.HandleFunc("/users", r.getList).Methods("GET")
	r.R.HandleFunc("/users/{id}", r.getUserByID).Methods("GET")
	r.R.HandleFunc("/users/{id}", r.updateUser).Methods("PUT")
	r.R.HandleFunc("/users/{id}", r.partUpdateUser).Methods("PATCH")
	r.R.HandleFunc("/users/{id}", r.deleteUser).Methods("DELETE")
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK) // http_test.go
	io.WriteString(w, "Welcome to our Web-site!")
}

func (rout *Router) signUp(w http.ResponseWriter, r *http.Request) { // register
	/*
			rout.logg.Errorln("Failed to create new user: ", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	*/
	//w.WriteHeader(http.StatusOK) // http_test.go
}

func (rout *Router) signIn(w http.ResponseWriter, r *http.Request) { // Entry
	w.WriteHeader(http.StatusOK) // http_test.go
	io.WriteString(w, "SignIn was successfully")
}

func (rout *Router) logOut(w http.ResponseWriter, r *http.Request) { //logOut
	w.WriteHeader(http.StatusOK) // http_test.go
	io.WriteString(w, "You have been logout")
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

func WriteJSON(w http.ResponseWriter, status int, a interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(a)
}

func (rout *Router) CreateAccount(w http.ResponseWriter, r *http.Request) {
	user := database.CreateNewUser(1, "vlad", "12345@gmail.com", "12345678")
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		rout.logg.Errorln(err)
	}

	_, err := rout.db.CreateNewUser(user)
	if err != nil {
		rout.logg.Errorln(err)
	}

	rout.logg.Errorln(WriteJSON(w, http.StatusOK, user))
}

func GetByIdAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func UpdateAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func PartUpdateAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func DeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}
