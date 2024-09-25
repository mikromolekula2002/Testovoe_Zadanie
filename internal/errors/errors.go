package errors

import "errors"

//Ошибки, которые мы будем выводить клиенту
var (
	ErrServer        = errors.New("на сервере произошел сбой, попробуйте позже")
	ErrMissingToken  = errors.New("у полученного guid не существует токен")
	ErrMissingUserID = errors.New("отсутствует параметр UserID")
	ErrDataToken     = errors.New("введены неверные данные, не совпадение данных с сервером")
	ErrExpiredToken  = errors.New("срок действия сессии истёк")
	ErrBlockedToken  = errors.New("ссесия заблокирована")
	ErrWrongIP       = errors.New("ip address не распознан, попробуйте еще раз")
	ErrValidToken    = errors.New("похоже на то, что Access Token был подделан")
)
