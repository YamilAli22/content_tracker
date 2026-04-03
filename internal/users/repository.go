package users

/* 
repository.go is responsible for managing database operations, implementing necessary operations, and interacting with the database.
*/

import (
	"github.com/jackc/pgx/v5"
	"github.com/google/uuid"
	"context"
)

func CreateUser(ctx context.Context, conn *pgx.Conn, user *UserRequestBody) (UserResponse, error) {
	var id uuid.UUID
	var email string
	query := `INSERT INTO users (email, hash) VALUES ($1, $2) RETURNING id, email`
	err := conn.QueryRow(ctx, query, user.Email, user.Password).Scan(&id, &email)
	if err != nil {
		response := UserResponse{
			Id: id,
		}
		return response, err 
	}
	response := UserResponse{
		Id: id,
		Email: email,
	}
	return response, err
}


// esta funcion devuelve los usuarios traidos de la db en un formato json 
func GetUsers(ctx context.Context, conn *pgx.Conn) ([]*UserResponse, error) {
    rows, _ := conn.Query(ctx, "SELECT id, email FROM users")
    users, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (*UserResponse, error) {
        user := new(UserResponse)
        err := row.Scan(&user.Id, &user.Email)
        return user, err
    })
    return users, err
}

// GetUserByEmail returns User (internal use, contains hash)
func GetUserByEmail(ctx context.Context, conn *pgx.Conn, email string) (User, error) {
    var user User
    query := `SELECT id, email, hash FROM users WHERE email = $1`
    err := conn.QueryRow(ctx, query, email).Scan(&user.Id, &user.Email, &user.Hash)
	if err != nil {
		return User{}, err
	}
    return user, err
}


