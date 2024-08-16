package doit

import (
	"time"
)

type NoteState struct {
	ID    int64
	State string
}

var (
	StateDone       = NoteState{State: "done"}
	StateToDo       = NoteState{State: "todo"}
	StateInProgress = NoteState{State: "in progress"}
	StatePaused     = NoteState{State: "paused"}
)

var States = []*NoteState{&StateDone, &StateToDo, &StateInProgress, &StatePaused}

type NotePriority struct {
	ID       int64
	Priority uint8
}

var (
	PriorityVeryLow  = NotePriority{Priority: 0}
	PriorityLow      = NotePriority{Priority: 50}
	PriorityMedium   = NotePriority{Priority: 100}
	PriorityHigh     = NotePriority{Priority: 150}
	PriorityVeryHigh = NotePriority{Priority: 200}
	PriorityMax      = NotePriority{Priority: 250}
)

var Priorities = []*NotePriority{&PriorityVeryLow, &PriorityLow, &PriorityMedium, &PriorityHigh, &PriorityVeryHigh, &PriorityMax}

// Hex of the color
type Color struct {
	ID  int64
	Hex string
}

var (
	ColorBlack = Color{Hex: "#000000"}
	ColorRed   = Color{Hex: "#ff0000"}
	ColorBlue  = Color{Hex: "#0000ff"}
	ColorGreen = Color{Hex: "#00ff00"}
	ColorWhite = Color{Hex: "#ffffff"}
)

var Colors = []*Color{&ColorBlack, &ColorRed, &ColorBlue, &ColorWhite}

type Note struct {
	ID             int64
	Title          string
	Description    string
	StateID        int64
	PriorityID     int64
	ColorID        int64
	ExpirationDate time.Time
	UserID         int64
}

type NoteResponse struct {
	ID             int64
	Title          string
	Description    string
	State          NoteState
	Priority       NotePriority
	Color          Color
	ExpirationDate time.Time
}

type User struct {
	ID       int64
	Username string
	Email    string
	Name     string
	Surname  string
	Admin    bool
	Active   bool
	Password string
}

// This is used during JSON unmarshaling to check if values are present
type UserUnmarshaling struct {
	ID       *int64
	Username *string
	Email    *string
	Name     *string
	Surname  *string
	Admin    *bool
	Active   *bool
	Password *string
}

type UserResponse struct {
	ID       int64
	Username string
	Email    string
	Name     string
	Surname  string
	Admin    bool
	External bool
	Active   bool
}

func UserToResponse(u *User) *UserResponse {
	return &UserResponse{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		Name:     u.Name,
		Surname:  u.Surname,
		Admin:    u.Admin,
		Active:   u.Active,
	}
}

func NoteToResponse(n *Note) *NoteResponse {
	return &NoteResponse{
		ID:          n.ID,
		Title:       n.Title,
		Description: n.Description,
	}
}
