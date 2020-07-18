package user

type User struct {
	Id           string   `json:"id"`
	Name         string   `json:"name"`
	FullName     string   `json:"full_name"`
	Address      string   `json:"address"`
	Email        string   `json:"email"`
	Telephone    string   `json:"telephone"`
	PasswordHash []byte `json:"password_hash"`
}
