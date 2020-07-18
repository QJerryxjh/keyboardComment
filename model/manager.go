package model

import "time"

type DbManager struct {
	ID        uint      `json:"id"`
	Account   string    `json:"account"`
	Password  string    `json:"password"`
	Name      string    `json:"name"`
	Gender    string    `json:"gender"`
	CreatedAt time.Time `json:"created_at"`
}
