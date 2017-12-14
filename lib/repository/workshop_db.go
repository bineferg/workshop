package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/workshop/lib/workshop"
	_"github.com/go-sql-driver/mysql"
)

type WorkshopDB interface {
	WorkshopByID(workshopID string) (workshop.Workshop, error)
	InsertWorkshop(workshop.Workshop) error
	GetWorkshopsAfterDate(date time.Time) ([]workshop.Workshop, error)
	UpdateWorkshop(workshop workshop.Workshop) error

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

	err := w.db.QueryRow(`SELECT workshop_id, name, description, start_time, end_time, created_at, updated_at, cap, cost, location, level WHERE id = ?`, workshopID).Scan(&ws.ID, &ws.Name, &ws.Description, &ws.StartTime, &ws.EndTime, &ws.CreatedAt, &ws.UpdatedAt, &ws.Cap, &ws.Cost, &ws.Location, &ws.Level)

	if err != nil {
		return ws, err
	}
	return ws, nil
}

func (w workshopDB) InsertWorkshop(ws workshop.Workshop) error {

	sqlCmd := "INSERT INTO workshops (workshop_id, name, description, start_time, end_time, created_at, updated_at, cap, cost, location, level) VALUES (?,?,?,?,?,NOW(),NOW(),?,?,?,?)"

	if _, err := w.db.Exec(
		sqlCmd,
		ws.ID,
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

func (w workshopDB) UpdateWorkshop(ws workshop.Workshop) error {

	sqlCmd := "INSERT INTO workshops (workshop_id, name, description, start_time, end_time, created_at, updated_at, cap, cost, location, level) VALUES (?,?,?,?,?,NOW(),NOW(),?,?,?,?) ON DUPLICATE KEY UPDATE name=VALUES(name), description=VALUES(description), start_time=VALUES(start_time), end_time=VALUES(end_time), updated_at=VALUES(NOW()), cap=VALUES(cap), cost=VALUES(cost), location=VALUES(location), level=VALUES(level)"

	if _, err := w.db.Exec(
		sqlCmd,
		ws.ID,
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
		err := rows.Scan(&id, &ws.ID, &ws.Name, &ws.Description, &ws.StartTime, &ws.EndTime, &ws.CreatedAt, &ws.UpdatedAt, &ws.Cap, &ws.Cost, &ws.Location, &ws.Level)
		if err != nil {
			return workshops, err
		}
		workshops = append(workshops, ws)

	}
	return workshops, nil
}
