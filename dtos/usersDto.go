package dtos

type RegisterRequestDto struct {
	FirstName string `form:"first_name" json:"first_name" xml:"first_name" binding:"required"`
	LastName  string `form:"last_name" json:"last_name" xml:"last_name" binding:"required"`
	Email     string `form:"email" json:"email" xml:"email" binding:"required"`
	Password  string `form:"password" json:"password" xml:"password" binding:"required"`
}

type LoginRequestDto struct {
	Email    string `form:"email" json:"email" xml:"email" binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}
