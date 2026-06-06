package model

import "errors"

var ErrOrderNotFound = errors.New("order not found")
var ErrPaidCanNotBeCancelled = errors.New("paid order cannot be cancelled")
var ErrOrderAlreadyPaid = errors.New("order already paid")
var ErrCancelledCanNotBePaid = errors.New("cancelled cannot be paid")

var ErrUserUUIDIsRequired = errors.New("user_uuid is required")
var ErrPartUUIDsIsRequired = errors.New("part_uuids is required")
var ErrPartNotFound = errors.New("part not found")
