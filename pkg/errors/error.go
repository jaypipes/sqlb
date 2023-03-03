//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package errors

import (
	"errors"
)

var (
	InvalidJoinNoSelect      = errors.New("Unable to join selection. There was no selection to join to.")
	InvalidJoinUnknownTarget = errors.New("Unable to join selection. Target selection was not found.")
	NoTargetTable            = errors.New("No target table supplied.")
	NoValues                 = errors.New("No values supplied.")
	UnknownColumn            = errors.New("Received an unknown column.")
)
