package games

import (
	"time"

	"github.com/google/uuid"
)

// for the DB
type Game struct {
	ID uuid.UUID `json:"id"` 
	UserID uuid.UUID `json:"user_id"`
	SteamAppID int `json:"steam_app_id"`
	Name string `json:"name"`
	CurrentPrice float64 `json:"current_price"`
	TargetPrice float64 `json:"target_price"`
	IsFree bool `json:"is_free"`
	Created_at time.Time `json:"created_at"`
}

// FOR THE API REQUESTS
type GameRequest struct {
	SteamAppID int `json:"steam_app_id"`
	TargetPrice float64 `json:"target_price"`
}

// FOR THE API RESPONSE 
type GameResponse struct {
	SteamAppID int `json:"steam_app_id"`
	Name string `json:"name"`
	CurrentPrice float64 `json:"current_price"`
	TargetPrice float64 `json:"target_price"`
	IsFree bool `json:"is_free"`
}
