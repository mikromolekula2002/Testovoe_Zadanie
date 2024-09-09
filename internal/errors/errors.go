package errors

import "errors"

//Ошибки, которые мы будем выводить клиенту
var (
	ErrServer        = errors.New("на сервере произошел сбой, попробуйте позже")
	ErrMissingToken  = errors.New("у полученного guid не существует токен")
	ErrMissingUserID = errors.New("отсутствует параметр UserID")
	ErrDataToken     = errors.New("введены неверные данные, не совпадение данных с сервером")
)
