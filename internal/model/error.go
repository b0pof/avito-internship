package model

import "github.com/pkg/errors"

var (
	ErrInvalidQueryParam = errors.New("невалидное значение query-параметра")
	ErrInvalidBody       = errors.New("невалидное тело запроса")
	ErrInvalidPathParam  = errors.New("невалидное значение path-параметра")
)

var (
	ErrNoOrganizationFound = errors.New("пользователь не является ответственным ни в одной организации")
)

var (
	ErrUserNotFound = errors.New("пользователь не найден")
)

var (
	ErrNoRights      = errors.New("доступ ограничен")
	ErrWrongDecision = errors.New("неверное значение решения")
)

var (
	ErrTenderNotFound        = errors.New("тендер не найден")
	ErrInvalidAttributeValue = errors.New("невалидное значение атрибута")
	ErrNoSuchVersion         = errors.New("версия не найдена")
)

var (
	ErrNoBidsFound = errors.New("предложений не найдено")
	ErrNoBidFound  = errors.New("предложение не найдено")
	ErrInternal    = errors.New("внутренняя ошибка")
)
