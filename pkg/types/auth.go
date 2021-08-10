package types

type Authentication struct {
	Account  string
	Email    string
	Password string
}

type Token struct {
	Email       string
	TokenString string
}
