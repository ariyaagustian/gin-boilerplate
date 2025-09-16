package dto

type RegisterReq struct {
	Name     string `json:"name"     example:"Ariya"`
	Email    string `json:"email"    example:"user@mail.com"`
	Password string `json:"password" example:"secret123"`
}

type User struct {
	ID    string `json:"id"    example:"8d7a9b6e-..."`
	Name  string `json:"name"  example:"Ariya"`
	Email string `json:"email" example:"user@mail.com"`
	Role  string `json:"role"  example:"user"`
}

type RegisterResp struct {
	User      User   `json:"user"`
	Token     string `json:"token"      example:"eyJhbGciOiJI..."`
	TokenType string `json:"token_type" example:"Bearer"`
}
type LoginReq struct {
	Email    string `json:"email"    example:"user@mail.com"`
	Password string `json:"password" example:"secret123"`
}

type TokenResp struct {
	User      User   `json:"user"`
	Token     string `json:"token"      example:"eyJhbGciOiJI..."`
	TokenType string `json:"token_type" example:"Bearer"`
}

type SetPasswordReq struct {
	Email       string `json:"email"        example:"user@mail.com"`
	NewPassword string `json:"new_password" example:"newsecret123"`
}

type SetPasswordResp struct {
	Message string `json:"message" example:"password updated"`
}
