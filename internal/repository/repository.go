package repository

import (
	"time"

	"github.com/acceleraterA/go_app_udemy/internal/models"
)

type DatabaseRepo interface {
	AllUsers() bool
	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(res models.RoomRestriction) error
	SearchAvalibilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error)
	SearchAvalibilityForAllRooms(start, end time.Time) ([]models.Room, error)
	GetRoomByID(id int) (models.Room, error)
}
