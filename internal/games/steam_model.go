package games

type SteamAppDetails struct {
    Success bool `json:"success"`
    Data *SteamAppData `json:"data"`
}

type SteamAppData struct {
    Name string `json:"name"`
    IsFree bool `json:"is_free"`
    PriceOverview *SteamPrice `json:"price_overview"`
}

type SteamPrice struct {
    Final int `json:"final"`
    Currency string `json:"currency"`
}
