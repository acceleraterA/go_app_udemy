package dbrepo

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/acceleraterA/go_app_udemy/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// define function for postgresDBrepo
func (m *postgresDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation into a database and return the new reservation id
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

// SearchAvailabilityByDatesByRoomID return s true if avaiability exists for roomID,
func (m *postgresDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	//close the transaction after the 5 minutes lifetime if nothing is happening
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `
		SELECT
			count(id) 
		FROM 
			room_restrictions 
		WHERE 
			$1 <end_date and $2 >start_date and room_id=$3;`
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

// returns a slice of available rooms, if any for given date range
func (m *postgresDBRepo) SearchAvalibilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	//close the transaction after the 5 minutes lifetime if nothing is happening
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var rooms []models.Room
	query :=
		`SELECT
    	r.id, r.room_name
	FROM
    	rooms r
	where
	r.id not in (select rr.room_id from room_restrictions rr where $1 <rr.end_date and $2 >rr.start_date)`

	rows, err := m.DB.QueryContext(ctx, query, start, end)

	if err != nil {
		log.Println(err)
		return rooms, err
	}
	defer rows.Close()
	for rows.Next() {
		var room models.Room
		err := rows.Scan(&room.ID, &room.RoomName)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}
	//check err again when done with scanning
	if err = rows.Err(); err != nil {
		return rooms, err
	}
	return rooms, nil
}

// getroombyid gets a room by id and return room
func (m *postgresDBRepo) GetRoomByID(id int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var room models.Room
	query := `
select id,room_name,created_at,updated_at
from rooms 
where id=$1`
	row := m.DB.QueryRowContext(ctx, query,
		id)
	err := row.Scan(&room.ID, &room.RoomName, &room.CreatedAt, &room.UpdatedAt)
	if err != nil {
		log.Println(err)
		return room, err
	}
	return room, nil
}

// GetUserByID gets a user by id and return user
func (m *postgresDBRepo) GetUserByID(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var user models.User
	query := `select id,first_name, last_name,email,password, access_level, created_at, updated_at
	 from users where id=$1`
	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.AccessLevel, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		log.Println(err)
		return user, err
	}
	return user, nil
}

func (m *postgresDBRepo) UpdateUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `update users set first_name=$1,last_name=$2,email=$3,access_level=$4,updated_at=$5 
	`
	_, err := m.DB.ExecContext(ctx, query, u.FirstName, u.LastName, u.Email, u.AccessLevel, time.Now())
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// Authenticate a user
func (m *postgresDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var id int
	var hashedPassword string

	row := m.DB.QueryRowContext(ctx, "select id,password from users where email=$1", email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		log.Println(err)
		return id, "", err
	} else {
		//check if the testpassword match with the one in database, cased with byte
		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
		if err == bcrypt.ErrMismatchedHashAndPassword {
			//_,err=m.DB.ExecContext(ctx,"update users set password=$1",)
			return 0, "", errors.New("incorrect password")
		} else if err != nil {
			return 0, "", err
		}
		return id, hashedPassword, nil
	}
}

// returns a slice of all reservations
func (m *postgresDBRepo) GetAllReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var reservations []models.Reservation
	query :=
		`SELECT
    	r.id, r.first_name,r.last_name,r.email,phone,r.start_date,r.end_date, r.room_id, r.created_at,r.updated_at,r.processed,
		rm.id,rm.room_name
	FROM
    	reservations r
	left join rooms rm on (r.room_id = rm.id)
	order by r.start_date asc`

	rows, err := m.DB.QueryContext(ctx, query)

	if err != nil {
		log.Println(err)
		return reservations, err
	}
	//prevent memory leak
	defer rows.Close()

	for rows.Next() {
		var res models.Reservation
		err := rows.Scan(
			&res.ID,
			&res.FirstName,
			&res.LastName,
			&res.Email,
			&res.Phone,
			&res.StartDate,
			&res.EndDate,
			&res.RoomID,
			&res.CreatedAt,
			&res.UpdatedAt,
			&res.Processed,
			&res.Room.ID,
			&res.Room.RoomName,
		)
		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, res)
	}
	//check err again when done with scanning
	if err = rows.Err(); err != nil {
		return reservations, err
	}
	return reservations, nil
}

// returns a slice of new reservations
func (m *postgresDBRepo) GetNewReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var reservations []models.Reservation
	query :=
		`SELECT
    	r.id, r.first_name,r.last_name,r.email,phone,r.start_date,r.end_date, r.room_id, r.created_at,r.updated_at,r.processed,
		rm.id,rm.room_name
	FROM
    	reservations r
	left join rooms rm on (r.room_id = rm.id)
	WHERE
		r.processed=0
	order by r.start_date asc`

	rows, err := m.DB.QueryContext(ctx, query)

	if err != nil {
		log.Println(err)
		return reservations, err
	}
	//prevent memory leak
	defer rows.Close()

	for rows.Next() {
		var res models.Reservation
		err := rows.Scan(
			&res.ID,
			&res.FirstName,
			&res.LastName,
			&res.Email,
			&res.Phone,
			&res.StartDate,
			&res.EndDate,
			&res.RoomID,
			&res.CreatedAt,
			&res.UpdatedAt,
			&res.Processed,
			&res.Room.ID,
			&res.Room.RoomName,
		)
		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, res)
	}
	//check err again when done with scanning
	if err = rows.Err(); err != nil {
		return reservations, err
	}
	return reservations, nil
}

// returns one reservation by ID
func (m *postgresDBRepo) GetReservationByID(id int) (models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var res models.Reservation
	query := `
		SELECT
			r.id, r.first_name,r.last_name,r.email,phone,r.start_date,r.end_date, r.room_id, r.created_at,r.updated_at,r.processed,
			rm.id,rm.room_name
		FROM 
			reservations r
			left join rooms rm on (r.room_id = rm.id) 
		WHERE 
			r.id=$1`
	row := m.DB.QueryRowContext(ctx, query,
		id)
	err := row.Scan(
		&res.ID,
		&res.FirstName,
		&res.LastName,
		&res.Email,
		&res.Phone,
		&res.StartDate,
		&res.EndDate,
		&res.RoomID,
		&res.CreatedAt,
		&res.UpdatedAt,
		&res.Processed,
		&res.Room.ID,
		&res.Room.RoomName,
	)
	if err != nil {
		log.Println(err)
		return res, err
	}
	return res, nil
}

func (m *postgresDBRepo) UpdateReservation(u models.Reservation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `update reservations set first_name=$1,last_name=$2,email=$3,phone=$4,updated_at=$5
	where id=$6 
	`
	_, err := m.DB.ExecContext(ctx, query, u.FirstName, u.LastName, u.Email, u.Phone, time.Now(), u.ID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (m *postgresDBRepo) DeleteReservation(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `delete from reservations 
	where id=$1 
	`
	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
func (m *postgresDBRepo) UpdateProcessedForReservation(id int, processed int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `update reservations set processed=$1
	where id=$2
	`
	_, err := m.DB.ExecContext(ctx, query, processed, id)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (m *postgresDBRepo) AllRooms() ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var rooms []models.Room
	query := `select id,room_name,created_at,updated_at from rooms order by room_name
	`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return rooms, err
	}
	defer rows.Close()

	for rows.Next() {
		var r models.Room
		err := rows.Scan(
			&r.ID,
			&r.RoomName,
			&r.CreatedAt,
			&r.UpdatedAt,
		)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, r)
	}
	if err = rows.Err(); err != nil {
		return rooms, err
	}
	return rooms, nil
}

func (m *postgresDBRepo) GetRestrictionsForRoomByDate(roomID int, start, end time.Time) ([]models.RoomRestriction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var restrictions []models.RoomRestriction
	query := `
	select id,coalesce(reservation_id,0), restriction_id,room_id,start_date,end_date 
	from room_restrictions 
	where $1<end_date AND $2 >start_date AND room_id=$3
	`
	//coalesce() is has value, use it, if not use 0
	rows, err := m.DB.QueryContext(ctx, query, start, end, roomID)
	if err != nil {
		return restrictions, err
	}
	defer rows.Close()

	for rows.Next() {
		var r models.RoomRestriction
		err := rows.Scan(
			&r.ID,
			&r.ReservationID,
			&r.RestrictionID,
			&r.RoomID,
			&r.StartDate,
			&r.EndDate,
		)
		if err != nil {
			return restrictions, err
		}
		restrictions = append(restrictions, r)
	}
	if err = rows.Err(); err != nil {
		return restrictions, err
	}
	return restrictions, nil
}

// insert a room restriction for block
func (m *postgresDBRepo) InsertBlockForRoom(id int, startDate time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `insert into room_restrictions (start_date, end_date, room_id, restriction_id, created_at,updated_at) values ($1,$2,$3,$4,$5,$6)`
	_, err := m.DB.ExecContext(ctx, query, startDate, startDate.AddDate(0, 0, 1), id, 2, time.Now(), time.Now())
	if err != nil {
		log.Println(err)
		return err
	}
	return nil

}

// delete a room restriction for block
func (m *postgresDBRepo) DeleteBlockByID(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `delete from room_restrictions where id=$1`
	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil

}
