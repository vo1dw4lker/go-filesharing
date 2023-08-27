package models

import (
	"time"
)

type File struct {
	ID         string
	FileName   string
	Expiration time.Time
	FileSize   int64
}
