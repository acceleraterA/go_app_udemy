package dbrepo

import (
	"errors"
	"log"
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

// SearchAvailabilityByDatesByRoomID return s true if avaiability exists for roomID,
func (m *testDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	// set up a test time
	layout := "2006-01-02"
	str := "2049-12-31"
	t, err := time.Parse(layout, str)
	if err != nil {
		log.Println(err)
	}

	// this is our test to fail the query -- specify 2060-01-01 as start
	testDateToFail, err := time.Parse(layout, "2060-01-01")
	if err != nil {
		log.Println(err)
	}

	if start == testDateToFail {
		return false, errors.New("some error")
	}

	// if the start date is after 2049-12-31, then return false,
	// indicating no availability;
	if start.After(t) {
		return false, nil
	}

	// otherwise, we have availability
	return true, nil
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
