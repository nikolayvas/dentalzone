package exceptions

import "errors"

// ErrAlreadyExists is used in case of duplications or unexpected existence
var ErrAlreadyExists = errors.New("")

// ErrNotSuch is used in case of non existence of something expected
var ErrNotSuch = errors.New("")
