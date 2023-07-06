package demisto

import (
	"context"
	"encoding/json"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const BaseURL = "https://localhost:8443"
const UsersBaseURL = BaseURL + "/users"
const RolesBaseURL = BaseURL + "/roles"

type Client struct {
	httpClient *http.Client
	Token      string
}

type UsersResponse = []User
type RolesResponse = []Role

func NewClient(httpClient *http.Client, token string) *Client {
	return &Client{
		httpClient: httpClient,
		Token:      token,
	}
}

func (c *Client) GetUsers(ctx context.Context) ([]User, error) {
	var usersResponse UsersResponse

	err := c.doRequest(ctx, UsersBaseURL, &usersResponse)
	if err != nil {
		return nil, err
	}

	return usersResponse, nil
}

func (c *Client) GetRoles(ctx context.Context) ([]Role, error) {
	var rolesResponse RolesResponse

	err := c.doRequest(ctx, RolesBaseURL, &rolesResponse)
	if err != nil {
		return nil, err
	}

	return rolesResponse, nil
}

func (c *Client) doRequest(
	ctx context.Context,
	urlAddress string,
	resourceResponse interface{},
) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlAddress, nil)
	if err != nil {
		return err
	}

	req.Header.Set("content-type", "application/json")
	req.Header.Set("Authorization", c.Token)
	req.Header.Set("Accept", "application/json")

	rawResponse, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer rawResponse.Body.Close()

	if rawResponse.StatusCode >= 300 {
		return status.Error(codes.Code(rawResponse.StatusCode), "Request failed")
	}

	if err := json.NewDecoder(rawResponse.Body).Decode(&resourceResponse); err != nil {
		return err
	}

	return nil
}
