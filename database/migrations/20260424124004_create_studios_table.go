package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"kuliah-web-bsm-go/app/facades"
)

type M20260424124004CreateStudiosTable struct{}

func (r *M20260424124004CreateStudiosTable) Signature() string {
	return "20260424124004_create_studios_table"
}

func (r *M20260424124004CreateStudiosTable) Up() error {
	if !facades.Schema().HasTable("studios") {
		if err := facades.Schema().Create("studios", func(table schema.Blueprint) {
			table.ID()
			table.String("nama_studio")
			table.String("lokasi")
			table.Integer("luas")
			table.Text("jam_operasional").Nullable()
			table.Integer("harga_per_jam").Default(75000)
			table.Boolean("is_active").Default(true)
			table.Timestamps()
		}); err != nil {
			return err
		}
	}
	return nil
}

func (r *M20260424124004CreateStudiosTable) Down() error {
	return facades.Schema().DropIfExists("studios")
}
