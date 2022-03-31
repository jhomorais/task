package entities

type Role string

const (
	ROLE_TECHNICIAN Role = "TECHNICIAN"
	ROLE_MANAGER    Role = "MANAGER"
)

type User struct {
	ID       string `gorm:"id"`
	Email    string `gorm:"index"`
	Password string
	Role     Role
}

func (r Role) String() string {
	return string(r)
}
