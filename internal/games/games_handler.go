package games

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type DBHandler struct {
	Conn *pgx.Conn
}

func HandleGameRequest(w http.ResponseWriter, r *http.Request) {
	//  TODO: REFACTOR. CALL THE STEAM ENDPOINT IN OTHER MODULE
	game := r.URL.Query().Get("term")
	baseURL := "https://store.steampowered.com/api/storesearch/"
	params := url.Values{}
	params.Add("term", game)
	params.Add("cc", "ar")
	params.Add("l", "es")
	fullURL := baseURL + "?" + params.Encode()
	resp, err := http.Get(fullURL)
	if err != nil {
    	w.WriteHeader(http.StatusInternalServerError)
    	w.Write([]byte("Failed to contact Steam"))
    	return
	}
	var steamResult any
	json.NewDecoder(resp.Body).Decode(&steamResult)
	defer resp.Body.Close()
	json.NewEncoder(w).Encode(steamResult)
}

func (h *DBHandler) HandleAddGame(w http.ResponseWriter, r *http.Request) {
	gameReq := new(GameRequest)
	err := json.NewDecoder(r.Body).Decode(gameReq)
	
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Request"))
		return 
	}
	//  TODO: REFACTOR. CALL THE STEAM ENDPOINT IN OTHER MODULE (handler -> service -> repository)  
	baseURL := "https://store.steampowered.com/api/appdetails"
	params := url.Values{}
	params.Add("appids", strconv.Itoa(gameReq.SteamAppID))
	params.Add("cc", "ar")
	params.Add("l", "es")
	fullURL := baseURL + "?" + params.Encode()
	resp, err := http.Get(fullURL)
	
	if err != nil {
    	w.WriteHeader(http.StatusInternalServerError)
    	w.Write([]byte("Failed to contact Steam"))
    	return
	}
	// use map and the steam details struct to store the response of the called steam endpoint (ej: steamResp["730"] = <game_data>
	var steamResp map[string]SteamAppDetails
	json.NewDecoder(resp.Body).Decode(&steamResp)
	defer resp.Body.Close()
	
	// obtain game data via steam response
	gameDetails := steamResp[strconv.Itoa(gameReq.SteamAppID)]
	//  TODO: HANDLE SUCCESS=FALSE CASE
	currentPrice := 0.0
	if gameDetails.Data.PriceOverview != nil {
		currentPrice = float64(gameDetails.Data.PriceOverview.Final) / 100
	} else {
		currentPrice = 0.0
	}

	targetPrice := gameReq.TargetPrice // for pass the target_price into the function that stores in DB
	if targetPrice == 0 && !gameDetails.Data.IsFree && gameDetails.Data.PriceOverview != nil {
    	targetPrice = float64(gameDetails.Data.PriceOverview.Final) / 100
	}
	
	// get the user id stored in the jwt 
	claims := r.Context().Value("claims").(jwt.MapClaims)
	userIDStr, ok := claims["sub"].(string)
	if !ok {
    	w.WriteHeader(http.StatusUnauthorized)
    	w.Write([]byte("Invalid token"))
    	return
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
    	w.WriteHeader(http.StatusUnauthorized)
    	w.Write([]byte("Invalid token"))
    	return
	}
	// store in DB
	inputGame := Game{
		UserID: userID,
		SteamAppID: gameReq.SteamAppID,
		Name: gameDetails.Data.Name,
		CurrentPrice: currentPrice,
		TargetPrice: targetPrice,
		IsFree: gameDetails.Data.IsFree,
	}
	DBResponse, err := StoreGameInDB(r.Context(), h.Conn, inputGame)  
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Somethig Went Wrong In The Server"))
		log.Println(err)
		return 
	}
	json.NewEncoder(w).Encode(DBResponse)
}
