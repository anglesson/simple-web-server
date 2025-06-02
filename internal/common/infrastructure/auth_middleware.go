package common_infrastructure

import "context"

func GetUserFromContext(ctx context.Context) (*string, bool) {
	user, ok := ctx.Value(LoggedUserKey).(*string)
	return user, ok
}
