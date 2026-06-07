package inventory

import "errors"

var (
	ErrEmptyPartID  = errors.New("empty part id")
	ErrPartNotFound = errors.New("part not found")
)
