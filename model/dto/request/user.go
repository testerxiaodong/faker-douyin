package request

type UserRegisterReq struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}
