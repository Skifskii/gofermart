package inmem

import "gophermart/internal/repository"

type User struct {
	Login     string
	Password  string
	Balance   float64
	Withdrawn float64
	OrderNums []string
}

type Repo struct {
	Users map[string]User
}

func New() *Repo {
	r := Repo{}
	r.Users = make(map[string]User)

	return &r
}

func (r *Repo) AddUser(login, password string) error {
	if _, ok := r.Users[login]; ok {
		return repository.ErrLoginAlreadyTaken
	}

	r.Users[login] = User{Login: login, Password: password}

	return nil
}

func (r *Repo) AuthenticateUser(login, password string) error {
	u, ok := r.Users[login]
	if !ok {
		return repository.ErrUserLoginNotFound
	}

	if u.Password != password {
		return repository.ErrWrongPassword
	}

	return nil
}

func (r *Repo) GetBalance(login string) (current, withdrawn float64, err error) {
	u, ok := r.Users[login]
	if !ok {
		return 0, 0, repository.ErrUserLoginNotFound
	}
	return u.Balance, u.Withdrawn, err
}

func (r *Repo) GetOrderNums(login string) (orderNums []string, err error) {
	u, ok := r.Users[login]
	if !ok {
		return nil, repository.ErrUserLoginNotFound
	}
	return u.OrderNums, nil
}

func (r *Repo) AddOrderNum(login, orderNum string) error {
	u, ok := r.Users[login]
	if !ok {
		return repository.ErrUserLoginNotFound
	}

	u.OrderNums = append(u.OrderNums, orderNum)
	r.Users[u.Login] = u

	return nil
}
