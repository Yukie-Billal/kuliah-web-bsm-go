package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"kuliah-web-bsm-go/app/facades"
)

type M20260424124002CreateUsersTable struct{}

func (r *M20260424124002CreateUsersTable) Signature() string {
	return "20260424124002_create_users_table"
}

func (r *M20260424124002CreateUsersTable) Up() error {
	if !facades.Schema().HasTable("users") {
		if err := facades.Schema().Create("users", func(table schema.Blueprint) {
			table.ID()
			table.UnsignedBigInteger("role_id").Default(2)
			table.Foreign("role_id").References("id").On("roles")
			table.String("username")
			table.String("email")
			table.Unique("email")
			table.String("password").Nullable()
			table.String("google_id").Nullable()
			table.Unique("google_id")
			table.Timestamp("email_verified_at").Nullable()
			table.String("avatar").Nullable()
			table.String("remember_token").Nullable()
			table.Timestamps()
		}); err != nil {
			return err
		}
	}
	return nil
}

func (r *M20260424124002CreateUsersTable) Down() error {
	return facades.Schema().DropIfExists("users")
}
