package connector

import (
	"context"
	"crypto/tls"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"github.com/conductorone/baton-sdk/pkg/uhttp"
	"github.com/conductorone/baton-xsoar/pkg/xsoar"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	resourceTypeUser = &v2.ResourceType{
		Id:          "user",
		DisplayName: "User",
		Traits: []v2.ResourceType_Trait{
			v2.ResourceType_TRAIT_USER,
		},
		Annotations: annotationsForUserResourceType(),
	}
	resourceTypeRole = &v2.ResourceType{
		Id:          "role",
		DisplayName: "Role",
		Traits: []v2.ResourceType_Trait{
			v2.ResourceType_TRAIT_ROLE,
		},
	}
)

type Xsoar struct {
	client *xsoar.Client
}

func (xs *Xsoar) ResourceSyncers(ctx context.Context) []connectorbuilder.ResourceSyncer {
	return []connectorbuilder.ResourceSyncer{
		userBuilder(xs.client),
		roleBuilder(xs.client),
	}
}

func (xs *Xsoar) Metadata(ctx context.Context) (*v2.ConnectorMetadata, error) {
	return &v2.ConnectorMetadata{
		DisplayName: "Xsoar",
		Description: "Connector syncing Xsoar/Cortex XSOAR users and their roles to Baton.",
	}, nil
}

// Validate method checks for compatible and valid credentials.
// Since each role can be configured to have different permissions,
// we need to check if the provided credentials have the required
// permissions to perform the operations.
func (xs *Xsoar) Validate(ctx context.Context) (annotations.Annotations, error) {
	currentUser, err := xs.client.GetCurrentUser(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Provided Access Token is invalid - unable to get current user")
	}

	users, err := xs.client.GetUsers(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Provided Access Token is invalid - unable to get users")
	}

	_, err = xs.client.GetRoles(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Provided Access Token is invalid - unable to get roles")
	}

	targetUsers := removeUsers(users, "admin", currentUser.Id)
	if len(targetUsers) == 0 {
		return nil, nil
	}

	err = xs.client.UpdateUserRoles(
		ctx,
		targetUsers[0].Id,
		flattenRoleNames(targetUsers[0].Roles),
	)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Provided Access Token is invalid - unable to update user roles")
	}

	return nil, nil
}

func New(ctx context.Context, token, apiUrl string, unsafe bool) (*Xsoar, error) {
	options := []uhttp.Option{
		uhttp.WithLogger(true, ctxzap.Extract(ctx)),
	}

	// Skip TLS verification if flag `unsafe` is specified.
	if unsafe { // #nosec G402
		options = append(
			options,
			uhttp.WithTLSClientConfig(
				&tls.Config{InsecureSkipVerify: true},
			),
		)
	}

	httpClient, err := uhttp.NewClient(
		ctx,
		options...,
	)
	if err != nil {
		return nil, err
	}

	return &Xsoar{
		client: xsoar.NewClient(httpClient, token, apiUrl),
	}, nil
}
