package auth

// maps tokens to user ids
var TokenStore map[string]string

type Authorizer struct{}

func NewAuthorizer() *Authorizer {
	return &Authorizer{}
}

func init() {
	TokenStore = make(map[string]string)
}
