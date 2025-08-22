package define

type SearchResp struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
}

type LoginResp struct {
	Token    string `json:"token"`
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Role     int    `json:"role"`
}
