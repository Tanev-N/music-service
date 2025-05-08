package models

import (
	"github.com/google/uuid"
)

type Genre struct {
	ID   uuid.UUID
	Name string
}

type TrackGenre struct {
	TrackID uuid.UUID
	GenreID uuid.UUID
}
