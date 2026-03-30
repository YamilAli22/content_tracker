package users

import (
	"encoding/json"
	"net/http"

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
	response, err := CreateUser(r.Context(), h.Conn, newUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong in the server"))
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
