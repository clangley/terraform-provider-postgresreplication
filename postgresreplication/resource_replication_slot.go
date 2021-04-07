package postgresreplication

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
)

func resourceReplicationSlot() *schema.Resource {
	return &schema.Resource{
		Create: resourceReplicationSlotCreate,
		Read:   resourceReplicationSlotRead,
		Delete: resourceReplicationSlotDelete,
		Importer: &schema.ResourceImporter{
			State: resourceReplicationSlotImport,
		},
		SchemaVersion: 0,
		Schema: map[string]*schema.Schema{
			slotNameAttributeName: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the slot to create. Must be a valid replication slot name.",
			},
			outputPluginAttributeName: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the output plugin used for logical decoding.",
			},
			databaseAttributeName: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the database this slot is associated with.",
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
	}
}

func resourceReplicationSlotCreate(d *schema.ResourceData, m interface{}) error {
	replConn, err := connect(d, m)
	if err != nil {
		return err
	}
	defer replConn.Close(context.Background())

	_, err = replConn.Exec(context.Background(), "SELECT * FROM pg_create_logical_replication_slot($1,$2);", getSlotName(d), getOutputPlugin(d))
	if err != nil {
		//Return error if it is NOT 42710, which is only returned if the replica slot already exists
		if !strings.Contains(err.Error(), "(SQLSTATE 42710)") {
			return errors.Wrap(err, "error creating replication slot.")
		}
	}

	d.SetId(getSlotName(d))

	return nil
}

func resourceReplicationSlotImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	err := resourceReplicationSlotRead(d, m)
	if err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func resourceReplicationSlotRead(d *schema.ResourceData, m interface{}) error {
	replConn, err := connect(d, m)
	if err != nil {
		return err
	}
	defer replConn.Close(context.Background())

	r, err := replConn.Query(context.Background(), "select slot_name, plugin, database from pg_replication_slots where slot_name=$1;", d.Id())
	if err != nil {
		return errors.Wrap(err, "error while trying to read existing replication slot")
	}
	defer r.Close()
	if r.Next() {
		v, _ := r.Values()
		err = d.Set(slotNameAttributeName, v[0])
		if err != nil {
			return errors.Wrap(err, "error reading slot name")
		}
		err = d.Set(outputPluginAttributeName, v[1])
		if err != nil {
			return errors.Wrap(err, "error reading output plugin")
		}
		err = d.Set(databaseAttributeName, v[2])
		if err != nil {
			return errors.Wrap(err, "error reading database")
		}
	}

	return nil
}

func resourceReplicationSlotDelete(d *schema.ResourceData, m interface{}) error {
	return resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		replConn, err := connect(d, m)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		defer replConn.Close(context.Background())

		_, err = replConn.Exec(context.Background(), "SELECT pg_drop_replication_slot($1);", getSlotName(d))
		if err != nil {
			return resource.RetryableError(errors.Wrap(err, "error dropping replication slot."))
		}

		return nil
	})
}
