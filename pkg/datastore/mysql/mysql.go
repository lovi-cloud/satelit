package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/whywaita/satelit/pkg/europa"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
	pb "github.com/whywaita/satelit/api"

	"github.com/whywaita/satelit/internal/client/teleskop"
	"github.com/whywaita/satelit/internal/config"
	"github.com/whywaita/satelit/internal/logger"
)

// A MySQL is backend of datastore by MySQL Server
type MySQL struct {
	Conn *sql.DB
}

// New create MySQL datastore
func New(c *config.MySQLConfig) (*MySQL, error) {
	dsn := c.DSN + "?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci"
	conn, err := sql.Open("mysql", dsn)
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
	row, err := m.Conn.Query("SELECT iqn FROM hypervisor WHERE hostname = ?", hostname)
	if err != nil {
		return "", fmt.Errorf("failed to exec query (host: %v): %w", hostname, err)
	}
	defer func() {
		if row != nil {
			errClose := row.Close()
			if errClose != nil {
				logger.Logger.Warn(err.Error())
			}
		}
	}()

	b := row.Next()
	if b == false {
		// not found
		iqn, err := teleskop.GetClient(hostname).GetISCSIQualifiedName(ctx, &pb.GetISCSIQualifiedNameRequest{})
		if err != nil {
			return "", fmt.Errorf("failed to get qin from Teleskop (host: %v): %w", hostname, err)
		}
		return iqn.Iqn, nil
	}

	var iqn string
	if err := row.Scan(iqn); err != nil {
		return "", fmt.Errorf("failed to scan MySQL response: %w", err)
	}

	return iqn, nil
}

// PutImage write image record
func (m *MySQL) PutImage(image europa.BaseImage) error {
	query, err := m.Conn.Prepare("INSERT INTO image(uuid, name, volume_id, description) VALUES (?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("failed to prepare query: %w", err)
	}

	_, err = query.Exec(image.ID, image.Name, image.CacheVolumeID, image.Description)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	return nil
}
