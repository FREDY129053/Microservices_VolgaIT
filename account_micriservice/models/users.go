package models

// Регистрация пользователя
type SignupUser struct {
	LastName  string `json:"last_name"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

// Информация о пользователи
type UserInfo struct {
	UUID      string   `json:"uuid"`
	LastName  string   `json:"last_name"`
	FirstName string   `json:"first_name"`
	Username  string   `json:"username"`
	Password  string   `json:"password,omitempty"`
	Roles     []string `json:"roles,omitempty"`
}

// Вход пользователя
type SigninUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Обновление пользователя(своего акка)
type UpdateUser struct {
	LastName  string `json:"last_name"`
	FirstName string `json:"first_name"`
	Password  string `json:"password"`
}

// Действия админов с акками
type AdminAccounts struct {
	LastName  string   `json:"last_name"`
	FirstName string   `json:"first_name"`
	Username  string   `json:"username"`
	Password  string   `json:"password"`
	Roles     []string `json:"roles"`
}

// Выборка докторов
type DoctorsInfo struct {
	UUID      string `json:"uuid"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}
