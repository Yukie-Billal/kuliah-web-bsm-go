package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"kuliah-web-bsm-go/app/facades"
)

type M20260424124006CreateBookingsTable struct{}

func (r *M20260424124006CreateBookingsTable) Signature() string {
	return "20260424124006_create_bookings_table"
}

func (r *M20260424124006CreateBookingsTable) Up() error {
	if !facades.Schema().HasTable("bookings") {
		if err := facades.Schema().Create("bookings", func(table schema.Blueprint) {
			table.ID()
			table.String("booking_id")
			table.Unique("booking_id")
			table.UnsignedBigInteger("customer_id")
			table.Foreign("customer_id").References("id").On("customers").CascadeOnDelete()
			table.UnsignedBigInteger("studio_id")
			table.Foreign("studio_id").References("id").On("studios")
			table.Date("tanggal")
			table.String("jam")
			table.Integer("durasi").Default(1)
			table.Integer("total_biaya").Default(75000)
			table.String("status").Default("pending")
			table.Text("catatan").Nullable()
			table.Timestamp("checked_in_at").Nullable()
			table.UnsignedBigInteger("approved_by").Nullable()
			table.Foreign("approved_by").References("id").On("users")
			table.Timestamps()

			// Composite unique key
			table.Unique("studio_id", "tanggal", "jam")
		}); err != nil {
			return err
		}
	}
	return nil
}

func (r *M20260424124006CreateBookingsTable) Down() error {
	return facades.Schema().DropIfExists("bookings")
}
