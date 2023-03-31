package users

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserInfo struct {
	Username string `json:"username"`
	Balance  string `json:"balance"`
	Currency string `json:"currency"`
}
