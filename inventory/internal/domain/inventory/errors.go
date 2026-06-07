package inventory

import "errors"

var ErrEmptyPartID = errors.New("empty part id")
var ErrPartNotFound = errors.New("part not found")
