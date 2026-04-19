package requests

type CreateUserRequest struct {
	Login string `json:"login" binding:"required"`
	Role  string `json:"role" binding:"required"`
}
