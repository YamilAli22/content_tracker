package users

/* 
repository.go is responsible for managing database operations, implementing necessary operations, and interacting with the database.
*/

import (
	"github.com/jackc/pgx/v5"
	"context"
)

func CreateUser(ctx context.Context, conn *pgx.Conn, user *UserRequestBody) (UserResponse, error) {
	var id int
	query := `INSERT INTO users (email, hash) VALUES ($1, $2) RETURNING id`
	err := conn.QueryRow(ctx, query, user.Email, user.Password).Scan(&id)
	if err != nil {
		response := UserResponse{
			Id: id,
		}
		return response, err 
	}
	response := UserResponse{
		Id: id,
	}
	return response, err
}


// esta funcion devuelve los usuarios traidos de la db en un formato json 
func GetUsers(ctx context.Context, conn *pgx.Conn) ([]*UserResponse, error) {
    rows, _ := conn.Query(ctx, "SELECT id, email, hash FROM users")
    users, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (*UserResponse, error) {
        user := new(UserResponse)
        err := row.Scan(&user.Id, &user.Email)
        return user, err
    })
    return users, err
}


