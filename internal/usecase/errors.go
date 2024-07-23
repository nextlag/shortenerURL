package usecase

import "errors"

// ErrConflict - ошибка конфликта данных
var ErrConflict = errors.New("data conflict in DBStorage")
