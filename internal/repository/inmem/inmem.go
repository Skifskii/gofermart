package inmem

import "gophermart/internal/repository"

type User struct {
	Login    string
	Password string
}

type Repo struct {
	Users []User
}

func New() *Repo {
	r := Repo{}
	r.Users = make([]User, 0, 100)

	return &r
}

func (r *Repo) AddUser(login, password string) error {
	if r.userExists(login) {
		return repository.ErrLoginAlreadyTaken
	}

	r.Users = append(r.Users, User{login, password})

	return nil
}

func (r *Repo) userExists(login string) bool {
	for _, user := range r.Users {
		if user.Login == login {
			return true
		}
	}

	return false
}
