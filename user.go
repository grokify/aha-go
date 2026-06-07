package aha

import (
	"context"
	"time"

	"github.com/grokify/aha-go/internal/api"
)

// User represents an Aha user.
type User struct {
	ID        string
	FirstName string
	LastName  string
	Email     string
	Role      string
	CreatedAt *time.Time
}

// Name returns the user's full name.
func (u *User) Name() string {
	if u.FirstName != "" && u.LastName != "" {
		return u.FirstName + " " + u.LastName
	}
	if u.FirstName != "" {
		return u.FirstName
	}
	return u.LastName
}

// UserList represents a paginated list of users.
type UserList struct {
	Users      []User
	Pagination Pagination
}

// GetUser retrieves a user by ID or email.
func (c *Client) GetUser(ctx context.Context, id string) (*User, error) {
	resp, err := c.apiClient.GetUser(ctx, api.GetUserParams{
		UserID: id,
	})
	if err != nil {
		return nil, wrapError("GetUser", err)
	}

	if u, ok := resp.User.Get(); ok {
		return userFromAPI(u), nil
	}
	return nil, &APIError{StatusCode: 404, Message: "user not found"}
}

// GetCurrentUser retrieves the currently authenticated user.
func (c *Client) GetCurrentUser(ctx context.Context) (*User, error) {
	resp, err := c.apiClient.GetCurrentUser(ctx)
	if err != nil {
		return nil, wrapError("GetCurrentUser", err)
	}

	if u, ok := resp.User.Get(); ok {
		return userFromAPI(u), nil
	}
	return nil, &APIError{StatusCode: 500, Message: "unexpected response: user not returned"}
}

// ListUsers lists all users in the account.
func (c *Client) ListUsers(ctx context.Context, opts ...ListOption) (*UserList, error) {
	listOpts := applyListOptions(opts...)

	params := api.ListUsersParams{}
	if listOpts.Page > 0 {
		params.Page = api.NewOptInt32(int32(listOpts.Page)) //nolint:gosec // G115: Page number bounded by API limits
	}
	if listOpts.PerPage > 0 {
		params.PerPage = api.NewOptInt32(int32(listOpts.PerPage)) //nolint:gosec // G115: PerPage bounded by API limits
	}

	resp, err := c.apiClient.ListUsers(ctx, params)
	if err != nil {
		return nil, wrapError("ListUsers", err)
	}

	return userListFromAPI(resp), nil
}

// userFromAPI converts an API user to a domain user.
func userFromAPI(u api.User) *User {
	user := &User{}
	if v, ok := u.ID.Get(); ok {
		user.ID = v
	}
	if v, ok := u.FirstName.Get(); ok {
		user.FirstName = v
	}
	if v, ok := u.LastName.Get(); ok {
		user.LastName = v
	}
	if v, ok := u.Email.Get(); ok {
		user.Email = v
	}
	if v, ok := u.Role.Get(); ok {
		user.Role = string(v)
	}
	if v, ok := u.CreatedAt.Get(); ok {
		user.CreatedAt = &v
	}
	return user
}

// userListFromAPI converts an API users response to a domain user list.
func userListFromAPI(resp *api.UsersResponse) *UserList {
	list := &UserList{}

	list.Users = make([]User, len(resp.Users))
	for i, u := range resp.Users {
		list.Users[i] = *userFromAPI(u)
	}

	if v, ok := resp.Pagination.Get(); ok {
		list.Pagination = paginationFromAPI(v)
	}

	return list
}
