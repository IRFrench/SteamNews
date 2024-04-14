package steam

type User struct {
	key string
	id  int
}

func NewUser(key string, id int) User {
	return User{
		key: key,
		id:  id,
	}
}
