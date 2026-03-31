package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/jackc/pgx/v5"

	"github.com/YamilAli22/content_tracker/internal/db"
	"github.com/YamilAli22/content_tracker/internal/users"
)

var PORT string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Couldn't load .env file!!")
	}
	PORT = os.Getenv("PORT")
}

type Server struct {
	Router *chi.Mux
	DB *pgx.Conn
}

func newServer(db *pgx.Conn) *Server {
	server := &Server{
		Router: chi.NewRouter(), // NewRouter() returns a new Mux
		DB: db,	
	}

	usersHandler := &users.Handler{Conn: server.DB}
	server.Router.Post("/users", usersHandler.HandleCreateUser)
	
	return server
}

func main() {
	conn, err := db.Connect()
	if err != nil {
		log.Fatalln(err)
	}

	defer conn.Close(context.Background())

	err = conn.Ping(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Successfully Connected to Postgresql")

	server := newServer(conn)

	fmt.Println("server running on port", PORT)
	http.ListenAndServe(":"+PORT, server.Router)
}
