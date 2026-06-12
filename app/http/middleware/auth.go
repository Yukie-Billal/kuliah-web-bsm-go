package middleware

import (
	"kuliah-web-bsm-go/app/facades"

	"github.com/goravel/framework/contracts/http"
)

func Auth() http.Middleware {
	return func(ctx http.Context) {
		if !facades.Auth(ctx).Check() {
			ctx.Response().Json(http.StatusUnauthorized, http.Json{
				"error": "Unauthorized",
			})
			return
		}

		ctx.Request().Next()
	}
}
