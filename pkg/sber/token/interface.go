package token

type TokenGetter interface {
	GetToken() (string, error)
}
