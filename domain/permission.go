package domain

// Permission is a bitmask type for user permissions.
// Creating random URLs is a baseline action available to guests and all authenticated users,
// and thus is not governed by a specific permission bit.
type Permission uint

const (
	PermCreatePrefix Permission = 1 << iota // 1
	PermCreateAny                           // 2
	PermDeleteOwn                           // 4
	PermDeleteAny                           // 8
	PermViewOwnStats                        // 16
	PermViewAnyStats                        // 32
	PermUserManage                          // 64
)

// Permission sets (pre-configured permission bundles)
const (
	RoleGuest      Permission = 0
	RoleRegular    Permission = PermCreatePrefix | PermDeleteOwn | PermViewOwnStats
	RolePrivileged Permission = RoleRegular | PermCreateAny
	RoleEditor     Permission = RolePrivileged | PermDeleteAny | PermViewAnyStats
	RoleAdmin      Permission = RoleEditor | PermUserManage
)

// Has checks if the user's permissions (p) include the required permission (required).
func (p Permission) Has(required Permission) bool {
	// If required is 0, it's an invalid permission to check.
	if required == 0 {
		return false
	}
	return (p & required) == required
}

// Add grants a new permission.
func (p Permission) Add(perm Permission) Permission {
	return p | perm
}

// Remove revokes a permission.
func (p Permission) Remove(perm Permission) Permission {
	return p &^ perm
}
