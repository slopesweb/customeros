package repository

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/postgres"
	"github.com/testcontainers/testcontainers-go"
	"log"
	"os"
	"testing"
)

var (
	neo4jContainer testcontainers.Container
	driver         *neo4j.DriverWithContext

	repositories *Repositories
)

const tenantName = "openline"

func TestMain(m *testing.M) {
	neo4jContainer, driver = neo4jt.InitTestNeo4jDB()
	defer func(dbContainer testcontainers.Container, driver neo4j.DriverWithContext, ctx context.Context) {
		neo4jt.CloseDriver(driver)
		neo4jt.Terminate(dbContainer, ctx)
	}(neo4jContainer, *driver, context.Background())

	postgresContainer, postgresGormDB, _ := postgres.InitTestDB()
	defer func(postgresContainer testcontainers.Container, ctx context.Context) {
		err := postgresContainer.Terminate(ctx)
		if err != nil {
			log.Fatal("Error during container termination")
		}
	}(postgresContainer, context.Background())

	repositories = InitRepos(driver, "neo4j", postgresGormDB)

	os.Exit(m.Run())
}
