package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	_ "github.com/Vladroon22/REST-API/docs"
	db "github.com/Vladroon22/REST-API/internal/database"
	"github.com/Vladroon22/REST-API/internal/service"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	TTL = time.Hour
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

// SayHello godoc
// @Summary Say Hello
// @ID create-account
// @Accept  json
// @Produce  json
// @Description Main page
// @Tags example
// @Success 200 {string} string "MAIN PAGE"
// @Router / [get]

func (r *Router) SayHello() {
	r.R.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("MAIN PAGE")) }).Methods("GET")
}

// Swagger godoc
// @Summary Swagger documentation
// @Description Swagger documentation endpoint
// @Tags swagger
// @Router /swagger/ [get]

func (r *Router) Swagger() {
	r.R.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://127.0.0.1:8000/swagger/doc.json"), //The url pointing to API definition
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods("GET")
}

func (r *Router) AuthEndPoints() {
	r.R.HandleFunc("/sign-up", r.CreateAccount).Methods("POST")
	r.R.HandleFunc("/sign-in", r.signIn).Methods("POST")

}

func (r *Router) UserEndPoints(sub *mux.Router) {
	sub.HandleFunc("/{id:[0-9]+}", r.GetAccount).Methods("GET")
	sub.HandleFunc("/{id:[0-9]+}", r.UpdateAccount).Methods("PUT, POST")
	sub.HandleFunc("/name/{id:[0-9]+}", r.PartUpdateAccountName).Methods("POST, PATCH")
	sub.HandleFunc("/email/{id:[0-9]+}", r.PartUpdateAccountEmail).Methods("POST, PATCH")
	sub.HandleFunc("/pass/{id:[0-9]+}", r.PartUpdateAccountPass).Methods("POST, PATCH")
	sub.HandleFunc("/{id:[0-9]+}", r.DeleteAccount).Methods("DELETE")
	sub.HandleFunc("/logout/{id:[0-9]+}", r.logout).Methods("GET")
}

func SetCookie(w http.ResponseWriter, cookieName string, cookies string) {
	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    cookies,
		Path:     "/",
		Secure:   false,
		HttpOnly: true,
		Expires:  time.Now().Add(TTL),
	}
	http.SetCookie(w, cookie)
}

func ClearCookie(w http.ResponseWriter, cookieName string, cookies string) {
	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    cookies,
		Path:     "/",
		Expires:  time.Unix(0, 0),
		Secure:   false,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
}

func (rout *Router) AuthMiddleWare(next http.Handler) http.Handler {
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

// logout godoc
// @Summary Log out a user
// @Description Log out a user from the app and clear cookies
// @Tags auth
// @Accept json
// @Produce json
// @Success 303 {string} string "See Other"
// @Failure 500 {string} string "Internal server error"
// @Router /logout [get]

func (rout *Router) logout(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, _ := strconv.Atoi(vars["id"])
	rout.logg.Infof("Has been logout the ID: %d", userID)

	ClearCookie(w, "jwt", "")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// signIn godoc
// @Summary Sign in a user
// @Description Log in a user by email and password
// @Tags auth
// @Accept  json
// @Produce  json
// @Param input body AuthInput true "Credentials"
// @Success 200 {string} string "token"
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {string} string "Internal server error"
// @Router /sign-in [post]

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

// CreateAccount godoc
// @Summary Create a new account
// @Description Register a new user
// @Tags auth
// @Accept  json
// @Produce  json
// @Param user body RegInput true "Registration data"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {string} string "Internal server error"
// @Router /sign-up [post]

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

// DeleteAccount godoc
// @Summary Delete an account
// @Description Remove a user account by ID
// @Tags user
// @Param id path int true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {string} string "Internal server error"
// @Router /user/{id} [delete]

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

// GetAccount godoc
// @Summary Get an account
// @Description Get user account by ID
// @Tags user
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {string} string "Internal server error"
// @Router /user/{id} [get]

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

// UpdateAccount godoc
// @Summary Update an account
// @Description Fully update a user account
// @Tags user
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Param user body db.User true "User data"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {string} string "Internal server error"
// @Router /user/{id} [put]

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

// PartUpdateAccountName godoc
// @Summary Partially update account name
// @Description Update user account name
// @Tags user
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body db.User true "User data"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {string} "Internal server error"
// @Router /user/name/{id} [patch]

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

// PartUpdateAccountEmail godoc
// @Summary Partially update account email
// @Description Update user account email
// @Tags user
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body db.User true "User data"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {string} "Internal server error"
// @Router /user/email/{id} [patch]

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

// PartUpdateAccountPass godoc
// @Summary Partially update account password
// @Description Update user account password
// @Tags user
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body db.User true "User data"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {string} "Internal server error"
// @Router /user/pass/{id} [patch]

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
