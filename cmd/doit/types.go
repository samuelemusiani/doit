package doit

type Note struct {
	ID          int64
	Title       string
	Description string
	UserID      int64
}

type NoteResponse struct {
	ID          int64
	Title       string
	Description string
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
