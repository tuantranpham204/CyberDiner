package request

type SignUp struct {
	Name            string `json:"name" binding:"required,personname"`
	Surname         string `json:"surname" binding:"required,personname"`
	Username        string `json:"username" binding:"required,username"`
	Email           string `json:"email" binding:"required,email,max=255"`
	DateOfBirth     string `json:"dob" binding:"required,dob"`
	Gender          string `json:"gender" binding:"required,gender"`
	Password        string `json:"password" binding:"required,strongpwd"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}

type SignIn struct {
	Username string `json:"username" binding:"required,username"`
	Password string `json:"password" binding:"required,min=1,max=72"`
}
