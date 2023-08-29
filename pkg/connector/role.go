package connector

import (
	"context"
	"fmt"
	"strings"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	ent "github.com/conductorone/baton-sdk/pkg/types/entitlement"
	"github.com/conductorone/baton-sdk/pkg/types/grant"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
	"github.com/conductorone/baton-xsoar/pkg/xsoar"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
)

const roleMember = "member"
const defaultAdminUser = "admin"

type roleResourceType struct {
	resourceType *v2.ResourceType
	client       *xsoar.Client
}

func (r *roleResourceType) ResourceType(_ context.Context) *v2.ResourceType {
	return r.resourceType
}

// roleResource creates a new connector resource for a Xsoar Role.
func roleResource(ctx context.Context, role *xsoar.Role) (*v2.Resource, error) {
	rolePermissionsString := strings.Join(role.Permissions, ",")
	profile := map[string]interface{}{
		"role_id":          role.Id,
		"role_name":        role.Name,
		"role_permissions": rolePermissionsString,
	}

	resource, err := rs.NewRoleResource(
		role.Name,
		resourceTypeRole,
		role.Id,
		[]rs.RoleTraitOption{rs.WithRoleProfile(profile)},
	)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

func (r *roleResourceType) List(ctx context.Context, _ *v2.ResourceId, _ *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	roles, err := r.client.GetRoles(ctx)
	if err != nil {
		return nil, "", nil, fmt.Errorf("xsoar-connector: failed to list roles: %w", err)
	}

	rv := make([]*v2.Resource, 0, len(roles))
	for _, role := range roles {
		roleCopy := role

		rr, err := roleResource(ctx, &roleCopy)
		if err != nil {
			return nil, "", nil, err
		}

		rv = append(rv, rr)
	}

	return rv, "", nil, nil
}

func (r *roleResourceType) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	var rv []*v2.Entitlement

	entitlementOptions := []ent.EntitlementOption{
		ent.WithGrantableTo(resourceTypeUser),
		ent.WithDisplayName(fmt.Sprintf("%s role", resource.DisplayName)),
		ent.WithDescription(fmt.Sprintf("%s Xsoar role", resource.DisplayName)),
	}

	rv = append(rv, ent.NewAssignmentEntitlement(resource, roleMember, entitlementOptions...))

	return rv, "", nil, nil
}

func (r *roleResourceType) Grants(ctx context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	users, err := r.client.GetUsers(ctx)
	if err != nil {
		return nil, "", nil, fmt.Errorf("xsoar-connector: failed to get users: %w", err)
	}

	var rv []*v2.Grant
	for _, user := range users {
		userRoles := flattenRoleNames(user.Roles)

		if !containsRole(userRoles, resource.DisplayName) {
			continue
		}

		userCopy := user

		ur, err := userResource(ctx, &userCopy)
		if err != nil {
			return nil, "", nil, fmt.Errorf("xsoar-connector: failed to build user resource: %w", err)
		}

		rv = append(rv, grant.NewGrant(
			resource,
			roleMember,
			ur.Id,
		))
	}

	return rv, "", nil, nil
}

func (r *roleResourceType) Grant(ctx context.Context, principal *v2.Resource, entitlement *v2.Entitlement) (annotations.Annotations, error) {
	l := ctxzap.Extract(ctx)

	if principal.Id.Resource == defaultAdminUser {
		l.Warn(
			"xsoar-connector: cannot grant role memberships to default admin user",
			zap.String("principal_id", principal.Id.Resource),
		)

		return nil, fmt.Errorf("xsoar-connector: cannot grant role memberships to default admin user")
	}

	if principal.Id.ResourceType != resourceTypeUser.Id {
		l.Warn(
			"xsoar-connector: only users can be granted role membership",
			zap.String("principal_type", principal.Id.ResourceType),
			zap.String("principal_id", principal.Id.Resource),
		)

		return nil, fmt.Errorf("xsoar-connector: only users can be granted role membership")
	}

	// fetch the current user
	currentUser, err := r.client.GetCurrentUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("xsoar-connector: failed to get current user: %w", err)
	}

	// check if the principal is current user
	if principal.Id.Resource == currentUser.Id {
		l.Warn(
			"xsoar-connector: cannot grant role membership to current user",
			zap.String("principal_id", principal.Id.Resource),
			zap.String("current_user_id", currentUser.Id),
		)

		return nil, fmt.Errorf("xsoar-connector: cannot grant role membership to current user")
	}

	users, err := r.client.GetUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("xsoar-connector: failed to get users: %w", err)
	}

	targetUser := findUser(users, principal.Id.Resource)
	if targetUser == nil {
		l.Warn(
			"xsoar-connector: failed to find user to grant role membership",
			zap.String("principal_id", principal.Id.Resource),
		)

		return nil, fmt.Errorf("xsoar-connector: failed to find user to grant role membership")
	}

	userRoles := flattenRoleNames(targetUser.Roles)
	targetRole := entitlement.Resource

	// check if role to be granted is already present
	if containsRole(userRoles, targetRole.DisplayName) {
		l.Warn(
			"xsoar-connector: role membership already granted",
			zap.String("principal_id", principal.Id.Resource),
			zap.String("role", targetRole.DisplayName),
		)

		return nil, fmt.Errorf("xsoar-connector: role membership %s already granted", targetRole.DisplayName)
	}

	err = r.client.UpdateUserRoles(
		ctx,
		targetUser.Id,
		append(userRoles, targetRole.DisplayName),
	)
	if err != nil {
		return nil, fmt.Errorf("xsoar-connector: failed to update user roles: %w", err)
	}

	return nil, nil
}

func (r *roleResourceType) Revoke(ctx context.Context, grant *v2.Grant) (annotations.Annotations, error) {
	l := ctxzap.Extract(ctx)

	entitlement := grant.Entitlement
	principal := grant.Principal

	if principal.Id.Resource == defaultAdminUser {
		l.Warn(
			"xsoar-connector: cannot revoke role memberships from default admin user",
			zap.String("principal_id", principal.Id.Resource),
		)

		return nil, fmt.Errorf("xsoar-connector: cannot revoke role memberships from default admin user")
	}

	if principal.Id.ResourceType != resourceTypeUser.Id {
		l.Warn(
			"xsoar-connector: only users can have role membership revoked",
			zap.String("principal_type", principal.Id.ResourceType),
			zap.String("principal_id", principal.Id.Resource),
		)
	}

	currentUser, err := r.client.GetCurrentUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("xsoar-connector: failed to get current user: %w", err)
	}

	// check if the principal is current user
	if principal.Id.Resource == currentUser.Id {
		l.Warn(
			"xsoar-connector: cannot revoke role membership from current user",
			zap.String("principal_id", principal.Id.Resource),
			zap.String("current_user_id", currentUser.Id),
		)

		return nil, fmt.Errorf("xsoar-connector: cannot revoke role membership from current user")
	}

	users, err := r.client.GetUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("xsoar-connector: failed to get users: %w", err)
	}

	targetUser := findUser(users, principal.Id.Resource)
	if targetUser == nil {
		l.Warn(
			"xsoar-connector: failed to find user to revoke role membership",
			zap.String("principal_id", principal.Id.Resource),
		)

		return nil, fmt.Errorf("xsoar-connector: failed to find user to revoke role membership")
	}

	userRoles := flattenRoleNames(targetUser.Roles)
	targetRole := entitlement.Resource

	// check if role to be revoked is not present
	if !containsRole(userRoles, targetRole.DisplayName) {
		l.Warn(
			"xsoar-connector: role membership already revoked",
			zap.String("principal_id", principal.Id.Resource),
			zap.String("role", targetRole.DisplayName),
		)

		return nil, fmt.Errorf("xsoar-connector: %s role membership already revoked", targetRole.DisplayName)
	}

	// check if revoked role is not last one existing
	if len(userRoles) == 1 {
		l.Warn(
			"xsoar-connector: cannot revoke last role membership",
			zap.String("principal_id", principal.Id.Resource),
			zap.String("role", targetRole.DisplayName),
		)

		return nil, fmt.Errorf("xsoar-connector: cannot revoke last role membership")
	}

	// remove the role from the user roles
	updatedUserRoles := removeRole(userRoles, targetRole.DisplayName)

	err = r.client.UpdateUserRoles(
		ctx,
		targetUser.Id,
		updatedUserRoles,
	)
	if err != nil {
		return nil, fmt.Errorf("xsoar-connector: failed to update user roles: %w", err)
	}

	return nil, nil
}

func roleBuilder(client *xsoar.Client) *roleResourceType {
	return &roleResourceType{
		resourceType: resourceTypeRole,
		client:       client,
	}
}
