package repository

import (
	"database/sql"
	"errors"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/workshop/lib/workshop"
)

type WorkshopDB interface {
	WorkshopByID(workshopID string) (workshop.Workshop, error)
	InsertWorkshop(workshop.Workshop) error
	GetWorkshopsAfterDate(date time.Time) ([]workshop.Workshop, error)
	UpdateWorkshop(workshop workshop.Workshop) error
	GetEventsAfterDate(date time.Time) ([]workshop.Event, error)
	InsertEvent(event workshop.Event) error
	UpdateEvent(event workshop.Event) error
	EventByID(eventID string) (workshop.Event, error)
	SignUp(signup workshop.SignUp) error
	GetSignUpsByWorkshopID(workshopID string) ([]workshop.SignUp, error)
	GetNumSignUpsByWorkshopID(workshopID string) (int, error)

	GetDB() interface{}
}

type workshopDB struct {
	db *sql.DB
}

func (w workshopDB) GetDB() interface{} {
	return w.db
}

func NewWorkshopDB(dns string) (*workshopDB, error) {
	if dns == "" {
		return nil, errors.New("db dns not found")
	}

	db, err := sql.Open("mysql", dns)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("DESCRIBE workshops")
	if err != nil {
		return nil, err
	}

	rows.Close()

	return &workshopDB{db: db}, nil
}

func (w workshopDB) WorkshopByID(workshopID string) (workshop.Workshop, error) {
	var ws workshop.Workshop

	err := w.db.QueryRow(`SELECT workshop_id, name, description, start_time, end_time, created_at, updated_at, cap, cost, location, level FROM workshops WHERE id = ?`, workshopID).Scan(&ws.WorkshopID, &ws.Name, &ws.Description, &ws.StartTime, &ws.EndTime, &ws.CreatedAt, &ws.UpdatedAt, &ws.Cap, &ws.Cost, &ws.Location, &ws.Level)

	if err != nil {
		return ws, err
	}
	return ws, nil
}

func (w workshopDB) EventByID(eventID string) (workshop.Event, error) {
	var e workshop.Event

	err := w.db.QueryRow(`SELECT event_id, name, description, start_time, created_at, updated_at, cost, location FROM events WHERE id = ?`, eventID).Scan(&e.ID, &e.Name, &e.Description, &e.StartTime, &e.CreatedAt, &e.UpdatedAt, &e.Cost, &e.Location)

	if err != nil {
		return e, err
	}
	return e, nil
}

func (w workshopDB) InsertWorkshop(ws workshop.Workshop) error {

	sqlCmd := "INSERT INTO workshops (workshop_id, name, description, start_time, end_time, created_at, updated_at, cap, cost, location, level) VALUES (?,?,?,?,?,NOW(),NOW(),?,?,?,?)"

	if _, err := w.db.Exec(
		sqlCmd,
		ws.WorkshopID,
		ws.Name,
		ws.Description,
		ws.StartTime,
		ws.EndTime,
		ws.Cap,
		ws.Cost,
		ws.Location,
		ws.Level,
	); err != nil {
		return err
	}
	return nil
}

func (w workshopDB) InsertEvent(e workshop.Event) error {

	sqlCmd := "INSERT INTO events (event_id, name, description, start_time, created_at, updated_at, cost, location) VALUES (?,?,?,?,NOW(),NOW(),?,?)"

	if _, err := w.db.Exec(
		sqlCmd,
		e.ID,
		e.Name,
		e.Description,
		e.StartTime,
		e.Cost,
		e.Location,
	); err != nil {
		return err
	}
	return nil
}

func (w workshopDB) UpdateWorkshop(ws workshop.Workshop) error {

	sqlCmd := "INSERT INTO workshops (workshop_id, name, description, start_time, end_time, created_at, updated_at, cap, cost, location, level) VALUES (?,?,?,?,?,NOW(),NOW(),?,?,?,?) ON DUPLICATE KEY UPDATE name=VALUES(name), description=VALUES(description), start_time=VALUES(start_time), end_time=VALUES(end_time), updated_at=VALUES(NOW()), cap=VALUES(cap), cost=VALUES(cost), location=VALUES(location), level=VALUES(level)"

	if _, err := w.db.Exec(
		sqlCmd,
		ws.WorkshopID,
		ws.Name,
		ws.Description,
		ws.StartTime,
		ws.EndTime,
		ws.Cap,
		ws.Cost,
		ws.Location,
		ws.Level,
	); err != nil {
		return err
	}
	return nil

}

func (w workshopDB) UpdateEvent(e workshop.Event) error {

	sqlCmd := "INSERT INTO events (event_id, name, description, start_time, created_at, updated_at, cost, location) VALUES (?,?,?,?,NOW(),NOW(),?,?) ON DUPLICATE KEY UPDATE name=VALUES(name), description=VALUES(description), start_time=VALUES(start_time), updated_at=VALUES(NOW()), cost=VALUES(cost), location=VALUES(location)"

	if _, err := w.db.Exec(
		sqlCmd,
		e.ID,
		e.Name,
		e.Description,
		e.StartTime,
		e.Cost,
		e.Location,
	); err != nil {
		return err
	}
	return nil

}
func (w workshopDB) GetWorkshopsAfterDate(date time.Time) ([]workshop.Workshop, error) {
	sqlCmd := "SELECT * FROM workshops where start_time > ?"
	var workshops []workshop.Workshop
	rows, err := w.db.Query(
		sqlCmd,
		date,
	)
	if err != nil {
		return workshops, err
	}
	var id int
	for rows.Next() {
		var ws workshop.Workshop
		err := rows.Scan(&id, &ws.WorkshopID, &ws.Name, &ws.Description, &ws.StartTime, &ws.EndTime, &ws.CreatedAt, &ws.UpdatedAt, &ws.Cap, &ws.Cost, &ws.Location, &ws.Level)
		if err != nil {
			return workshops, err
		}
		workshops = append(workshops, ws)

	}
	for _, wk := range workshops {
		c, err := w.GetNumSignUpsByWorkshopID(wk.WorkshopID)
		if err != nil {
			return workshops, err
		}
		if wk.Cap == c {
			wk.IsFull = true
		}
	}
	return workshops, nil
}

func (w workshopDB) GetEventsAfterDate(date time.Time) ([]workshop.Event, error) {
	sqlCmd := "SELECT * FROM events where start_time > ?"
	var events []workshop.Event
	rows, err := w.db.Query(
		sqlCmd,
		date,
	)
	if err != nil {
		return events, err
	}
	var id int
	for rows.Next() {
		var e workshop.Event
		err := rows.Scan(&id, &e.ID, &e.Name, &e.Description, &e.StartTime, &e.CreatedAt, &e.UpdatedAt, &e.Cost, &e.Location)
		if err != nil {
			return events, err
		}
		events = append(events, e)

	}
	return events, nil
}

func (w workshopDB) SignUp(signup workshop.SignUp) error {
	sqlCmd := "INSERT INTO signups (workshop_id, first_name, last_name, email, created_at, updated_at) VALUES(?, ?, ?, ?, NOW(), NOW())"
	if _, err := w.db.Exec(
		sqlCmd,
		signup.WorkshopID,
		signup.FirstName,
		signup.LastName,
		signup.Email,
	); err != nil {
		return err
	}
	return nil

}

func (w workshopDB) GetSignUpsByWorkshopID(workshopID string) ([]workshop.SignUp, error) {
	sqlCmd := "SELECT * FROM signups WHERE workshop_id = ?"
	var signups []workshop.SignUp
	rows, err := w.db.Query(sqlCmd, workshopID)
	if err != nil {
		return signups, err
	}
	var id int
	for rows.Next() {
		var s workshop.SignUp
		if err := rows.Scan(&id, &s.WorkshopID, &s.FirstName, &s.LastName, &s.Email, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return signups, err
		}
		signups = append(signups, s)

	}
	return signups, nil
}
func (w workshopDB) GetNumSignUpsByWorkshopID(workshopID string) (int, error) {
	sqlCmd := "SELECT count(*) FROM signups WHERE workshop_id = ?"
	var count int
	rows, err := w.db.Query(sqlCmd, workshopID)
	if err != nil {
		return 0, err
	}
	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return 0, err

		}

	}
	return count, nil
}
