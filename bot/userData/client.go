package userdata

type StoreClient struct {
	users                map[int]User
	userInfoRequestChan  chan userInfoRequest
	userInforesponseChan chan User

	// Update steam info
	steamRequestChan  chan userSteamRequest
	steamResponseChan chan bool
}

type userInfoRequest struct {
	name string
	id   int
}

type userSteamRequest struct {
	id      int
	steamId *int
	game    *steamGame
}

type steamGame struct {
	add  bool
	name string
	id   int
}

func (s *StoreClient) Run() {
	for {
		select {
		case userRequest := <-s.userInfoRequestChan:
			user, ok := s.users[userRequest.id]
			if ok {
				s.userInforesponseChan <- user
				continue
			}

			createdUser := newUser(userRequest.name)
			s.users[userRequest.id] = createdUser
			s.userInforesponseChan <- createdUser

		case steamRequest := <-s.steamRequestChan:
			user, ok := s.users[steamRequest.id]
			if !ok {
				s.steamResponseChan <- false
			}

			if steamRequest.steamId != nil {
				user.Steam.Id = *steamRequest.steamId
			}

			if steamRequest.game != nil {
				if steamRequest.game.add {
					user.Steam.Added[steamRequest.game.id] = steamRequest.game.name
					delete(user.Steam.Removed, steamRequest.game.id)
				}
			}
		}
	}
}

func NewStoreClient() StoreClient {
	return StoreClient{
		users:                make(map[int]User, 1),
		userInfoRequestChan:  make(chan userInfoRequest, 5),
		userInforesponseChan: make(chan User, 5),
		steamRequestChan:     make(chan userSteamRequest, 5),
		steamResponseChan:    make(chan bool, 5),
	}
}
