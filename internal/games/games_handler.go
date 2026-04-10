package games

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"
)

type DBHandler struct {
	Conn *pgx.Conn
}

func HandleGameRequest(w http.ResponseWriter, r *http.Request) {
	game := r.URL.Query().Get("term")
	steamResp, err := CallSteamSearch(r.Context(), game)
	if err != nil {
		log.Printf("CallSteamSerch: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong"))
		return
	}
	json.NewEncoder(w).Encode(steamResp)
}

func (h *DBHandler) HandleAddGame(w http.ResponseWriter, r *http.Request) {
	gameReq := new(GameRequest)
	err := json.NewDecoder(r.Body).Decode(gameReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Request"))
		return 
	}

	steamResp, err := CallSteamSearchByID(r.Context(), gameReq.SteamAppID)
	if err != nil {
    	log.Printf("CallSteamSearchByID: %v", err)
    	w.WriteHeader(http.StatusInternalServerError)
    	w.Write([]byte("Something went wrong"))
    	return
	}

	gameDetails, err := GameDetailsLogic(steamResp, gameReq.SteamAppID, gameReq.TargetPrice)
	if err != nil {
 	   	log.Printf("CallSteamEndpointByID: %v", err)
    	w.WriteHeader(http.StatusInternalServerError)
    	w.Write([]byte("Something went wrong"))
   	 	return
	}
	
	userID, err := GetTokenClaims(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return 
	}

	// store in DB
	inputGame := Game{
		UserID: userID,
		SteamAppID: gameReq.SteamAppID,
		Name: gameDetails.Name,
		CurrentPrice: gameDetails.CurrentPrice,
		TargetPrice: gameDetails.TargetPrice,
		IsFree: gameDetails.IsFree,
	}
	log.Println("holi")
	DBResponse, err := StoreGameInDB(r.Context(), h.Conn, inputGame)  
	log.Printf("DB RESPONSE: %v\n Error: %v", DBResponse, err)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Somethig Went Wrong In The Server"))
		return 
	}

	json.NewEncoder(w).Encode(DBResponse)
}
