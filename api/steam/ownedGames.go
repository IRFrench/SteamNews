package steam

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type OverResponse struct {
	Response Response `json:"response,omitempty"`
}

type Response struct {
	GameCount int    `json:"game_count,omitempty"`
	Games     []Game `json:"games,omitempty"`
}

type Game struct {
	Appid                  int    `json:"appid,omitempty"`
	Name                   string `json:"name,omitempty"`
	PlaytimeForever        int    `json:"playtime_forever,omitempty"`
	PlaytimeWindowsForever int    `json:"playtime_windows_forever,omitempty"`
	PlaytimeMacForever     int    `json:"playtime_mac_forever,omitempty"`
	PlaytimeDeckForever    int    `json:"playtime_deck_forever,omitempty"`
	RtimeLastPlayed        int    `json:"rtime_last_played,omitempty"`
	PlaytimeDisconnected   int    `json:"playtime_disconnected,omitempty"`
}

func (s *SteamClient) GetOwnedGames(steamId int) ([]Game, error) {
	steamUrl := fmt.Sprintf(
		"http://api.steampowered.com/IPlayerService/GetOwnedGames/v0001/?key=%s&steamid=%d&format=json&include_appinfo=true",
		s.key,
		steamId,
	)

	request, err := http.NewRequest(http.MethodGet, steamUrl, nil)
	if err != nil {
		return nil, err
	}

	response, err := s.sendRequest(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("non 200 response recieved from API")
	}

	var jsonResponse OverResponse
	if err := json.NewDecoder(response.Body).Decode(&jsonResponse); err != nil {
		return nil, err
	}

	return jsonResponse.Response.Games, nil
}
