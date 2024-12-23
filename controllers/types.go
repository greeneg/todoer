package controllers

import "github.com/greeneg/todoer/globals"

type TodoerService struct {
	AppPath    string
	ConfigPath string
	ConfStruct globals.Config
}

type SafeUser struct {
	Id           int    `json:"Id"`
	UserName     string `json:"userName"`
	CreationDate string `json:"creationDate"`
}
