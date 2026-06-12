package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"kuliah-web-bsm-go/app/facades"
)

type M20260424124003CreateCustomersTable struct{}

func (r *M20260424124003CreateCustomersTable) Signature() string {
	return "20260424124003_create_customers_table"
}

func (r *M20260424124003CreateCustomersTable) Up() error {
	if !facades.Schema().HasTable("customers") {
		if err := facades.Schema().Create("customers", func(table schema.Blueprint) {
			table.ID()
			table.UnsignedBigInteger("user_id")
			table.Foreign("user_id").References("id").On("users").CascadeOnDelete()
			table.String("nama")
			table.String("telepon").Nullable()
			table.Timestamps()
		}); err != nil {
			return err
		}
	}
	return nil
}

func (r *M20260424124003CreateCustomersTable) Down() error {
	return facades.Schema().DropIfExists("customers")
}
