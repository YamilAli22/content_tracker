package games

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func StoreGameInDB(ctx context.Context, conn *pgx.Conn, input Game) (GameResponse, error) {
	var SteamAppId int 
	var Name string
	var CurrentPrice float64
	var TargetPrice float64
	var IsFree bool


	query := `INSERT INTO games (user_id, steam_app_id, name, current_price, taget_price, is_free) 
			  VALUES ($1, $2, $3, $4, $5, $6) 
			  RETURNING steam_app_id, name, current_price, target_price, is_free`
	
	err := conn.QueryRow(ctx, query, input.UserID, input.SteamAppID, input.Name, input.CurrentPrice, input.TargetPrice, input.IsFree).Scan(&SteamAppId, &Name, &CurrentPrice, &TargetPrice, &IsFree)
	if err != nil {
		response := GameResponse{}
		return response, err
	}
	response := GameResponse{
		SteamAppID: SteamAppId,
		Name: Name,
		CurrentPrice: CurrentPrice,
		TargetPrice: TargetPrice,
		IsFree: IsFree,
	}
	return response, nil
}
