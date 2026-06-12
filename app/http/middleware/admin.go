package middleware

import (
	"kuliah-web-bsm-go/app/facades"
	"kuliah-web-bsm-go/app/models"

	"github.com/goravel/framework/contracts/http"
)

func Admin() http.Middleware {
	return func(ctx http.Context) {
		if !facades.Auth(ctx).Check() {
			ctx.Response().Json(http.StatusUnauthorized, http.Json{
				"error": "Unauthorized",
			})
			return
		}

		var user models.User
		if err := facades.Auth(ctx).User(&user); err != nil {
			ctx.Response().Json(http.StatusUnauthorized, http.Json{
				"error": "Unauthorized",
			})
			return
		}

		var role models.Role
		if err := facades.Orm().Query().Where("id", user.RoleID).First(&role); err != nil || role.Role != "admin" {
			ctx.Response().Json(http.StatusForbidden, http.Json{
				"error": "Forbidden: Admin access required",
			})
			return
		}

		ctx.Request().Next()
	}
}
