package model

// primary object structs

type HealthCheck struct {
	Db           string `json:"db"`
	DiskSpace    string `json:"diskSpace"`
	DiskWritable string `json:"diskWritable"`
	Health       string `json:"health"`
	Status       int    `json:"status"`
}

type PasswordChange struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type Status struct {
	Id           int    `json:"Id"`
	StatusString string `json:"statusString"`
}

type Todo struct {
	Id           int    `json:"Id"`
	Description  string `json:"description"`
	Status       string `json:"status"`
	CreationDate string `json:"creationDate"`
}

type TodoList struct {
	Data []Todo `json:"data"`
}

type User struct {
	Id              int    `json:"Id"`
	UserName        string `json:"userName"`
	Status          string `json:"status"`
	PasswordHash    string `json:"passwordHash"`
	CreationDate    string `json:"creationDate"`
	LastChangedDate string `json:"lastChangedDate"`
}

type UserStatus struct {
	Status string `json:"status" enum:"enabled,disabled"`
}

type UserStatusMsg struct {
	Message    string `json:"message"`
	UserStatus string `json:"userStatus" enum:"enabled,disabled"`
}

// proposed object structs. Normally used when creating new DB entries

type ProposedTodo struct {
	Description string `json:"description"`
}

type ProposedUser struct {
	Id       int    `json:"Id"`
	UserName string `json:"userName"`
	Status   string `json:"status" enum:"enabled,disabled"`
	Password string `json:"password"`
}

// list object structs

type UsersList struct {
	Data []User `json:"data"`
}

// generic message structs

type FailureMsg struct {
	Error string `json:"error"`
}

type SuccessMsg struct {
	Message string `json:"message"`
}
