package postgresreplication

import (
	"context"
	"fmt"
	"net/url"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/jackc/pgx/v4"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

const (
	testAccReplicaIdentityCreate = `
resource "postgresreplication_replica_identity" "test" {
  table_name = "foo"
  database = "postgres"
}
`
)

func TestAccResourceReplicaIdentityCreate(t *testing.T) {
	resourceName := "postgresreplication_replica_identity.test"

	resource.Test(t, resource.TestCase{
		IDRefreshName: resourceName,
		Providers: map[string]terraform.ResourceProvider{
			"postgresreplication": Provider(),
		},
		CheckDestroy: testAccCheckSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccReplicaIdentityCreate,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccReplicaIdentityExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, tableNameAttributeName, "foo"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func replicaIdentityExists(ID string) (bool, error) {
	u, err := url.Parse(fmt.Sprintf("postgres://%s:%d/%s?sslmode=%s",
		defaultUser,
		5432,
		"postgres",
		"disable"))
	if err != nil {
		return false, errors.Wrap(err, "error contructing database connection uri.")
	}

	u.User = url.UserPassword(defaultUser, defaultPassword)
	conn, err := pgx.Connect(context.Background(), u.String())
	if err != nil {
		return false, errors.Wrap(err, "error connecting to database.")
	}

	return IsReplicaIdentityFull(conn, "foo")

}

func testAccReplicaIdentityExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		_, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %q", resourceName)
		}

		return nil
	}
}

func testAccReplicaIdentityDestroy(s *terraform.State) error {
	for _, r := range s.RootModule().Resources {
		if r.Type == replicationSlotName {
			exists, err := replicaIdentityExists(r.Primary.ID)
			if err != nil {
				return err
			}
			if exists {
				return fmt.Errorf("still exists: %q", r.Primary.ID)
			}
		}
	}
	return nil
}
