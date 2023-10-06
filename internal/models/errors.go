package models

import "fmt"

var (
	ErrorNotFoundName      = fmt.Errorf("not found name")
	ErrorAlreadyExistsName = fmt.Errorf("already exists name")
)
