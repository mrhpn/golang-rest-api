package security

// UserStatus represents the status of users
type UserStatus string

const (
	// UserStatusActive indicates a user is active & currently using the system
	UserStatusActive UserStatus = "active"
	// UserStatusInactive indicates a user is inactive & currently not using the system
	UserStatusInactive UserStatus = "inactive"
	// UserStatusBlocked indicates a user gets blocked by admin/superadmin and cannot access the system
	UserStatusBlocked UserStatus = "blocked"
)

func (r UserStatus) String() string {
	return string(r)
}

// IsValidUserStatus reports whether the given status is supported by the system.
func IsValidUserStatus(status UserStatus) bool {
	switch status {
	case UserStatusActive, UserStatusInactive, UserStatusBlocked:
		return true
	default:
		return false
	}
}
