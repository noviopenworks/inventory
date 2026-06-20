package models

// User represents a person who can be assigned devices or software.
type User struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Surname   *string `json:"surname"`
	Status    string  `json:"status"`
	Notes     *string `json:"notes"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
}

// UserInput carries the user-supplied fields for creating or updating a user.
type UserInput struct {
	Name    string  `json:"name"`
	Surname *string `json:"surname"`
	Status  string  `json:"status"`
	Notes   *string `json:"notes"`
}
