package dbrepo

import (
	"context"
	"log"
	"time"

	"github.com/acceleraterA/go_app_udemy/internal/models"
)

// define function for postgresDBrepo
func (m *postgresDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation into a database
func (m *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {
	//close the transaction after the 5 minutes lifetime if nothing is happening
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var newID int

	//sql query
	stmt := `insert into reservations (first_name,last_name,email,phone,start_date,end_date, room_id, created_at,updated_at) values($1,$2,$3,$4,$5,$6,$7,$8,$9) returning id`

	err := m.DB.QueryRowContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	).Scan(&newID)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return newID, nil
}

// Inserts a room restriction into the database
func (m *postgresDBRepo) InsertRoomRestriction(res models.RoomRestriction) error {
	//close the transaction after the 5 minutes lifetime if nothing is happening
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `insert into room_restrictions (start_date,end_date, room_id, reservation_id, restriction_id, created_at,updated_at) values($1,$2,$3,$4,$5,$6,$7)`

	_, err := m.DB.ExecContext(ctx, stmt,

		res.StartDate,
		res.EndDate,
		res.RoomID,
		res.ReservationID,
		res.RestrictionID,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// SearchAvalibilityByDatesByRoomID return s true if avaiability exists for roomID,
func (m *postgresDBRepo) SearchAvalibilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	//close the transaction after the 5 minutes lifetime if nothing is happening
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `
		SELECT
			count(id) 
		FROM 
			room_restrictions 
		WHERE 
			'$1 <end_date and $2 >start_date and room_id=$3;`
	var numRows int
	row := m.DB.QueryRowContext(ctx, query,
		start, end, roomID)
	err := row.Scan(&numRows)
	if err != nil {
		log.Println(err)
		return false, err
	}
	return numRows == 0, nil
}
