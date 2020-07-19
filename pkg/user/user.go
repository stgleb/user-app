package user

type User struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Address      string `json:"address"`
	Email        string `json:"email"`
	Telephone    string `json:"telephone"`
	PasswordHash string `json:"password_hash"`
}
