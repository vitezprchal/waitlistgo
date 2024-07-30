package forms

type EmailForm struct {
	Name  string `form:"name" binding:"required"`
	Email string `form:"email" binding:"required,email"`
	Terms string `form:"terms" binding:"required"`
}
