package dto

type TaskAssigneeRequest struct {
	TaskID uint `json:"task_id"`
	UserID uint `json:"user_id"`
}

type TaskAssigneeResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
}