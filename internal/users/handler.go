package users

/*
handler.go for handles the http requests and serving the responses, interacting with the database if necessary
*/

import (
	"encoding/json"
	"net/http"
	
	"github.com/YamilAli22/content_tracker/internal/auth"

	"github.com/jackc/pgx/v5"
)

type Handler struct {
	Conn *pgx.Conn
}

func (h *Handler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	newUser := new(UserRequestBody)
	err := json.NewDecoder(r.Body).Decode(newUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Please enter a valid input"))
		return
	}
	newUser.Password, err = auth.HashPassword(newUser.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong in the server during the user creation"))
		return
	}
	response, err := CreateUser(r.Context(), h.Conn, newUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong in the server"))
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) HandleUserLogin(w http.ResponseWriter, r *http.Request) {
	user := new(UserRequestBody)
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Please enter a valid input"))
		return
	}
	user_check, err := GetUserByEmail(r.Context(), h.Conn, user.Email)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid credentials"))
		return
	}
	if !auth.CheckPasswordHash(user.Password, user_check.Hash) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid credentials"))
		return 
	}
	JWT, err := auth.CreateJWT(user_check.Id, user_check.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong in the server"))
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(JWT)
}




