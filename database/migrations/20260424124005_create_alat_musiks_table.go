package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"kuliah-web-bsm-go/app/facades"
)

type M20260424124005CreateAlatMusiksTable struct{}

func (r *M20260424124005CreateAlatMusiksTable) Signature() string {
	return "20260424124005_create_alat_musiks_table"
}

func (r *M20260424124005CreateAlatMusiksTable) Up() error {
	if !facades.Schema().HasTable("alat_musiks") {
		if err := facades.Schema().Create("alat_musiks", func(table schema.Blueprint) {
			table.ID()
			table.UnsignedBigInteger("studio_id")
			table.Foreign("studio_id").References("id").On("studios").CascadeOnDelete()
			table.String("nama_alat")
			table.String("kondisi").Default("Baik")
			table.Text("keterangan").Nullable()
			table.Timestamps()
		}); err != nil {
			return err
		}
	}
	return nil
}

func (r *M20260424124005CreateAlatMusiksTable) Down() error {
	return facades.Schema().DropIfExists("alat_musiks")
}
