package security

// Role represents the roles that app supports
type Role string

const (
	// RoleSuperAdmin indicates a super administrator with full access to all features.
	RoleSuperAdmin Role = "superadmin"
	// RoleAdmin indicates an administrator with access to manage users and content.
	RoleAdmin Role = "admin"
	// RoleEmployee indicates a regular employee with limited access to the system.
	RoleEmployee Role = "employee"
	// RoleUser indicates a regular end-user with very limited access to the system.
	RoleUser Role = "user"
)

func (r Role) String() string {
	return string(r)
}

// IsValidRole reports whether the given role is supported by the system.
func IsValidRole(role Role) bool {
	switch role {
	case RoleSuperAdmin, RoleAdmin, RoleEmployee, RoleUser:
		return true
	default:
		return false
	}
}
