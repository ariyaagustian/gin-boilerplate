package dto

type ListUsersResp struct {
	Total int    `json:"total" example:"100"`
	Limit int    `json:"limit" example:"20"`
	Page  int    `json:"page"  example:"1"`
	Data  []User `json:"data"`
}

type GetUserResp struct {
	Data User `json:"data"`
}

type DeleteUserResp struct {
	Message string `json:"message" example:"user deleted"`
}

type CreateUserReq struct {
	Name  string `json:"name"  example:"Ariya"`
	Email string `json:"email" example:"test@mail.com"`
}

type UpdateUserReq struct {
	Name  string `json:"name"  example:"Ariya"`
	Email string `json:"email" example:"test@mail.com"`
}

type UpdateUserResp struct {
	Data User `json:"data"`
}
