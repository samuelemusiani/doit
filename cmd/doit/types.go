package doit

import (
	"time"
)

type TodoState struct {
	ID    int64
	State string
}

var (
	StateToDo       = TodoState{State: "todo"}
	StateInProgress = TodoState{State: "in progress"}
	StatePaused     = TodoState{State: "paused"}
	StateDone       = TodoState{State: "done"}
)

var States = []*TodoState{&StateToDo, &StateInProgress, &StatePaused, &StateDone}

type TodoPriority struct {
	ID       int64
	Priority string
}

var (
	PriorityVeryLow  = TodoPriority{Priority: "very low"}
	PriorityLow      = TodoPriority{Priority: "low"}
	PriorityMedium   = TodoPriority{Priority: "medium"}
	PriorityHigh     = TodoPriority{Priority: "high"}
	PriorityVeryHigh = TodoPriority{Priority: "very high"}
	PriorityMax      = TodoPriority{Priority: "max"}
)

var Priorities = []*TodoPriority{&PriorityVeryLow, &PriorityLow, &PriorityMedium, &PriorityHigh, &PriorityVeryHigh, &PriorityMax}

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

type Todo struct {
	ID          int64
	Title       string
	Description string
	StateID     int64
	PriorityID  int64
	ColorID     int64
	Expiration  Expiration
	UserID      int64
}

type TodoResponse struct {
	ID             int64
	Title          string
	Description    string
	State          TodoState
	Priority       TodoPriority
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

func TodoToResponse(n *Todo) *TodoResponse {
	return &TodoResponse{
		ID:          n.ID,
		Title:       n.Title,
		Description: n.Description,
	}
}

func UserUnmarshalingToUser(n *UserUnmarshaling) *User {
	var u User
	if n.ID != nil {
		u.ID = *n.ID
	}

	if n.Username != nil {
		u.Username = *n.Username
	}

	if n.Email != nil {
		u.Email = *n.Email
	}

	if n.Name != nil {
		u.Name = *n.Name
	}

	if n.Surname != nil {
		u.Surname = *n.Surname
	}

	if n.Admin != nil {
		u.Admin = *n.Admin
	}

	if n.Active != nil {
		u.Active = *n.Active
	}

	if n.Password != nil {
		u.Password = *n.Password
	}

	return &u
}
