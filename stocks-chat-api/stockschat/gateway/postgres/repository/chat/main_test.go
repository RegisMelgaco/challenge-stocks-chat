package chat_test

import (
	"local/challengestockschat/stockschat/gateway/postgres/migration"
	"os"
	"testing"

	"github.com/regismelgaco/go-sdks/postgres"
)

func TestMain(m *testing.M) {
	postgres.SetMigrationFunc(migration.Migrate)
	teardown := postgres.SetupPgContainer()

	code := m.Run()

	teardown()

	os.Exit(code)
}
