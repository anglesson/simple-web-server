package middlewares

import (
	"context"

	"github.com/anglesson/simple-web-server/internal/infrastructure/http_server/utils"
)

func GetUserFromContext(ctx context.Context) (*string, bool) {
	user, ok := ctx.Value(utils.LoggedUserKey).(*string)
	return user, ok
}
