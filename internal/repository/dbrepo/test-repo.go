package dbrepo

import (
	"errors"
	"time"

	"github.com/acceleraterA/go_app_udemy/internal/models"
)

// define function for postgresDBrepo
func (m *testDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation into a database and return the new reservation id
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	//if the roomid is 3, fail
	if res.RoomID == 3 {
		return 0, errors.New("some error")
	}
	return 1, nil
}

// Inserts a room restriction into the database
func (m *testDBRepo) InsertRoomRestriction(res models.RoomRestriction) error {
	if res.RoomID == 11 {
		return errors.New("some error")
	}
	return nil
}

// SearchAvalibilityByDatesByRoomID return s true if avaiability exists for roomID,
func (m *testDBRepo) SearchAvalibilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {

	return false, nil
}

// returns a slice of available rooms, if any for given date range
func (m *testDBRepo) SearchAvalibilityForAllRooms(start, end time.Time) ([]models.Room, error) {

	var rooms []models.Room

	return rooms, nil
}

// getroombyid gets a room by id and return room
func (m *testDBRepo) GetRoomByID(id int) (models.Room, error) {

	var room models.Room
	if id != 1 && id != 3 {
		return room, errors.New("some error")
	}
	return room, nil
}
