package postgresreplication

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

func resourceReplicaIdentity() *schema.Resource {
	return &schema.Resource{
		Create: resourceReplicaIdentityCreate,
		Read:   resourceReplicaIdentityRead,
		Delete: resourceReplicaIdentityDelete,
		Importer: &schema.ResourceImporter{
			State: resourceReplicaIdentityImport,
		},
		SchemaVersion: 0,
		Schema: map[string]*schema.Schema{
			tableNameAttributeName: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the table you want to alter replication identity to full",
			},
			databaseAttributeName: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The database where the table resides",
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
	}
}

func IsReplicaIdentityFull(conn *pgx.Conn, tableName string) (bool, error) {
	//Case statement needed, otherwise single characters (f) return as ints (ascii character code and not strings, aka 102)
	//Using case forces the postgres driver to convert to a string that can be stored in a *string
	query := `SELECT CASE relreplident
	WHEN 'd' THEN 'default'
	WHEN 'n' THEN 'nothing'
	WHEN 'f' THEN 'full'
	WHEN 'i' THEN 'index'
END AS replica_identity
FROM pg_class
WHERE oid = $1::regclass`

	var resp string
	err := conn.QueryRow(context.Background(), query, tableName).Scan(&resp)
	if err != nil {
		return false, err
	}
	return resp == "full", nil
}

func resourceReplicaIdentityCreate(d *schema.ResourceData, m interface{}) error {
	conn, err := connect(d, m)
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	tableName := getTableName(d)

	boolean, err := IsReplicaIdentityFull(conn, tableName)
	if err != nil {
		return errors.Wrap(err, "resourceReplicaIdentityCreate: Error checking identity status")
	}

	if !boolean {
		_, err = conn.Exec(context.Background(), fmt.Sprintf("alter table %s replica identity full", tableName))
		if err != nil {
			return errors.Wrap(err, "resourceReplicaIdentityCreate: Error altering table")
		}
	}

	d.SetId(tableName)
	return nil
}

func resourceReplicaIdentityImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	err := resourceReplicationSlotRead(d, m)
	if err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func resourceReplicaIdentityRead(d *schema.ResourceData, m interface{}) error {
	conn, err := connect(d, m)
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	_, err = IsReplicaIdentityFull(conn, getTableName(d))
	if err != nil {
		return errors.Wrap(err, "resourceReplicaIdentityRead: Unable to read replica identity")
	}

	return nil
}

//Do not alter the identity if the terraform is destroyed
func resourceReplicaIdentityDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
