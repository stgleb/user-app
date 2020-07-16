package user

type User struct {
	Id           string `json:"id"`
	UserName     string `json:"user_name"`
	IsActive     bool   `json:"is_active"`
	FullName     string `json:"full_name"`
	Address      string `json:"address"`
	Email        string `json:"email"`
	Telephone    string `json:"telephone"`
	PasswordHash string `json:"password_hash"`
}
