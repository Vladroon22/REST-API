package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	db "github.com/Vladroon22/REST-API/internal/database"
	"github.com/Vladroon22/REST-API/internal/service"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Router struct {
	R    mux.Router
	logg logrus.Logger
	srv  *service.Service
}

type UserInput struct {
	Email    string `json:"email"`
	Password string `json:"pass"`
}

func NewRouter(srv *service.Service) *Router {
	return &Router{
		R:    *mux.NewRouter(),
		logg: *logrus.New(),
		srv:  srv,
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
	r.R.HandleFunc("/{id}", r.UpdateAccount).Methods("PUT")
	r.R.HandleFunc("/name/{id}", r.PartUpdateAccountName).Methods("PATCH")
	r.R.HandleFunc("/email/{id}", r.PartUpdateAccountEmail).Methods("PATCH")
	r.R.HandleFunc("/pass/{id}", r.PartUpdateAccountPass).Methods("PATCH")
	r.R.HandleFunc("/{id}", r.DeleteAccount).Methods("DELETE")
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK) // http_test.go
	io.WriteString(w, "Welcome to our Web-site!")
}

func (rout *Router) signIn(w http.ResponseWriter, r *http.Request) { // Entry
	var input UserInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		rout.logg.Errorln(err)
		return
	}

	token, err := rout.srv.Accounts.GenerateJWT(input.Email, input.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		WriteJSON(w, http.StatusUnauthorized, map[string]interface{}{"error": err.Error()})
		rout.logg.Errorln(err)
		return
	}
	if token == "" {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		WriteJSON(w, http.StatusUnauthorized, map[string]interface{}{"token": token})
		rout.logg.Errorln(err)
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"token": token,
	})
}

func WriteJSON(w http.ResponseWriter, status int, a interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(a)
}

func (rout *Router) CreateAccount(w http.ResponseWriter, r *http.Request) {
	user := &db.User{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		rout.logg.Errorln(err)
		return
	}

	id, err := rout.srv.Accounts.CreateNewUser(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		rout.logg.Errorln(err)
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (rout *Router) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		WriteJSON(w, http.StatusMethodNotAllowed, map[string]interface{}{"status": http.StatusMethodNotAllowed})
		rout.logg.Errorln(http.StatusMethodNotAllowed)
		return
	}

	user := &db.User{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		rout.logg.Errorln(err)
		return
	}

	id, err := rout.srv.Accounts.DeleteUser(r.Context(), user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		rout.logg.Errorln(err)
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (rout *Router) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		WriteJSON(w, http.StatusMethodNotAllowed, map[string]interface{}{"status": http.StatusMethodNotAllowed})
		rout.logg.Errorln(http.StatusMethodNotAllowed)
		return
	}

	user := &db.User{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		rout.logg.Errorln(err)
		return
	}

	id, err := rout.srv.Accounts.UpdateUserFully(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		rout.logg.Errorln(err)
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (rout *Router) PartUpdateAccountName(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PATCH" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		WriteJSON(w, http.StatusMethodNotAllowed, map[string]interface{}{"status": http.StatusMethodNotAllowed})
		rout.logg.Errorln(http.StatusMethodNotAllowed)
		return
	}

	user := &db.User{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		WriteJSON(w, http.StatusMethodNotAllowed, map[string]interface{}{"error": err.Error()})
		rout.logg.Errorln(err)
		return
	}

	id, err := rout.srv.Accounts.PartUpdateUserName(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		rout.logg.Errorln(err)
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (rout *Router) PartUpdateAccountEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PATCH" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		WriteJSON(w, http.StatusMethodNotAllowed, map[string]interface{}{"status": http.StatusMethodNotAllowed})
		rout.logg.Errorln(http.StatusMethodNotAllowed)
		return
	}

	user := &db.User{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		rout.logg.Errorln(err)
		return
	}

	id, err := rout.srv.Accounts.PartUpdateUserEmail(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		rout.logg.Errorln(err)
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (rout *Router) PartUpdateAccountPass(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PATCH" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		rout.logg.Errorln(http.StatusMethodNotAllowed)
		return
	}

	user := &db.User{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		rout.logg.Errorln(err)
		return
	}

	id, err := rout.srv.Accounts.PartUpdateUserPass(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		rout.logg.Errorln(err)
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"id": id,
	})
}
