package models

import "time"

// Reservation is the reservation model
type Reservation struct {
	ID int
	FirstName string
	LastName string
	Email string
	Phone string
	StartDate time.Time
	EndDate time.Time
	RoomID int
	CreatedAt time.Time
	UpdatedAt time.Time
	Room Room
}

// Restriction is the restriction model
type Restriction struct {
	ID int
	RestrictionName string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// RoomRestriction is the room restriction model
type RoomRestriction struct {
	ID int
	StartDate time.Time
	EndDate time.Time
	RoomID int
	ReservationID int
	RestrictionID int
	CreatedAt time.Time
	UpdatedAt time.Time
	Room Room
	Reservation Reservation
	Restriction Restriction
}

// Room is the room model
type Room struct {
	ID int
	RoomName string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// User is the user model
type User struct{
	ID int
	FirstName string
	LastName string
	Email string
	Password string
	AccessLevel int
	CreatedAt time.Time
	UpdatedAt time.Time
}