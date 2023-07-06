package connector

import (
	"context"
	"crypto/tls"

	"github.com/ConductorOne/baton-demisto/pkg/demisto"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"github.com/conductorone/baton-sdk/pkg/uhttp"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
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

type Demisto struct {
	client *demisto.Client
}

func (de *Demisto) ResourceSyncers(ctx context.Context) []connectorbuilder.ResourceSyncer {
	return []connectorbuilder.ResourceSyncer{
		userBuilder(de.client),
		roleBuilder(de.client),
	}
}

func (de *Demisto) Metadata(ctx context.Context) (*v2.ConnectorMetadata, error) {
	return &v2.ConnectorMetadata{
		DisplayName: "Demisto",
		Description: "Connector syncing Demisto/Cortex XSOAR users and their roles to Baton.",
	}, nil
}

func (de *Demisto) Validate(ctx context.Context) (annotations.Annotations, error) {
	// TODO: implement validation

	return nil, nil
}

func New(ctx context.Context, token string, unsafe bool) (*Demisto, error) {
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

	return &Demisto{
		client: demisto.NewClient(httpClient, token),
	}, nil
}
