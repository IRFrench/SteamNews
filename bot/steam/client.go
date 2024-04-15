package steam

import (
	"net/http"
	"time"
)

type SteamClient struct {
	client *http.Client
	key    string
}

func (s *SteamClient) sendRequest(req *http.Request) (*http.Response, error) {
	return s.client.Do(req)
}

func NewClient(key string) SteamClient {
	return SteamClient{
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		key: key,
	}
}
