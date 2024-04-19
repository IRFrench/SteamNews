package steam

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	ARTICLE_COUNT  = 5
	CONTENT_LENGTH = 10
)

type OverNews struct {
	Appnews News `json:"appnews,omitempty"`
}

type News struct {
	Appid     int       `json:"appid,omitempty"`
	Newsitems []Article `json:"newsitems,omitempty"`
	Count     int       `json:"count,omitempty"`
}

type Article struct {
	Gid           string `json:"gid,omitempty"`
	Title         string `json:"title,omitempty"`
	Url           string `json:"url,omitempty"`
	IsExternalUrl bool   `json:"is_external_url,omitempty"`
	Author        string `json:"author,omitempty"`
	Contents      string `json:"contents,omitempty"`
	Feedlabel     string `json:"feedlabel,omitempty"`
	Date          int    `json:"date,omitempty"`
	Feedname      string `json:"feedname,omitempty"`
	FeedType      int    `json:"feed_type,omitempty"`
	Appid         int    `json:"appid,omitempty"`
}

func (s *SteamClient) GetAppNews(appId int) ([]Article, error) {
	steamUrl := fmt.Sprintf(
		"http://api.steampowered.com/ISteamNews/GetNewsForApp/v0002/?appid=%d&count=%d&maxlength=%d&format=json",
		appId,
		ARTICLE_COUNT,
		CONTENT_LENGTH,
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

	var jsonResponse OverNews
	if err := json.NewDecoder(response.Body).Decode(&jsonResponse); err != nil {
		return nil, err
	}

	return jsonResponse.Appnews.Newsitems, nil
}
