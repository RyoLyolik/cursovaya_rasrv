package domain

type RoleOut struct {
	Name string `json:"title"`
}

type Role struct {
	ID   int
	Name string
}

var (
	Admin    = Role{ID: 0, Name: "admin"}
	Employee = Role{ID: 1, Name: "employee"}
)
