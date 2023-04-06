package domain

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserRes struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
}

type LoginUserReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserRes struct {
	AccessToken string `json:"accessToken"`
	ID          string `json:"id"`
	Username    string `json:"username"`
}

type UpdateUsernameReq struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type PublicUser struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
