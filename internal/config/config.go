package config

import (
	"fmt"
	"io/ioutil"

	"github.com/goccy/go-yaml"

	"github.com/pkg/errors"
)

// A API is config of Satelit API Server
type API struct {
	Listen string `yaml:"listen"`
}

// A MySQLConfig is config of MySQL
type MySQLConfig struct {
	DSN                   string `yaml:"dsn"`
	MaxIdleConn           int    `yaml:"max_idle_conn"`
	ConnMaxLifetimeSecond int    `yaml:"conn_max_lifetime_second"`
}

// A Teleskop is endpoints of Teleskop
// example: host1: "teleskop_ip_and_port"
type Teleskop struct {
	Endpoints map[string]string `yaml:"endpoints"`
}

// A Dorado is config for Dorado
type Dorado struct {
	Username             string   `yaml:"username"`
	Password             string   `yaml:"password"`
	LocalIps             []string `yaml:"local_ip"`
	RemoteIps            []string `yaml:"remote_ip"`
	PortGroupName        string   `yaml:"portgroup_name"`
	StoragePoolName      string   `yaml:"storagepool_name"`
	HyperMetroDomainName string   `yaml:"hypermetrodomain_name"`
}

// A YAML is top element of config.yaml
type YAML struct {
	API         API         `yaml:"api"`
	MySQLConfig MySQLConfig `yaml:"mysql"`
	Teleskop    Teleskop    `yaml:"teleskop"`
	Dorado      Dorado      `yaml:"dorado"`
	LogLevel    string      `yaml:"log_level"`
}

var configContent YAML

// Load is to load yaml file
func Load(filepath *string) error {
	d, err := ioutil.ReadFile(*filepath)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to read config file: %s", *filepath))
	}

	y := &YAML{}
	if err := yaml.Unmarshal(d, y); err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to parse config file: %s", *filepath))
	}

	configContent = *y

	return nil
}

// GetValue return config values
func GetValue() YAML {
	return configContent
}
