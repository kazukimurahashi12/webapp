package domain

type User struct {
	ID       string
	Password string
}

type FormUser struct {
	UserID   string `json:"userId"`
	Password string `json:"password"`
}

type UserIdChange struct {
	ChangeID string `json:"changeId"`
	NowID    string `json:"nowId"`
}

type UserPwChange struct {
	UserID         string `json:"userId"`
	NowPassword    string `json:"nowPassword"`
	ChangePassword string `json:"changePassword"`
}
