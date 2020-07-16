package user

type User struct {
	FullName  string `json:"full_name"`
	Address   string `json:"address"`
	Email     string `json:"email"`
	Telephone string `json:"telephone"`
}
