package model

import (
	"time"
)

type ValueWithExpiration struct {
	Value   string
	Expires *time.Time
}
