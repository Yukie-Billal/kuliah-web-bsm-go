package bootstrap

import (
	"github.com/goravel/framework/contracts/database/schema"

	"kuliah-web-bsm-go/database/migrations"
)

func Migrations() []schema.Migration {
	return []schema.Migration{
		&migrations.M20210101000001CreateJobsTable{},
		&migrations.M20260424124001CreateRolesTable{},
		&migrations.M20260424124002CreateUsersTable{},
		&migrations.M20260424124003CreateCustomersTable{},
		&migrations.M20260424124004CreateStudiosTable{},
		&migrations.M20260424124005CreateAlatMusiksTable{},
		&migrations.M20260424124006CreateBookingsTable{},
	}
}
