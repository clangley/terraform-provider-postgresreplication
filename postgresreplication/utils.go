package postgresreplication

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

const (
	slotNameAttributeName     = "slot_name"
	outputPluginAttributeName = "output_plugin"
	databaseAttributeName     = "database"

	tableNameAttributeName = "table_name"
)

func getTableName(d *schema.ResourceData) string {
	return d.Get(tableNameAttributeName).(string)
}

func getSlotName(d *schema.ResourceData) string {
	return d.Get(slotNameAttributeName).(string)
}

func getOutputPlugin(d *schema.ResourceData) string {
	return d.Get(outputPluginAttributeName).(string)
}

func connect(d *schema.ResourceData, m interface{}) (r *pgx.Conn, err error) {
	c := m.(*providerConfiguration)

	u, err := url.Parse(fmt.Sprintf("postgres://%s:%d/%s?sslmode=%s",
		c.host,
		c.port,
		d.Get(databaseAttributeName).(string),
		c.sslMode))
	if err != nil {
		return nil, errors.Wrap(err, "error contructing database connection uri.")
	}

	u.User = url.UserPassword(c.user, c.password)
	Conn, err := pgx.Connect(context.Background(), u.String())
	if err != nil {
		return nil, errors.Wrap(err, "error connecting to database.")
	}

	return Conn, nil
}
