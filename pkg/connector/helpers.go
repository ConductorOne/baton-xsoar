package connector

import (
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-xsoar/pkg/xsoar"
)

const ResourcesPageSize = 50

func annotationsForUserResourceType() annotations.Annotations {
	annos := annotations.Annotations{}
	annos.Update(&v2.SkipEntitlementsAndGrants{})
	return annos
}

func flattenRoleNames(data map[string][]string) []string {
	var roles []string

	for _, values := range data {
		roles = append(roles, values...)
	}

	return roles
}

func containsRole(roles []string, role string) bool {
	for _, r := range roles {
		if r == role {
			return true
		}
	}

	return false
}

func findUser(users []xsoar.User, id string) *xsoar.User {
	for _, user := range users {
		if user.Id == id {
			return &user
		}
	}

	return nil
}

func removeRole(roles []string, targetRole string) []string {
	var newRoles []string

	for _, role := range roles {
		if role == targetRole {
			continue
		}

		newRoles = append(newRoles, role)
	}

	return newRoles
}

func containsTargetUser(users []string, targetUser string) bool {
	for _, user := range users {
		if user == targetUser {
			return true
		}
	}

	return false
}

func removeUsers(users []xsoar.User, targetUsers ...string) []xsoar.User {
	var newUsers []xsoar.User

	for _, user := range users {
		if containsTargetUser(targetUsers, user.Id) {
			continue
		}

		newUsers = append(newUsers, user)
	}

	return newUsers
}
