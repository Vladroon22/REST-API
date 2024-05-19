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
	r.R.HandleFunc("/{id}/name", r.PartUpdateAccountName).Methods("PATCH")
	r.R.HandleFunc("/{id}/email", r.PartUpdateAccountEmail).Methods("PATCH")
	r.R.HandleFunc("/{id}/pass", r.PartUpdateAccountPass).Methods("PATCH")
	r.R.HandleFunc("/{id}", r.DeleteAccount).Methods("DELETE")
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK) // http_test.go
	io.WriteString(w, "Welcome to our Web-site!")
}

func (rout *Router) signIn(w http.ResponseWriter, r *http.Request) { // Entry
	w.WriteHeader(http.StatusOK) // http_test.go
	io.WriteString(w, "SignIn was successfully")
}

func WriteJSON(w http.ResponseWriter, status int, a interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(a)
}

func (rout *Router) CreateAccount(w http.ResponseWriter, r *http.Request) {
	user := &db.User{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		rout.logg.Errorln(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	id, err := rout.srv.CreateNewUser(r.Context(), user)
	if err != nil {
		rout.logg.Errorln(err)
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (rout *Router) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		rout.logg.Errorln(http.StatusMethodNotAllowed)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	user := &db.User{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		rout.logg.Errorln(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	id, err := rout.srv.DeleteUser(r.Context(), user.ID)
	if err != nil {
		rout.logg.Errorln(err)
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (rout *Router) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		rout.logg.Errorln(http.StatusMethodNotAllowed)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	user := &db.User{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		rout.logg.Errorln(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	id, err := rout.srv.UpdateUserFully(r.Context(), user)
	if err != nil {
		rout.logg.Errorln(err)
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (rout *Router) PartUpdateAccountName(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PATCH" {
		rout.logg.Errorln(http.StatusMethodNotAllowed)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	user := &db.User{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		rout.logg.Errorln(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	id, err := rout.srv.PartUpdateUserName(r.Context(), user)
	if err != nil {
		rout.logg.Errorln(err)
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (rout *Router) PartUpdateAccountEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PATCH" {
		rout.logg.Errorln(http.StatusMethodNotAllowed)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	user := &db.User{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		rout.logg.Errorln(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	id, err := rout.srv.PartUpdateUserEmail(r.Context(), user)
	if err != nil {
		rout.logg.Errorln(err)
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (rout *Router) PartUpdateAccountPass(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PATCH" {
		rout.logg.Errorln(http.StatusMethodNotAllowed)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	user := &db.User{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		rout.logg.Errorln(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	id, err := rout.srv.PartUpdateUserPass(r.Context(), user)
	if err != nil {
		rout.logg.Errorln(err)
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"id": id,
	})
}