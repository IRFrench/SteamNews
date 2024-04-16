package userdata

type User struct {
	Name  string
	Steam SteamInfo
}

type SteamInfo struct {
	Id       int
	UseOwned bool
	Added    map[int]string
	Removed  map[int]string
}

func newUser(name string) User {
	return User{
		Name: name,
		Steam: SteamInfo{
			Added:   make(map[int]string, 1),
			Removed: make(map[int]string, 1),
		},
	}
}
