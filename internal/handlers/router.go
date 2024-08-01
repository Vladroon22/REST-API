package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	db "github.com/Vladroon22/REST-API/internal/database"
	"github.com/Vladroon22/REST-API/internal/service"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Router struct {
	R    *mux.Router
	logg *logrus.Logger
	srv  *service.Service
}

type AuthInput struct {
	Email    string `json:"email"`
	Password string `json:"pass"`
}

type RegInput struct {
	Name     string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"pass"`
}

func NewRouter(srv *service.Service) *Router {
	return &Router{
		R:    mux.NewRouter(),
		logg: logrus.New(),
		srv:  srv,
	}
}

func (r *Router) SayHello() {
	r.R.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("MAIN PAGE")) }).Methods("GET")
}

func (r *Router) AuthEndPoints() {
	r.R.HandleFunc("/sign-up", r.CreateAccount).Methods("POST")
	r.R.HandleFunc("/sign-in", r.signIn).Methods("POST")
}

func (r *Router) UserEndPoints(sub *mux.Router) {
	sub.HandleFunc("/{id:[0-9]+}", r.GetAccount).Methods("GET")
	sub.HandleFunc("/{id:[0-9]+}", r.UpdateAccount).Methods("PUT")
	sub.HandleFunc("/name/{id:[0-9]+}", r.PartUpdateAccountName).Methods("PATCH")
	sub.HandleFunc("/email/{id:[0-9]+}", r.PartUpdateAccountEmail).Methods("PATCH")
	sub.HandleFunc("/pass/{id:[0-9]+}", r.PartUpdateAccountPass).Methods("PATCH")
	sub.HandleFunc("/{id:[0-9]+}", r.DeleteAccount).Methods("DELETE")
}

func SetCookie(w http.ResponseWriter, cookieName string, cookies string) {
	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    cookies,
		Path:     "/",
		Secure:   false,
		HttpOnly: true,
		Expires:  time.Now().Add(time.Hour),
	}
	http.SetCookie(w, cookie)
}

func AuthMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("jwt")
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if cookie.Value == "" {
			http.Error(w, "Cookie is empty", http.StatusUnauthorized)
			return
		}
		claims, err := db.ValidateToken(cookie.Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), "id", claims.UserId))
		next.ServeHTTP(w, r)
	})
}

func (rout *Router) signIn(w http.ResponseWriter, r *http.Request) { // Entry
	var input AuthInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		rout.logg.Errorln(err)
		return
	}

	token, err := rout.srv.Accounts.GenerateJWT(r.Context(), input.Email, input.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		rout.logg.Errorln(err)
		return
	}
	if token == "" {
		http.Error(w, "token is empty", http.StatusUnauthorized)
		rout.logg.Errorln(err)
		return
	}

	SetCookie(w, "jwt", token)
}

func WriteJSON(w http.ResponseWriter, status int, a interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(a)
}

func (rout *Router) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var user RegInput
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		rout.logg.Errorln(err)
		return
	}

	id, err := rout.srv.Accounts.CreateNewUser(r.Context(), user.Name, user.Email, user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		rout.logg.Errorln(err)
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (rout *Router) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, _ := strconv.Atoi(vars["id"])
	rout.logg.Infof("ID: %d", userID)

	id, err := rout.srv.Accounts.DeleteUser(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		rout.logg.Errorln(err)
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (rout *Router) GetAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID, _ := strconv.Atoi(vars["id"])
	rout.logg.Infof("ID: %d", ID)

	user, err := rout.srv.Accounts.GetUser(r.Context(), ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		rout.logg.Errorln(err)
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	})
}

func (rout *Router) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID, _ := strconv.Atoi(vars["id"])
	rout.logg.Infof("ID: %d", ID)

	user := &db.User{}
	user.ID = ID

	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		rout.logg.Errorln(err)
		return
	}

	id, err := rout.srv.Accounts.UpdateUserFully(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		rout.logg.Errorln(err)
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (rout *Router) PartUpdateAccountName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID, _ := strconv.Atoi(vars["id"])
	rout.logg.Infof("ID: %d", ID)

	user := &db.User{}
	user.ID = ID

	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		rout.logg.Errorln(err)
		return
	}

	id, err := rout.srv.Accounts.PartUpdateUserName(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		rout.logg.Errorln(err)
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (rout *Router) PartUpdateAccountEmail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID, _ := strconv.Atoi(vars["id"])
	rout.logg.Infof("ID: %d", ID)

	user := &db.User{}
	user.ID = ID

	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		rout.logg.Errorln(err)
		return
	}

	id, err := rout.srv.Accounts.PartUpdateUserEmail(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		rout.logg.Errorln(err)
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (rout *Router) PartUpdateAccountPass(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID, _ := strconv.Atoi(vars["id"])
	rout.logg.Infof("ID: %d", ID)

	user := &db.User{}
	user.ID = ID

	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		rout.logg.Errorln(err)
		return
	}

	id, err := rout.srv.Accounts.PartUpdateUserPass(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		rout.logg.Errorln(err)
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"id": id,
	})
}
