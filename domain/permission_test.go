package domain

import (
	"testing"
)

func TestPermission_Has(t *testing.T) {
	testCases := []struct {
		name           string
		userPerms      Permission
		requiredPerm   Permission
		expectedResult bool
	}{
		{
			name:           "Admin has PermUserManage",
			userPerms:      RoleAdmin,
			requiredPerm:   PermUserManage,
			expectedResult: true,
		},
		{
			name:           "Admin has PermCreateAny",
			userPerms:      RoleAdmin,
			requiredPerm:   PermCreateAny,
			expectedResult: true,
		},
		{
			name:           "Editor does not have PermUserManage",
			userPerms:      RoleEditor,
			requiredPerm:   PermUserManage,
			expectedResult: false,
		},
		{
			name:           "Regular user has PermCreatePrefix",
			userPerms:      RoleRegular,
			requiredPerm:   PermCreatePrefix,
			expectedResult: true,
		},
		{
			name:           "Regular user does not have PermCreateAny",
			userPerms:      RoleRegular,
			requiredPerm:   PermCreateAny,
			expectedResult: false,
		},
		{
			name:           "Guest has no permissions",
			userPerms:      RoleGuest,
			requiredPerm:   PermCreatePrefix,
			expectedResult: false,
		},
		{
			name:           "Checking for zero permission should be false",
			userPerms:      RoleAdmin,
			requiredPerm:   0,
			expectedResult: false,
		},
		{
			name:           "User with single permission has it",
			userPerms:      PermDeleteOwn,
			requiredPerm:   PermDeleteOwn,
			expectedResult: true,
		},
		{
			name:           "User with single permission does not have another",
			userPerms:      PermDeleteOwn,
			requiredPerm:   PermDeleteAny,
			expectedResult: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.userPerms.Has(tc.requiredPerm); got != tc.expectedResult {
				t.Errorf("Permission.Has() = %v, want %v", got, tc.expectedResult)
			}
		})
	}
}

func TestPermission_Add(t *testing.T) {
	testCases := []struct {
		name           string
		initialPerms   Permission
		permToAdd      Permission
		expectedPerms  Permission
	}{
		{
			name:           "Add PermCreateAny to Regular user",
			initialPerms:   RoleRegular,
			permToAdd:      PermCreateAny,
			expectedPerms:  RolePrivileged,
		},
		{
			name:           "Add existing permission does not change anything",
			initialPerms:   RoleRegular,
			permToAdd:      PermCreatePrefix,
			expectedPerms:  RoleRegular,
		},
		{
			name:           "Add permission to Guest",
			initialPerms:   RoleGuest,
			permToAdd:      PermCreatePrefix,
			expectedPerms:  PermCreatePrefix,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.initialPerms.Add(tc.permToAdd); got != tc.expectedPerms {
				t.Errorf("Permission.Add() = %v, want %v", got, tc.expectedPerms)
			}
		})
	}
}

func TestPermission_Remove(t *testing.T) {
	testCases := []struct {
		name           string
		initialPerms   Permission
		permToRemove   Permission
		expectedPerms  Permission
	}{
		{
			name:           "Remove PermCreateAny from Privileged user",
			initialPerms:   RolePrivileged,
			permToRemove:   PermCreateAny,
			expectedPerms:  RoleRegular,
		},
		{
			name:           "Remove non-existing permission does not change anything",
			initialPerms:   RoleRegular,
			permToRemove:   PermCreateAny,
			expectedPerms:  RoleRegular,
		},
		{
			name:           "Remove permission from Admin",
			initialPerms:   RoleAdmin,
			permToRemove:   PermUserManage,
			expectedPerms:  RoleEditor,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.initialPerms.Remove(tc.permToRemove); got != tc.expectedPerms {
				t.Errorf("Permission.Remove() = %v, want %v", got, tc.expectedPerms)
			}
		})
	}
}
