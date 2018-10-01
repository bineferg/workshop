package repository

import (
	"database/sql"
	"errors"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/workshop/lib/workshop"
)

type WorkshopDB interface {
	WorkshopByID(workshopID string) (workshop.Workshop, error)
	InsertWorkshop(workshop.Workshop) error
	GetWorkshopsAfterDate(date time.Time) ([]workshop.Workshop, error)
	GetWorkshops() ([]workshop.Workshop, error)
	GetEvents() ([]workshop.Event, error)
	UpdateWorkshop(workshop workshop.Workshop) error
	DeleteWorkshop(workshopID string) error
	GetEventsAfterDate(date time.Time) ([]workshop.Event, error)
	InsertEvent(event workshop.Event) error
	DeleteEvent(eventID string) error
	UpdateEvent(event workshop.Event) error
	EventByID(eventID string) (workshop.Event, error)
	SignUp(signup workshop.SignUp) error
	GetSignUpsByWorkshopID(workshopID string) ([]workshop.SignUp, error)
	GetNumSignUpsByWorkshopID(workshopID string) (int, error)
	GetAllSignUps() ([]workshop.SignUpTable, error)

	GetDB() interface{}
}

type workshopDB struct {
	db *sql.DB
}

func (w workshopDB) GetDB() interface{} {
	return w.db
}

func transact(db *sql.DB, txFunc func(*sql.Tx) error) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	err = txFunc(tx)
	return err
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

	err := w.db.QueryRow(`SELECT workshop_id, name, description, time, created_at, updated_at, cap, cost, location, level, caption FROM workshops WHERE id = ?`, workshopID).Scan(&ws.WorkshopID, &ws.Name, &ws.Description, &ws.Time, &ws.CreatedAt, &ws.UpdatedAt, &ws.Cap, &ws.Cost, &ws.Location, &ws.Level, &ws.Caption)

	if err != nil {
		return ws, err
	}
	return ws, nil
}

func (w workshopDB) EventByID(eventID string) (workshop.Event, error) {
	var e workshop.Event

	err := w.db.QueryRow(`SELECT event_id, name, description, time, updated_at, cost, location, caption FROM events WHERE id = ?`, eventID).Scan(&e.ID, &e.Name, &e.Description, &e.Time, &e.CreatedAt, &e.UpdatedAt, &e.Cost, &e.Location, &e.Caption)

	if err != nil {
		return e, err
	}
	return e, nil
}

func (w workshopDB) DeleteEvent(eventID string) error {
	stmt, err := w.db.Prepare("DELETE from events where event_id=?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(eventID)
	if err != nil {
		return err
	}
	return nil
}

func (w workshopDB) DeleteWorkshop(workshopID string) error {
	stmt, err := w.db.Prepare("DELETE from workshops where workshop_id=?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(workshopID)
	if err != nil {
		return err
	}
	return nil
}
func (w workshopDB) InsertWorkshop(ws workshop.Workshop) error {

	sqlCmd := "INSERT INTO workshops (workshop_id, name, description, time, created_at, updated_at, cap, cost, location, level, caption) VALUES (?,?,?,?,NOW(),NOW(),?,?,?,?,?)"

	if _, err := w.db.Exec(
		sqlCmd,
		ws.WorkshopID,
		ws.Name,
		ws.Description,
		ws.Time,
		ws.Cap,
		ws.Cost,
		ws.Location,
		ws.Level,
		ws.Caption,
	); err != nil {
		return err
	}
	return nil
}

func (w workshopDB) InsertEvent(e workshop.Event) error {

	sqlCmd := "INSERT INTO events (event_id, name, description, time, created_at, updated_at, cost, location, caption) VALUES (?,?,?,?,NOW(),NOW(),?,?,?)"

	if _, err := w.db.Exec(
		sqlCmd,
		e.ID,
		e.Name,
		e.Description,
		e.Time,
		e.Cost,
		e.Location,
		e.Caption,
	); err != nil {
		return err
	}
	return nil
}

func (w workshopDB) UpdateWorkshop(ws workshop.Workshop) error {
	sqlCmd := "UPDATE workshops SET name=?, description=?, time=?, cap=?, level=?, cost=?, location=?, caption=? WHERE workshop_id=?"

	if _, err := w.db.Exec(
		sqlCmd,
		ws.Name,
		ws.Description,
		ws.Time,
		ws.Cap,
		ws.Level,
		ws.Cost,
		ws.Location,
		ws.WorkshopID,
		ws.Caption,
	); err != nil {
		return err
	}
	return nil

}

func (w workshopDB) UpdateEvent(e workshop.Event) error {
	sqlCmd := "UPDATE events SET name=?, description=?, time=?, cost=?, location=?, caption=? WHERE event_id=?"
	if _, err := w.db.Exec(
		sqlCmd,
		e.Name,
		e.Description,
		e.Time,
		e.Cost,
		e.Location,
		e.Caption,
		e.ID,
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
		err := rows.Scan(&id, &ws.WorkshopID, &ws.Name, &ws.Description, &ws.Time, &ws.CreatedAt, &ws.UpdatedAt, &ws.Cap, &ws.Cost, &ws.Location, &ws.Level, &ws.Caption)
		if err != nil {
			return workshops, err
		}
		workshops = append(workshops, ws)

	}
	for index, _ := range workshops {
		c, err := w.GetNumSignUpsByWorkshopID(workshops[index].WorkshopID)
		if err != nil {
			return workshops, err
		}
		if workshops[index].Cap == c {
			workshops[index].IsFull = true
		}
	}
	return workshops, nil
}

func (w workshopDB) GetWorkshops() ([]workshop.Workshop, error) {
	sqlCmd := "SELECT * FROM workshops"
	var workshops []workshop.Workshop
	rows, err := w.db.Query(
		sqlCmd,
	)
	if err != nil {
		return workshops, err
	}
	var id int
	for rows.Next() {
		var ws workshop.Workshop
		err := rows.Scan(&id, &ws.WorkshopID, &ws.Name, &ws.Description, &ws.CreatedAt, &ws.UpdatedAt, &ws.Cap, &ws.Location, &ws.Level, &ws.Time, &ws.Cost, &ws.Caption)
		if err != nil {
			return workshops, err
		}
		workshops = append(workshops, ws)

	}
	for index, _ := range workshops {
		c, err := w.GetNumSignUpsByWorkshopID(workshops[index].WorkshopID)
		if err != nil {
			return workshops, err
		}
		if workshops[index].Cap == c {
			workshops[index].IsFull = true
		}
	}
	return workshops, nil
}
func (w workshopDB) GetEvents() ([]workshop.Event, error) {
	sqlCmd := "SELECT * FROM events"
	var events []workshop.Event
	rows, err := w.db.Query(
		sqlCmd,
	)
	if err != nil {
		return events, err
	}
	var id int
	for rows.Next() {
		var e workshop.Event
		err := rows.Scan(&id, &e.ID, &e.Name, &e.Description, &e.CreatedAt, &e.UpdatedAt, &e.Location, &e.Time, &e.Cost, &e.Caption)
		if err != nil {
			return events, err
		}
		events = append(events, e)

	}
	return events, nil
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
		err := rows.Scan(&id, &e.ID, &e.Name, &e.Description, &e.Time, &e.CreatedAt, &e.UpdatedAt, &e.Cost, &e.Location, &e.Caption)
		if err != nil {
			return events, err
		}
		events = append(events, e)

	}
	return events, nil
}

func (w workshopDB) SignUp(signup workshop.SignUp) error {
	log.Println(signup.WorkshopID)
	log.Println(signup.FirstName)
	log.Println(signup.LastName)
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
		if err := rows.Scan(&id, &s.WorkshopID, &s.FirstName, &s.LastName, &s.Email, &s.CreatedAt, &s.UpdatedAt, &s.Message); err != nil {
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

func (w workshopDB) GetAllSignUps() ([]workshop.SignUpTable, error) {
	var table []workshop.SignUpTable
	wNames, err := w.getWorkshopNames()
	if err != nil {
		return table, err
	}
	for _, n := range wNames {
		sups, err := w.GetSignUpsByWorkshopID(n.ID)
		if err != nil {
			return table, err
		}
		table = append(table, workshop.SignUpTable{
			WorkshopName: n.Name,
			SignUps:      sups,
		})

	}
	return table, nil

}

type workshopName struct {
	ID   string
	Name string
}

func (w workshopDB) getWorkshopNames() ([]workshopName, error) {

	sqlCmd := "SELECT workshop_id, name FROM workshops"
	var (
		id     string
		name   string
		wNames []workshopName
	)
	rows, err := w.db.Query(sqlCmd)
	if err != nil {
		return wNames, err
	}
	for rows.Next() {
		if err := rows.Scan(&id, &name); err != nil {
			return wNames, err
		}
		wNames = append(wNames, workshopName{ID: id, Name: name})
	}
	return wNames, nil

}
