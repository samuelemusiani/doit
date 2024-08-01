package doit

type Note struct {
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
	External bool
	Active   bool
	Password string
}
