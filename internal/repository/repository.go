package repository

import "errors"

var ErrLoginAlreadyTaken = errors.New("a user with this login already exists")
