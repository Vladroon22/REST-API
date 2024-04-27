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
}

func NewRouter() Router {
	return Router{
		R:    mux.NewRouter(),
		logg: logrus.New(),
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
	r.R.HandleFunc("/sign-up", r.signUp).Methods("GET")
	r.R.HandleFunc("/sign-in", r.signIn).Methods("GET")
	r.R.HandleFunc("/logout", r.logOut).Methods("GET")
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

func (rout Router) signUp(w http.ResponseWriter, r *http.Request) { // register
	if err := CreateAccount(w, r); err != nil {
		rout.logg.Errorln("Failed to create new user: ", err)
		http.Error(w, "Error: ", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK) // http_test.go
	io.WriteString(w, "You have registered")
}

func (rout Router) signIn(w http.ResponseWriter, r *http.Request) { // Entry
	w.WriteHeader(http.StatusOK) // http_test.goa
	io.WriteString(w, "SignIn was successfully")
}

func (rout Router) logOut(w http.ResponseWriter, r *http.Request) { //logOut
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

func WriteJSON(w http.ResponseWriter, status int, a any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(a)
}

func CreateAccount(w http.ResponseWriter, r *http.Request) error {
	input := &database.User{
		ID:       1,
		Name:     "Vlad",
		Email:    "12345@gmail.com",
		Password: "12345678",
	}
	if err := json.NewDecoder(r.Body).Decode(input); err != nil {
		return err
	}
	_, err := database.CreateNewUser(input)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, input)
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
