package config

import (
	"fmt"
	"io/ioutil"

	"github.com/goccy/go-yaml"

	"github.com/pkg/errors"
)

type Dorado struct {
	Username      string   `yaml:"username"`
	Password      string   `yaml:"password"`
	LocalIps      []string `yaml:"local_ip"`
	RemoteIps     []string `yaml:"remote_ip"`
	PortGroupName string   `yaml:"portgroup_name"`
}

type yml struct {
	Dorado   Dorado `yaml:"dorado"`
	LogLevel string `yaml:"log_level"`
}

var configContent yml

func Load(filepath *string) error {
	d, err := ioutil.ReadFile(*filepath)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to read config file: %s", filepath))
	}

	y := &yml{}
	if err := yaml.Unmarshal(d, y); err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to parse config file: %s", filepath))
	}

	configContent = *y

	return nil
}

func GetValue() yml {
	return configContent
}
