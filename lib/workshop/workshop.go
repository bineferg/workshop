package workshop

import "time"

type Workshop struct {
	WorkshopID  string
	Name        string
	Description string
	Time        string
	Caption     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Cap         int
	IsFull      bool
	Cost        string
	Location    string
	Level       string
}

type Event struct {
	ID          string
	Name        string
	Description string
	Caption     string
	Time        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Cost        string
	Location    string
}

type SignUp struct {
	WorkshopID string
	FirstName  string
	LastName   string
	Email      string
	Message    string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type SignUpTable struct {
	WorkshopName string
	SignUps      []SignUp
}

func (w Workshop) New() Workshop {
	return Workshop{}
}

func (e Event) New() Event {
	return Event{}
}

func (su SignUp) New() SignUp {
	return SignUp{}
}
