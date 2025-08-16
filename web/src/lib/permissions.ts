// Corresponds to the permission system in the backend (domain/permission.go)

export enum Permission {
	PermNone = 0, // 0
	PermCreatePrefix = 1 << 0, // 1
	PermCreateAny = 1 << 1, // 2
	PermDeleteOwn = 1 << 2, // 4
	PermDeleteAny = 1 << 3, // 8
	PermViewOwnStats = 1 << 4, // 16
	PermViewAnyStats = 1 << 5, // 32
	PermUserManage = 1 << 6, // 64
}

export const permissionNames: { [key in Permission]?: string } = {
	[Permission.PermCreatePrefix]: "Create Custom Prefix",
	[Permission.PermCreateAny]: "Create for Any User",
	[Permission.PermDeleteOwn]: "Delete Own Links",
	[Permission.PermDeleteAny]: "Delete Any Link",
	[Permission.PermViewOwnStats]: "View Own Stats",
	[Permission.PermViewAnyStats]: "View Any Stats",
	[Permission.PermUserManage]: "Manage Users",
  };


/**
 * Checks if a user's permission set includes the required permission.
 * @param userPermission The user's permission bitmask.
 * @param requiredPermission The permission to check for.
 * @returns True if the user has the permission, false otherwise.
 */
export function hasPermission(userPermission: number, requiredPermission: Permission): boolean {
	if (requiredPermission === 0) {
		return false
	}
	return (userPermission & requiredPermission) === requiredPermission
}

// --- Specific Permission Checkers ---

/**
 * Checks if the user can create short URLs with a custom prefix.
 */
export function canCreatePrefix(userPermission: number): boolean {
	return hasPermission(userPermission, Permission.PermCreatePrefix)
}

/**
 * Checks if the user can create short URLs for any user.
 */
export function canCreateAny(userPermission: number): boolean {
	return hasPermission(userPermission, Permission.PermCreateAny)
}

/**
 * Checks if the user can delete their own short URLs.
 */
export function canDeleteOwn(userPermission: number): boolean {
	return hasPermission(userPermission, Permission.PermDeleteOwn)
}

/**
 * Checks if the user can delete any short URL.
 */
export function canDeleteAny(userPermission: number): boolean {
	return hasPermission(userPermission, Permission.PermDeleteAny)
}

/**
 * Checks if the user can view stats for their own short URLs.
 */
export function canViewOwnStats(userPermission: number): boolean {
	return hasPermission(userPermission, Permission.PermViewOwnStats)
}

/**
 * Checks if the user can view stats for any short URL.
 */
export function canViewAnyStats(userPermission: number): boolean {
	return hasPermission(userPermission, Permission.PermViewAnyStats)
}

/**
 * Checks if the user can manage other users.
 */
export function canManageUsers(userPermission: number): boolean {
	return hasPermission(userPermission, Permission.PermUserManage)
}