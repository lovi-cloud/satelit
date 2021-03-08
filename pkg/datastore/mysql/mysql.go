package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	// mysql driver
	_ "github.com/go-sql-driver/mysql"
	pb "github.com/lovi-cloud/teleskop/protoc/agent"

	"github.com/lovi-cloud/satelit/internal/client/teleskop"
	"github.com/lovi-cloud/satelit/pkg/config"
)

// A MySQL is backend of datastore by MySQL Server
type MySQL struct {
	Conn *sqlx.DB
}

// New create MySQL datastore
func New(c *config.MySQLConfig) (*MySQL, error) {
	dsn := c.DSN + "?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci"
	conn, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect : %w", err)
	}

	conn.SetMaxIdleConns(c.MaxIdleConn)
	conn.SetConnMaxLifetime(time.Duration(c.ConnMaxLifetimeSecond) * time.Second)

	return &MySQL{
		Conn: conn,
	}, nil
}

// GetIQN get IQN from MySQL
func (m *MySQL) GetIQN(ctx context.Context, hostname string) (string, error) {
	var iqn string

	query := fmt.Sprintf(`SELECT iqn FROM hypervisor WHERE hostname = "%s"`, hostname)
	err := m.Conn.Get(&iqn, query)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		break
	case err != nil:
		return "", fmt.Errorf("failed to exec query (hostname: %v): %w", hostname, err)
	}

	if iqn != "" {
		return iqn, nil
	}

	// not found in mysql
	teleskopClient, err := teleskop.GetClient(hostname)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve teleskop client: %w", err)
	}

	resp, err := teleskopClient.GetISCSIQualifiedName(ctx, &pb.GetISCSIQualifiedNameRequest{})
	if err != nil {
		return "", fmt.Errorf("failed to get IQN from Teleskop (host: %v): %w", hostname, err)
	}
	return resp.Iqn, nil
}
