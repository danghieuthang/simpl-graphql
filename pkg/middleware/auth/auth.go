package auth

type AuthenticatedUser struct {
	Id    int
	Email string
	Name  string
	Role  string
}
