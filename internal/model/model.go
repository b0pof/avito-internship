package model

import "time"

type Bid struct {
	ID         string    `db:"id" json:"id"`
	Name       string    `db:"name" json:"name"`
	Status     string    `db:"status" json:"status"`
	AuthorType string    `db:"author_type" json:"authorType"`
	AuthorID   string    `db:"author_id" json:"authorId"`
	Version    int       `db:"version" json:"version"`
	CreatedAt  time.Time `db:"created_at" json:"createdAt"`
}

type Tender struct {
	ID          string    `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	Status      string    `db:"status" json:"status"`
	ServiceType string    `db:"service_type" json:"serviceType"`
	Version     string    `db:"version" json:"version"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
}
