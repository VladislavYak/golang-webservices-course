package models

type Task struct {
	TaskID     int    `json:"task_id"`
	TaskName   string `json:"task_name"`
	Asignee    string `json:"asignee"`
	ByUserID   string `json:"creator_id"`
	ByUserName string `json:"creator_username"`
}

type User struct {
	ChatID       int    `json:"chat_id"`
	UserID       string `json:"user_id"`
	UserName     string `json:"username"`
	LanguageCode string `json:"language_code"`
}

type NewParams struct {
	UserId   int    `json:"user_id"`
	Username string `json:"username"`
	TaskName string `json:"taskname"`
}

type ParamsOut struct {
}
