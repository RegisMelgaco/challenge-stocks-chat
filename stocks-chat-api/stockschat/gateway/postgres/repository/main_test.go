package repository_test

import (
	"os"
	"testing"

	"local/challengestockschat/stockschat/gateway/postgres/migration"

	"github.com/regismelgaco/go-sdks/postgres"
)

func TestMain(m *testing.M) {
	postgres.SetMigrationFunc(migration.Migrate)
	teardown := postgres.SetupPgContainer()

	code := m.Run()

	teardown()

	os.Exit(code)
}
