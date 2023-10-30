package model

type TaskAssignee struct {
	UserID uint `json:"user_id"`
	TaskID uint `json:"task_id"`
}
