package postgresreplication

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	_ "github.com/lib/pq"
)

const (
	testAccSlotCreate = `
        resource "postgresreplication_slot" "test_slot" {
            slot_name 	   = "test_slot"
            output_plugin  = "wal2json"
            database       = "postgres"
        }
`
)

func TestAccResourceReplicationSlot(t *testing.T) {
	resourceName := "postgresreplication_slot.test_slot"

	resource.Test(t, resource.TestCase{
		IDRefreshName: resourceName,
		Providers: map[string]terraform.ResourceProvider{
			"postgresreplication": Provider(),
		},
		CheckDestroy: testAccCheckSlotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSlotCreate,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckSlotExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, slotNameAttributeName, "test_slot"),
					resource.TestCheckResourceAttr(resourceName, outputPluginAttributeName, "wal2json"),
					resource.TestCheckResourceAttr(resourceName, databaseAttributeName, "postgres"),
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

func slotExists(ID string) (bool, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", defaultUser, defaultPassword, "postgres", "disable"))
	if err != nil {
		return false, err
	}
	defer db.Close()
	rows, err := db.Query("select * from pg_replication_slots where slot_name=$1", ID)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	if rows.Next() {
		return true, nil
	}
	return false, nil
}

func testAccCheckSlotExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		r, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %q", resourceName)
		}
		exists, err := slotExists(r.Primary.ID)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("not created: %q", resourceName)
		}
		return nil
	}
}

func testAccCheckSlotDestroy(s *terraform.State) error {
	for _, r := range s.RootModule().Resources {
		if r.Type == replicationSlotName {
			exists, err := slotExists(r.Primary.ID)
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
