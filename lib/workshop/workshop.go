package workshop

import "time"

type Workshop struct {
	ID          string
	Name        string
	Description string
	StartTime   time.Time
	EndTime     time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Cap         int
	Cost        float64
	Location    string
	Level       string
}

func (w Workshop) New() Workshop {
	return Workshop{}
}
