package common_infrastructure

// UserContextKey é a chave usada para armazenar o User logado no context.Context.
type UserContextKey string

const LoggedUserKey UserContextKey = "loggedInUser"
