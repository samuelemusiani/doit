package doit

import (
	"time"
)

type NoteState struct {
	ID    int64
	State string
}

var (
	StateToDo       = NoteState{State: "todo"}
	StateInProgress = NoteState{State: "in progress"}
	StatePaused     = NoteState{State: "paused"}
	StateDone       = NoteState{State: "done"}
)

var States = []*NoteState{&StateToDo, &StateInProgress, &StatePaused, &StateDone}

type NotePriority struct {
	ID       int64
	Priority string
}

var (
	PriorityVeryLow  = NotePriority{Priority: "very low"}
	PriorityLow      = NotePriority{Priority: "low"}
	PriorityMedium   = NotePriority{Priority: "medium"}
	PriorityHigh     = NotePriority{Priority: "high"}
	PriorityVeryHigh = NotePriority{Priority: "very high"}
	PriorityMax      = NotePriority{Priority: "max"}
)

var Priorities = []*NotePriority{&PriorityVeryLow, &PriorityLow, &PriorityMedium, &PriorityHigh, &PriorityVeryHigh, &PriorityMax}

// Hex of the color
type Color struct {
	ID  int64
	Hex string
}

var (
	ColorGreen  = Color{Hex: "#b0c5a4"}
	ColorRed    = Color{Hex: "#d37676"}
	ColorOrange = Color{Hex: "#f6995c"}
	ColorBlue   = Color{Hex: "#51829b"}
	ColorYellow = Color{Hex: "#ffc96f"}
	ColorWhite  = Color{Hex: "#ffffff"}
)

var Colors = []*Color{&ColorGreen, &ColorRed, &ColorOrange, &ColorBlue, &ColorYellow, &ColorWhite}

type Expiration struct {
	DoesExpire bool
	Date       time.Time
}

type Note struct {
	ID          int64
	Title       string
	Description string
	StateID     int64
	PriorityID  int64
	ColorID     int64
	Expiration  Expiration
	UserID      int64
}

type NoteResponse struct {
	ID             int64
	Title          string
	Description    string
	State          NoteState
	Priority       NotePriority
	Color          Color
	ExpirationDate Expiration
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
