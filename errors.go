package queueing

import (
	"errors"
)

var (

	// ErrModel is returned when a queueing model is not supported. See
	// constants.
	//
	ErrModel = errors.New("queueing: model is not supported")

	// ErrUtilisation is returned when the arrival rate equals or exceeds the
	// service rate.
	//
	ErrUtilisation = errors.New("queueing: system is fully utilised")
)
