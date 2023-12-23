package userstorage

type UserRepository interface {
	Save(login, authHash string) error
	Find(login string) (string, error)
}
