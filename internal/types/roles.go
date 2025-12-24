package types

type Role string

const (
	RoleSuperAdmin Role = "superadmin"
	RoleAdmin      Role = "admin"
	RoleEmployee   Role = "employee"
	RoleUser       Role = "user"
)

func (r Role) String() string {
	return string(r)
}

var ValidRoles = map[Role]bool{
	RoleSuperAdmin: true,
	RoleAdmin:      true,
	RoleEmployee:   true,
	RoleUser:       true,
}
