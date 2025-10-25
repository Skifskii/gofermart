package repository

import "errors"

var ErrLoginAlreadyTaken = errors.New("a user with this login already exists")
var ErrUserLoginNotFound = errors.New("a user with this login not found")
var ErrWrongPassword = errors.New("wrong password")
