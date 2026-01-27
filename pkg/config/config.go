package config

import (
	"github.com/conductorone/baton-sdk/pkg/field"
)

var (
	AccessToken = field.StringField(
		"token",
		field.WithDescription("Access token used to connect to the Cortex XSOAR API."),
		field.WithRequired(true),
		field.WithIsSecret(true),
	)

	ApiUrl = field.StringField(
		"api-url",
		field.WithDescription("The API URL of the Cortex XSOAR instance."),
		field.WithRequired(true),
	)

	Unsafe = field.BoolField(
		"unsafe",
		field.WithDescription("Allow insecure TLS connections to Cortex XSOAR instance."),
	)
)

//go:generate go run ./gen
var Config = field.NewConfiguration([]field.SchemaField{
	AccessToken,
	ApiUrl,
	Unsafe,
})
