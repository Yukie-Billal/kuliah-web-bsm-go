package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"kuliah-web-bsm-go/app/facades"
)

type M20260424124001CreateRolesTable struct{}

func (r *M20260424124001CreateRolesTable) Signature() string {
	return "20260424124001_create_roles_table"
}

func (r *M20260424124001CreateRolesTable) Up() error {
	if !facades.Schema().HasTable("roles") {
		if err := facades.Schema().Create("roles", func(table schema.Blueprint) {
			table.ID()
			table.String("role")
			table.Timestamps()
		}); err != nil {
			return err
		}
	}
	return nil
}

func (r *M20260424124001CreateRolesTable) Down() error {
	return facades.Schema().DropIfExists("roles")
}
