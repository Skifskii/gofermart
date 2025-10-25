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
	if _, ok := r.getUserByLogin(login); ok {
		return repository.ErrLoginAlreadyTaken
	}

	r.Users = append(r.Users, User{login, password})

	return nil
}

func (r *Repo) AuthenticateUser(login, password string) error {
	u, ok := r.getUserByLogin(login)
	if !ok {
		return repository.ErrUserLoginNotFound
	}

	if u.Password != password {
		return repository.ErrWrongPassword
	}

	return nil
}

func (r *Repo) getUserByLogin(login string) (u User, ok bool) {
	for _, user := range r.Users {
		if user.Login == login {
			return user, true
		}
	}

	return User{}, false
}
