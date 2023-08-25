package demisto

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const BaseURL = "https://%s"
const CurrentUserBaseURL = BaseURL + "/user"
const UsersBaseURL = BaseURL + "/users"
const RolesBaseURL = BaseURL + "/roles"
const UpdateUserBaseURL = BaseURL + "/users/update"

type Client struct {
	httpClient *http.Client
	Token      string
	Domain     string
}

type UsersResponse = []User
type RolesResponse = []Role

func NewClient(httpClient *http.Client, token, domain string) *Client {
	return &Client{
		httpClient: httpClient,
		Token:      token,
		Domain:     domain,
	}
}

func (c *Client) GetUsers(ctx context.Context) ([]User, error) {
	var usersResponse UsersResponse

	err := c.doRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf(UsersBaseURL, c.Domain),
		&usersResponse,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return usersResponse, nil
}

func (c *Client) GetRoles(ctx context.Context) ([]Role, error) {
	var rolesResponse RolesResponse

	err := c.doRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf(RolesBaseURL, c.Domain),
		&rolesResponse,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return rolesResponse, nil
}

func (c *Client) GetCurrentUser(ctx context.Context) (*User, error) {
	var user User

	err := c.doRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf(CurrentUserBaseURL, c.Domain),
		&user,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

type Roles struct {
	Roles []string `json:"roles"`
}

type UpdateRolesBody struct {
	Id    string `json:"id"`
	Roles Roles  `json:"roles"`
}

func (c *Client) UpdateUserRoles(ctx context.Context, userId string, roleIds []string) error {
	data := UpdateRolesBody{
		Id: userId,
		Roles: Roles{
			Roles: roleIds,
		},
	}

	err := c.doRequest(
		ctx,
		http.MethodPost,
		fmt.Sprintf(UpdateUserBaseURL, c.Domain),
		nil,
		&data,
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) doRequest(
	ctx context.Context,
	method string,
	urlAddress string,
	resourceResponse interface{},
	data interface{},
) error {
	var body io.Reader

	if data != nil {
		jsonBody, err := json.Marshal(data)
		if err != nil {
			return err
		}

		body = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		method,
		urlAddress,
		body,
	)
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
