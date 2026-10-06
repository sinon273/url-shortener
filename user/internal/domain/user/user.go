package user_domain

type User struct {
	ID           string
	PasswordHash string
	Email        string
}

func NewUser(ID string, passwordHash string, email string) User {
	return User{
		ID:           ID,
		PasswordHash: passwordHash,
		Email:        email,
	}
}
