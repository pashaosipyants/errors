package errors

/*
	everything that standard errors package supplies
	to allow user of this lib not to import both
*/

import (
	"errors"
)

var NewE = errors.New
var AsE = errors.As
var IsE = errors.Is
var UnwrapE = errors.Unwrap
