package database

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, a interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(a)
}

func (rout *Router) CreateAccount(w http.ResponseWriter, r *http.Request) {
	user := &User{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		rout.logg.Errorln(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	_, err := rout.rp.CreateNewUser(r.Context(), user)
	if err != nil {
		rout.logg.Errorln(err)
	}

	WriteJSON(w, http.StatusOK, user)
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
