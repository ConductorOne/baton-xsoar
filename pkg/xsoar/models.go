package xsoar

type BaseResource struct {
	Id      string `json:"id"`
	Version int    `json:"version"`
}

type User struct {
	BaseResource

	Username  string `json:"username"`
	Name      string `json:"name"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`

	Email string              `json:"email"`
	Roles map[string][]string `json:"roles"`

	Disabled bool `json:"disabled"`
}

type Role struct {
	BaseResource
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
}
