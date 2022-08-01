package model

func (User) IsNode() {}

type User struct {
	ID                      string                  `json:"id"`
	DatabaseID              int                     `json:"databaseId"`
	Name                    string                  `json:"name"`
	Email                   string                  `json:"email"`
	UnverifiedEmail         *string                 `json:"unverifiedEmail"`
	Avatar                  string                  `json:"avatar"`
	Introduction            *string                 `json:"introduction"`
	Role                    Role                    `json:"role"`
	EmailVerificationStatus EmailVerificationStatus `json:"emailVerificationStatus"`
}
