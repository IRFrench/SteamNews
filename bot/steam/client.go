package steam

import (
	"net/http"
	"time"
)

type SteamClient struct {
	client *http.Client
	user   User
}

func (s *SteamClient) sendRequest(req *http.Request) (*http.Response, error) {
	return s.client.Do(req)
}

func NewClient(user User) SteamClient {
	return SteamClient{
		&http.Client{
			Timeout: 5 * time.Second,
		},
		user,
	}
}
