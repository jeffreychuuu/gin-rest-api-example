package models

type Result struct {
	Code    int         `json:"code" example:"000"`
	Message string      `json:"message" example:"Request Message"`
	Data    interface{} `json:"data" `
}
