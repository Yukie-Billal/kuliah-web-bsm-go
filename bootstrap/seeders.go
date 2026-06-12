package bootstrap

import (
	"github.com/goravel/framework/contracts/database/seeder"

	"kuliah-web-bsm-go/database/seeders"
)

func Seeders() []seeder.Seeder {
	return []seeder.Seeder{
		&seeders.DatabaseSeeder{},
	}
}
