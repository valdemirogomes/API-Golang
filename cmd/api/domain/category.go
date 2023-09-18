package domain

import "time"

type Category struct {
	ID        int64     `json:"id,omitempty"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}
