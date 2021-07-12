package enums

type RoleIndex int

const (
	_ RoleIndex = iota
	User
	Admin
)

func (r RoleIndex) String() string {
	return [...]string{"User", "Admin"}[r-1]
}
