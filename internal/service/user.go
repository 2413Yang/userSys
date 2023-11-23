package service

// RegisterRequest 注册请求
type RegisterRequest struct {
	UserName string `json:"user_name"`
	Password string `json:"pass_word"`
	Age      int    `json:"age"`
	Gender   string `json:"gender"`
	NickName string `json:"nikc_name"`
}
