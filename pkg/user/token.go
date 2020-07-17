package user

type Token struct {
	Value string `json:"value"`
	Email string `json:"email"`
	Used  bool   `json:"used"`
}
