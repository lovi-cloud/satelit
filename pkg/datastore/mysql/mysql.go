package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	pb "github.com/whywaita/satelit/api"

	"github.com/pkg/errors"
	"github.com/whywaita/satelit/internal/client/teleskop"
	"github.com/whywaita/satelit/internal/config"
	"github.com/whywaita/satelit/internal/logger"
)

type MySQL struct {
	Conn *sql.DB
}

func New(c *config.MySQLConfig) (*MySQL, error) {
	dsn := c.DSN + "?charset=utf8mb4&collation=utf8mb4_unicode_ci"
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect ")
	}

	conn.SetMaxIdleConns(c.MaxIdleConn)
	conn.SetConnMaxLifetime(time.Duration(c.ConnMaxLifetimeSecond) * time.Second)

	return &MySQL{
		Conn: conn,
	}, nil
}

func (m *MySQL) GetIQN(ctx context.Context, hostname string) (string, error) {
	row, err := m.Conn.Query("SELECT iqn FROM hypervisor WHERE hostname = ?", hostname)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("failed to exec query. hostname: %s", hostname))
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
			return "", errors.Wrap(err, fmt.Sprintf("failed to get iqn from Teleskop. hostname: %s", hostname))
		}
		return iqn.Iqn, nil
	}

	var iqn string
	if err := row.Scan(iqn); err != nil {
		return "", errors.Wrap(err, "failed to scan MySQL response")
	}

	return iqn, nil
}
